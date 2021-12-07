package main

import (
	"go/ast"
	"testing"
)

func Test_parsePackage(t *testing.T) {
	pkg := parsePackage([]string{"D:/work/P/robot/db/"}, nil)
	for i := range pkg.GoFiles {
		syntax := pkg.Syntax[i]
		for _, decl := range syntax.Decls {
			if funcDecl, ok := decl.(*ast.FuncDecl); ok {

				for _, field := range funcDecl.Type.Params.List {
					for _, name := range field.Names {
						objDecl := name.Obj.Decl.(*ast.Field)
						cut := sourceCut(pkg.Fset, objDecl.Type)
						t.Log(cut)
					}
				}
			}
		}
	}
}

func TestGenAsync(t *testing.T) {
	GenAsync("D:/work/P/robot/db/")
}
