package main

import (
	"testing"
)

func Test_loadPackage(t *testing.T) {
	structs := parseLogStruct()
	for _, decl := range structs {
		for _, field := range decl.Fields {
			t.Log(field.JsonTag)
		}
	}
}

func TestExtractJsonTag(t *testing.T) {
	t.Log(extractJsonTag("`json:\"sauth_login_type\"`"))
}
