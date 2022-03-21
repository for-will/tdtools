package main

import (
	"fmt"
	"log"
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

func (m *Model) DbUpdate(columns ...string) (sql string, args string, argsIn string) {

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
		if len(col.Type) > typeLen {
			typeLen = len(col.Type)
		}
	}

	var sb strings.Builder

	// Sql for create table
	sb.WriteString(fmt.Sprintf("CREATE OR REPLACE TABLE %s\n(\n", m.DbTableName()))
	format := fmt.Sprintf("    %%-%ds %%-%ds", nameLen, typeLen)
	for i, col := range m.Fields {
		sb.WriteString(fmt.Sprintf(format, snakeCase(col.Name), col.Type))
		if col.Name == "Id" {
			sb.WriteString(" not null auto_increment\n")
			sb.WriteString("        primary key")
		} else {
			sb.WriteString(" not null")
		}
		//if col.Default != "" {
		//	sb.WriteString(" default ")
		//	sb.WriteString(col.Default)
		//}
		if i != len(m.Fields)-1 {
			sb.WriteString(",\n")
		}
	}
	//
	//if len(tm.Indies) == 0 {
	//	sb.WriteString("\n")
	//}
	//
	//// Generate sql for create index
	//for _, v := range tm.Indies {
	//	sb.WriteString(",\n")
	//	if v.Unique {
	//		sb.WriteString("\tunique index ")
	//	} else {
	//		sb.WriteString("\tindex ")
	//	}
	//	sb.WriteString(v.Index)
	//	sb.WriteString("(")
	//	for i, col := range v.Keys {
	//		if i != 0 {
	//			sb.WriteString(", ")
	//		}
	//		sb.WriteString(col.KeyName)
	//	}
	//	sb.WriteString(") using btree")
	//}
	//
	//sb.WriteString("\n)")
	//
	return sb.String()
}
