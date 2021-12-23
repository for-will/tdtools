package main

import (
	"log"
	"strings"
	"text/template"
)

type Model struct {
	Name         string
	Fields       []*ModelField
	SelectFields []*ModelField
	Sql          string
	Conditions   []Conditioner
}

type ModelField struct {
	Name string
	Type string
	Tag  string
}

func (m *Model) clone() *Model {

	cloneModel := *m
	cloneModel.Fields = []*ModelField{}
	for _, field := range m.Fields {
		cloneField := *field
		cloneModel.Fields = append(cloneModel.Fields, &cloneField)
	}
	return &cloneModel
}

func (m *Model) DbTableName() string {
	return snakeCase(m.Name)
}

func (m *Model) SqlQuery() string {
	out := m.Sql
	if len(m.Conditions) != 0 {
		out += " WHERE"
		for _, cond := range m.Conditions {
			out += " "
			out += cond.Condition()
		}
	}
	return out
}

func (m *Model) CondBuild() string {

	var sb strings.Builder
	for _, cond := range m.Conditions {
		sb.WriteString(cond.Condition())
	}
	return sb.String()
}

func (m *Model) DbSelect() *Model {

	var sb strings.Builder
	sb.WriteString("SELECT ")
	for i, field := range m.Fields {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(snakeCase(field.Name))
	}
	sb.WriteString(" FROM ")
	sb.WriteString(m.DbTableName())

	tx := m.clone()
	tx.Sql = sb.String()
	return tx
}

func (m *Model) DbInsert() (insert string, place string) {

	var sb strings.Builder
	sb.WriteString("INSERT INTO ")
	sb.WriteString(m.DbTableName())
	sb.WriteString("(")
	for i, field := range m.Fields[1:] {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(snakeCase(field.Name))
	}
	sb.WriteString(" VALUES ")

	return sb.String(), "(?" + strings.Repeat(", ?", len(m.Fields)-2) + ")"
}

func (m *Model) DbUpdate(columns ...string) (sql string, args string) {

	var sb strings.Builder
	sb.WriteString("UPDATE ")
	sb.WriteString(m.DbTableName())
	sb.WriteString(" SET ")
	for i, col := range columns {
		if i != 0 {
			sb.WriteString(", ")
			args += ", "
		}
		field := m.GetField(col)
		if field == nil {
			log.Fatalf("DbUpdate %s GetField(%s) nil", m.Name, col)
			return "", ""
		}
		sb.WriteString(snakeCase(field.Name) + "=?")
		args += "obj." + field.Name
	}
	sb.WriteString(" WHERE id=?")
	args += ", obj.Id"

	return sb.String(), args
}

func (m *Model) DbDelete() (sql string) {

	var sb strings.Builder
	sb.WriteString("DELETE FROM ")
	sb.WriteString(m.DbTableName())
	if len(m.Conditions) > 0 {
		sb.WriteString(" WHERE ")
	}
	return sb.String()
}

func (m *Model) Where(conds ...Conditioner) *Model {
	tx := m.clone()
	tx.Conditions = append(tx.Conditions, conds...)
	return tx
}

func (m *Model) GetField(name string) *ModelField {
	for _, field := range m.Fields {
		if field.Name == name {
			return field
		}
	}
	return nil
}

func (m *Model) FieldEqualCond(fieldName string) Conditioner {
	field := m.GetField(fieldName)
	if field == nil {
		log.Fatalf("FieldEqualCond GetField(%s) nil", fieldName)
	}
	return &FieldEqual{
		Field: field,
	}
}

func (m *Model) FieldInCond(fieldName string) Conditioner {
	field := m.GetField(fieldName)
	if field == nil {
		log.Fatalf("FieldEqualCond GetField(%s) nil", fieldName)
	}
	return &FieldIn{
		Field: field,
	}
}

func (m *Model) GenFixedQueryFunc(funcName string) string {

	text := `func (conn *DbConnect) {{.FUNC}}({{.IN}}) []*{{.STRUCT}} {

	rows, err1 := conn.db.Query(
		"{{.SQL}}",
		{{.ARGS}})

	if rows != nil {
		defer rows.Close()
	}

	if err1 != nil {
		log.Error("{{.FUNC}}({{.DUMP_FMT}}) failed: %v", {{.ARGS}}, err1)
		return nil
	}

	var retList []*{{.STRUCT}}
	for rows.Next() {
		item := &{{.STRUCT}}{}
		err2 := rows.Scan({{.SCAN_LIST}})
		if err2 != nil {
			log.Error("{{.FUNC}}({{.DUMP_FMT}}) Scan error: %v", {{.ARGS}}, err2)
			continue
		}
		retList = append(retList, item)
	}
	return retList
}`
	tpl := template.New(funcName)
	tpl.Parse(text)

	var sb = &strings.Builder{}
	tpl.Execute(sb, &struct {
		STRUCT    string
		FUNC      string
		IN        string
		SQL       string
		ARGS      string
		SCAN_LIST string
		DUMP_FMT  string
	}{
		STRUCT:    m.Name,
		FUNC:      funcName,
		IN:        m.FuncIn(),
		SQL:       m.SqlQuery(),
		ARGS:      m.FuncArgs(),
		SCAN_LIST: m.FuncScanList(),
		DUMP_FMT:  m.FuncDumpFmt(),
	})
	return sb.String()
}

func (m *Model) GenCreateFunc() string {

	text := `func (conn *DbConnect) {{.FUNC}}(obj *{{.STRUCT}}) *{{.STRUCT}} {

	res, err1 := conn.db.Exec("{{.SQL}}", 
		{{.ARGS}})

	if err1 != nil {
		log.Error("{{.FUNC}} (%+v) err:%v", obj, err1)
		return nil
	}

	if id, err2 := res.LastInsertId(); err2 != nil {
		log.Error("{{.FUNC}} (%+v) err:%v", obj, err2)
		return nil
	} else {
		obj.Id = int32(id)
		return obj
	}
}
`
	tpl := template.New("GenCreateFunc:" + m.Name)
	tpl.Parse(text)

	var args []string
	for _, field := range m.Fields[1:] {
		args = append(args, "obj."+field.Name)
	}
	var sb = &strings.Builder{}
	insert, place := m.DbInsert()
	tpl.Execute(sb, &struct {
		STRUCT string
		FUNC   string
		IN     string
		SQL    string
		ARGS   string
	}{
		STRUCT: m.Name,
		FUNC:   "Create" + m.Name,
		IN:     m.FuncIn(),
		SQL:    insert + place,
		ARGS:   strings.Join(args, ", "),
	})
	return sb.String()
}

func (m *Model) GenBatchInsertFunc() string {

	text := `func (conn *DbConnect) {{.FUNC}}(objList []*{{.STRUCT}}) ([]*{{.STRUCT}}, error) {

	if len(objList) == 0 {
		return nil, nil
	}

	var sqlSb strings.Builder
	var sqlArgs = make([]interface{}, 0, {{.FIELD_CNT}}*len(objList))
	sqlSb.WriteString("{{.SQL}}")
	for i, obj := range objList {
		if i != 0 {
			sqlSb.WriteString(", ")
		}
		sqlSb.WriteString("{{.PLACEHOLDERS}}")
		sqlArgs = append(sqlArgs, {{.ARGS}})
	}

	result, err1 := conn.db.Exec(sqlSb.String(), sqlArgs...)

	if err1 != nil {
		log.Error("{{.FUNC}} exec %s, %+v error: %v",
			sqlSb.String(), sqlArgs, err1)
		return nil, err1
	}

	if id, err2 := result.LastInsertId(); err2 != nil {
		log.Error("{{.FUNC}} get last insert id error: %+v", err2)
		return nil, err2
	} else {
		for i, obj := range objList {
			obj.Id = int32(i) + int32(id)
		}
	}
	return objList, nil
}
`
	tpl := template.New("GenBatchInsertFunc:" + m.Name)
	tpl.Parse(text)

	var args []string
	for _, field := range m.Fields[1:] {
		args = append(args, "obj."+field.Name)
	}
	var sb = &strings.Builder{}
	insert, place := m.DbInsert()
	tpl.Execute(sb, &struct {
		STRUCT       string
		FUNC         string
		SQL          string
		PLACEHOLDERS string
		ARGS         string
		FIELD_CNT    int
	}{
		STRUCT:       m.Name,
		FUNC:         "BatchInsert" + m.Name,
		SQL:          insert,
		PLACEHOLDERS: place,
		ARGS:         strings.Join(args, ", "),
		FIELD_CNT:    len(m.Fields) - 1,
	})
	return sb.String()
}

func (m *Model) GenUpdateFunc(funcName string, columns ...string) string {
	if len(columns) == 0 {
		log.Fatal("GenUpdateFunc columns nil.")
		return ""
	}

	if funcName == "" {
		funcName = "Update" + m.Name
	}

	text := `func (conn *DbConnect) {{.FUNC}}(obj *{{.STRUCT}}) bool {

	result, err1 := conn.db.Exec("{{.SQL}}", 
		{{.ARGS}})

	if err1 != nil {
		log.Error("{{.FUNC}} (%+v) failed:%v", obj, err1)
		return false
	}

	if r, err2 := result.RowsAffected(); err2 != nil || r != 1 {
		log.Error("{{.FUNC}} (%+v) rows_affected=%d, err:%v", obj, r, err2)
		return false
	}
	return true
}
`
	tpl := template.New("GenUpdateFunc:" + m.Name)
	tpl.Parse(text)

	var sb = &strings.Builder{}
	sql, args := m.DbUpdate(columns...)
	tpl.Execute(sb, &struct {
		STRUCT string
		FUNC   string
		SQL    string
		ARGS   string
	}{
		STRUCT: m.Name,
		FUNC:   funcName,
		SQL:    sql,
		ARGS:   args,
	})
	return sb.String()
}

func (m *Model) GenDeleteFunc(funcName string) string {

	text := `func (conn *DbConnect) {{.FUNC}}({{.IN}}) bool {

	var args []interface{}
	var sqlSb strings.Builder
	sqlSb.WriteString("{{.SQL}}")
	{{.COND_BUILD}}
	result, err1 := conn.db.Exec(sqlSb.String(), args...)
	if err1 != nil {
		log.Error("{{.FUNC}}({{.DUMP_FMT}}) failed: %v", {{.ARGS}}, err1)
		return false
	}

	if r, err2 := result.RowsAffected(); err2 != nil || r == 0 {
		log.Error("{{.FUNC}}({{.DUMP_FMT}}) rows_affected=%d, err:%v", {{.ARGS}}, r, err2)
		return false
	}
	return true
}
`
	tpl := template.New(funcName)
	tpl.Parse(text)

	var sb = &strings.Builder{}
	tpl.Execute(sb, &struct {
		STRUCT     string
		FUNC       string
		IN         string
		SQL        string
		COND_BUILD string
		ARGS       string
		DUMP_FMT   string
	}{
		STRUCT:     m.Name,
		FUNC:       funcName,
		IN:         m.FuncIn(),
		SQL:        m.DbDelete(),
		COND_BUILD: m.CondBuild(),
		ARGS:       m.FuncArgs(),
		DUMP_FMT:   m.FuncDumpFmt(),
	})
	return sb.String()
}

func (m *Model) FuncIn() string {
	var args []string
	for _, c := range m.Conditions {
		if a := c.Argument(); a != "" {
			args = append(args, a)
		}
	}
	return strings.Join(args, ", ")
}

func (m *Model) FuncArgs() string {
	var args []string
	for _, c := range m.Conditions {
		args = append(args, c.ArgumentName())
	}
	return strings.Join(args, ", ")
}

func (m *Model) FuncScanList() string {
	var args []string
	for _, f := range m.Fields {
		args = append(args, "&item."+f.Name)
	}
	return strings.Join(args, ", ")
}

func (m *Model) FuncDumpFmt() string {
	var args []string
	for _, c := range m.Conditions {
		args = append(args, c.ArgumentName()+"=%v")
	}
	return strings.Join(args, ", ")
}

func snakeCase(s string) string {

	var sb strings.Builder
	for i, r := range s {
		if r >= 'A' && r <= 'Z' && i != 0 {
			sb.WriteString("_")
		}
		sb.WriteString(strings.ToLower(string(r)))
	}
	return sb.String()
}
