package main

import (
	_ "github.com/go-sql-driver/mysql"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/packages"
	"io/ioutil"
	"log"
	"reflect"
	"strings"
)

func loadPackage(dir string, patterns ...string) *packages.Package {
	cfg := &packages.Config{
		Mode:  packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedName,
		Tests: false,
		Dir:   dir,
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages found", len(pkgs))
	}
	return pkgs[0]
}

func parseFile(path string) *FileSyntax {
	fset := token.NewFileSet()
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	file, err2 := parser.ParseFile(fset, path, content, parser.ParseComments)
	if err2 != nil {
		log.Fatal(err2)
	}
	return &FileSyntax{
		Fset:   fset,
		Syntax: file,
	}
}

func parseModel(file *ast.File) []*Model {

	var structs []*Model

	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			structs = append(structs, &Model{
				Name:   typeSpec.Name.Name,
				Fields: extractStructFields(structType.Fields),
			})
		}
	}
	return structs
}

func extractStructFields(fl *ast.FieldList) []*ModelField {
	var fields []*ModelField
	for _, field := range fl.List {
		//typ, ok := field.Type.(*ast.Ident)
		//if !ok {
		//	continue
		//}
		typeName := typeExprName(field.Type)
		if typeName == "" {
			continue
		}

		for _, name := range field.Names {
			fields = append(fields, &ModelField{
				Name: name.Name,
				Type: typeName,
				Tag:  extractTag(field.Tag),
			})
		}
	}
	return fields
}

func extractFunc(file *ast.File, funcName string) *ast.FuncDecl {

	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Name == nil || funcDecl.Name.Name != funcName {
			continue
		}

		return funcDecl
	}
	return nil
}

func editFunction(pkg *packages.Package, funcName string, funcText string) bool {
	for _, syntax := range pkg.Syntax {
		funcDecl := extractFunc(syntax, funcName)
		if funcDecl != nil {
			a, b := pkg.Fset.Position(funcDecl.Pos()), pkg.Fset.Position(funcDecl.End())
			text, _ := ioutil.ReadFile(a.Filename)
			src := string(text[:a.Offset]) + funcText + string(text[b.Offset:])
			out, err := format.Source([]byte(src))
			if err != nil {
				out = []byte(src)
			}

			ioutil.WriteFile(a.Filename, out, 0644)
			return true
		}
	}
	return false
}

func replaceFunction(fset *token.FileSet, syntax *ast.File, funcName string, funcText string) bool {
	funcDecl := extractFunc(syntax, funcName)
	if funcDecl != nil {
		a, b := fset.Position(funcDecl.Pos()), fset.Position(funcDecl.End())
		text, _ := ioutil.ReadFile(a.Filename)
		src := string(text[:a.Offset]) + funcText + string(text[b.Offset:])
		out, err := format.Source([]byte(src))
		if err != nil {
			out = []byte(src)
		}

		ioutil.WriteFile(a.Filename, out, 0644)
		return true
	}
	return false
}

func extractTag(Tag *ast.BasicLit) reflect.StructTag {
	if Tag == nil {
		return ""
	}

	return reflect.StructTag(strings.Trim(Tag.Value, "`"))

}

func typeExprName(Type ast.Expr) string {
	if ident, ok := Type.(*ast.Ident); ok && ident != nil {
		return ident.Name
	}
	if array, ok := Type.(*ast.ArrayType); ok && array != nil {
		return "[]" + typeExprName(array.Elt)
	}
	if star, ok := Type.(*ast.StarExpr); ok && star != nil {
		return "*" + typeExprName(star.X)
	}
	if sel, ok := Type.(*ast.SelectorExpr); ok && sel != nil {
		return typeExprName(sel.X) + "." + sel.Sel.Name
	}
	log.Fatalf("typeExprName: %+v", Type)
	return ""
}
