package main

import (
	"log"
	"strings"
)

func (m *Model) SqlUpdate(columns ...string) (sql string, args string, argsIn string) {

	var sb strings.Builder
	sb.WriteString("UPDATE ")
	sb.WriteString(m.DbTableName())
	sb.WriteString(" SET ")

	var ArgsName []string
	var ArgsDecl []string
	var ArgsSet []string
	for _, col := range columns {
		field := m.GetField(col)
		if field == nil {
			log.Fatalf("DbUpdate %s GetField(%s) nil", m.Name, col)
		}
		ArgsSet = append(ArgsSet, field.SqlName()+"=?")
		ArgsName = append(ArgsName, field.Name)
		ArgsDecl = append(ArgsDecl, field.Name+" "+field.Type)
	}
	sb.WriteString(strings.Join(ArgsSet, ", "))
	sb.WriteString(" WHERE")

	return sb.String(), strings.Join(ArgsName, ", "), strings.Join(ArgsDecl, ", ")
}

func (m *Model) SqlUpdateById(columns ...string) (sql string, args string, argsIn string) {

	var sb strings.Builder
	sb.WriteString("UPDATE ")
	sb.WriteString(m.DbTableName())
	sb.WriteString(" SET ")

	argsIn = "Id int32"
	for i, col := range columns {
		if i != 0 {
			sb.WriteString(", ")
			args += ", "
		}
		argsIn += ", "

		field := m.GetField(col)
		if field == nil {
			log.Fatalf("DbUpdate %s GetField(%s) nil", m.Name, col)
			return "", "", ""
		}
		sb.WriteString(snakeCase(field.Name) + "=?")
		args += field.Name
		argsIn += field.Name + " " + field.Type
	}
	sb.WriteString(" WHERE id=?")
	args += ", Id"

	return sb.String(), args, argsIn
}
