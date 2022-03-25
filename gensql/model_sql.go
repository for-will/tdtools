package main

import (
	"fmt"
	"strings"
)

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
	sb.WriteString(")")
	sb.WriteString(" VALUES ")

	return sb.String(), "(?" + strings.Repeat(", ?", len(m.Fields)-2) + ")"
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

func (m *Model) DbCreateTbl() (sql string) {
	var nameLen, typeLen int
	for _, col := range m.Fields {
		if len(snakeCase(col.Name)) > nameLen {
			nameLen = len(snakeCase(col.Name))
		}
		if len(col.SqlType()) > typeLen {
			typeLen = len(col.SqlType())
		}
	}

	var sb strings.Builder

	// Sql for create table
	sb.WriteString(fmt.Sprintf("CREATE OR REPLACE TABLE %s\n(\n", m.DbTableName()))
	format := fmt.Sprintf("    %%-%ds %%-%ds", nameLen, typeLen)
	for i, col := range m.Fields {
		sb.WriteString(fmt.Sprintf(format, snakeCase(col.Name), col.SqlType()))
		if col.Name == "Id" {
			sb.WriteString(" not null auto_increment\n")
			sb.WriteString("        primary key")
		} else {
			sb.WriteString(" not null")
		}
		if col.SqlDefault() != "" {
			sb.WriteString(" default ")
			sb.WriteString(col.SqlDefault())
		}
		if i != len(m.Fields)-1 {
			sb.WriteString(",\n")
		}
	}

	indies := modelIndies(m)

	if len(indies) == 0 {
		sb.WriteString("\n")
	}

	// Generate sql for create index
	for _, v := range indies {
		sb.WriteString(",\n")
		if v.Unique {
			sb.WriteString("\tunique index ")
		} else {
			sb.WriteString("\tindex ")
		}
		sb.WriteString(v.Index)
		sb.WriteString("(")
		for i, col := range v.Keys {
			if i != 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(col.KeyName)
		}
		sb.WriteString(") using btree")
	}

	sb.WriteString("\n)")

	return sb.String()
}
