package main

import (
	"fmt"
	"strings"
	"text/template"
)

type Conditioner interface {
	Argument() string
	ArgumentName() string
	Condition() string
}

type FieldEqual struct {
	Field *ModelField
}

func (fe *FieldEqual) Condition() string {
	return fmt.Sprintf("%s = ?", snakeCase(fe.Field.Name))
}

func (fe *FieldEqual) ArgumentName() string {
	return PrivateFieldCase(fe.Field.Name)
}

func (fe *FieldEqual) Argument() string {
	return fe.ArgumentName() + " " + fe.Field.Type
}

type FieldIn struct {
	Field *ModelField
}

func (in *FieldIn) ArgumentName() string {
	return strings.ToLower(in.Field.Name[:1]) + in.Field.Name[1:] + "List"
}

func (in *FieldIn) ArgumentType() string {
	return "[]" + in.Field.Type
}

func (in *FieldIn) Argument() string {
	return in.ArgumentName() + " " + in.ArgumentType()
}

func (in *FieldIn) Condition() string {

	var sb = &strings.Builder{}

	text := `
	sqlSb.WriteString("{{.COLUMN}} IN(")
	for i, v := range {{.LIST}} {
		if i != 0 {
			sqlSb.WriteString(", ")
		}
		sqlSb.WriteString("?")
		args = append(args, v)
	}
	sqlSb.WriteString(")")
`
	tpl := template.New("sqlIn")
	tpl.Parse(text)

	tpl.Execute(sb, struct {
		COLUMN string
		LIST   string
	}{
		COLUMN: snakeCase(in.Field.Name),
		LIST:   in.ArgumentName(),
	})

	return sb.String()
}

func PrivateFieldCase(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}
