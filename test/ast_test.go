package main

import (
	"go/ast"
	"go/parser"
	"testing"
)

func TestParseExpr(t *testing.T) {
	expr, _ := parser.ParseExpr("\"123\"")
	ast.Print(nil, expr)

	expr, _ = parser.ParseExpr(`100 +		200*300`)
	ast.Print(nil, expr)
}
