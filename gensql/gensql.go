package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"go/ast"
	"go/format"
	"go/token"
	"golang.org/x/tools/go/packages"
	"io/ioutil"
	"log"
	"regexp"
)

const dataSourceName = "game:game123@tcp(127.0.0.1:3306)/game?charset=utf8mb4&parseTime=True&loc=Local"

func main() {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	rows, err2 := db.Query("SELECT id, card_sn FROM hero_talent_page WHERE id IN (?)",
		"'2',1")
	if rows != nil {
		defer rows.Close()
	}
	if err2 != nil {
		log.Fatal(err2)
	}

	for rows.Next() {
		var id, cardSn int32
		if err3 := rows.Scan(&id, &cardSn); err3 != nil {
			log.Printf("scan error: %v", err3)
		} else {
			log.Printf("id = %d, card_sn = %d", id, cardSn)
		}
	}

}

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
		typ, ok := field.Type.(*ast.Ident)
		if !ok {
			continue
		}
		typeName := typ.Name

		for _, name := range field.Names {
			fields = append(fields, &ModelField{
				Name: name.Name,
				Type: typeName,
				//JsonTag: extractJsonTag(field.Tag.Value),
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

func extractTag(meta string) string {
	re := regexp.MustCompile(`json:"(\w+)"`)
	matches := re.FindAllStringSubmatch(meta, -1)
	if len(matches) > 0 {
		return matches[0][1]
	}
	return ""
}
