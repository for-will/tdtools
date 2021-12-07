package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"golang.org/x/tools/go/packages"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

var typeNames = flag.String("type", "", "comma-separated list of type names; must be set")

func main() {
	flag.Parse()
	args := flag.Args()
	var dir string
	if len(args) != 0 {
		dir = filepath.Dir(args[0])
	} else {
		dir = "."
	}

	GenAsync(dir)
}

func GenAsync(dir string) {
	var funcDecls []*FuncDecl
	pkg := parsePackage([]string{dir}, nil)
	for _, file := range pkg.Syntax {
		funcDecls = append(funcDecls, parseFuncDecls(pkg.Fset, file)...)
	}

	var sb strings.Builder
	sb.WriteString(`package db

import (
	"github.com/name5566/leaf/log"
)
`)

	for _, decl := range funcDecls {
		sb.WriteString(GenTaskAsyncWrap(decl))
	}

	out, err := format.Source([]byte(sb.String()))
	if err != nil {
		out = []byte(sb.String())
	}
	ioutil.WriteFile("task_async.go", out, 0664)
}

func parsePackage(patterns []string, tags []string) *packages.Package {
	cfg := &packages.Config{
		Mode:       packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes,
		Tests:      false,
		BuildFlags: []string{fmt.Sprintf("-tags=%s", strings.Join(tags, " "))},
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

func parseFuncDecls(fset *token.FileSet, f *ast.File) []*FuncDecl {

	var funcs []*FuncDecl
	for _, decl := range f.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			if funcDecl.Doc == nil {
				continue
			}
			var toGenerate = false
			for _, comment := range funcDecl.Doc.List {
				fmt.Println(comment.Text)
				if comment.Text == "//go:generate asyncdb" {
					toGenerate = true
				}
			}
			if !toGenerate {
				continue
			}

			fd := &FuncDecl{
				FuncName: funcDecl.Name.String(),
			}
			for _, field := range funcDecl.Type.Params.List {
				for _, name := range field.Names {
					objDecl := name.Obj.Decl.(*ast.Field)
					fieldType := sourceCut(fset, objDecl.Type)
					fd.In = append(fd.In, &FuncField{
						FieldName: name.String(),
						FieldType: fieldType,
					})
				}
			}
			if funcDecl.Type.Results != nil {

				for _, field := range funcDecl.Type.Results.List {
					fieldType := sourceCut(fset, field.Type)
					for _, name := range field.Names {
						fd.Out = append(fd.In, &FuncField{
							FieldName: name.String(),
							FieldType: fieldType,
						})
					}
					if field.Names == nil {
						fd.Out = append(fd.Out, &FuncField{
							FieldType: fieldType,
						})
					}
				}
			}
			fd.ParamsIn = fd.MakeParamsIn()
			fd.ReturnsDecl = fd.MakeReturnsDecl()
			fd.ReturnsRcv = fd.MakeReturnsRcv()
			fd.CbFuncType = fd.MakeCbFuncType()
			funcs = append(funcs, fd)
		}
	}
	return funcs
}

func GenTaskAsyncWrap(f *FuncDecl) string {
	tplText := `

func (conn *DbConnect) {{.FuncName}}Async(task DbTask) {
	req, ok := task.ReqValue.(*struct {
		{{- range $val := .In}}
		{{$val.FieldName}} {{$val.FieldType}}
		{{- end}}
	})
	if !ok || req == nil {
		log.Error("{{.FuncName}}Async task: %#v", task.ReqValue)
		return
	}
	w := conn.worker(task.ObjID)
	if w == nil {
		log.Error("{{.FuncName}}Async conn.worker(%v) nil", task.ObjID)
	}

	{{.ReturnsDecl}}
	w.Go(func() {
		{{if .ReturnsRcv}}{{.ReturnsRcv}} = {{end}}conn.{{.FuncName}}({{.ParamsIn}})
	}, func() {
		if callback, ok := task.Cbi.({{.CbFuncType}}); ok && callback != nil {
			callback({{.ReturnsRcv}})
		}
	})
}

func (conn *DbConnect) Go{{.FuncName}}(account string, 
			{{- range $val := .In -}}
			{{$val.FieldName}} {{$val.FieldType}}, 
			{{- end -}} 
			cb {{.CbFuncType}}) {
	conn.{{.FuncName}}Async(DbTask{
		ObjID: account,
		ReqValue: &struct {
			{{- range $val := .In}}
			{{$val.FieldName}} {{$val.FieldType}}
			{{- end}}
		}{
			{{- range $val := .In}}
			{{$val.FieldName}}: {{$val.FieldName}},
			{{- end}}
		},
		Cbi: cb,
	})
}
`

	tpl := template.New("handler")
	tpl.Parse(tplText)

	var sb = &strings.Builder{}
	template.Must(tpl, tpl.Execute(sb, f))

	return sb.String()
}

func sourceCut(fset *token.FileSet, node ast.Node) string {
	p1 := fset.Position(node.Pos())
	p2 := fset.Position(node.End())
	src, err := ioutil.ReadFile(p1.Filename)
	if err != nil {
		return ""
	}
	return string(src[p1.Offset:p2.Offset])
}

type FuncField struct {
	FieldName string
	FieldType string
}

type FuncDecl struct {
	FuncName    string
	In          []*FuncField
	Out         []*FuncField
	ParamsIn    string
	ReturnsDecl string
	ReturnsRcv  string
	CbFuncType  string
}

func (decl FuncDecl) MakeParamsIn() string {
	var sb strings.Builder
	for i, v := range decl.In {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("req.")
		sb.WriteString(v.FieldName)
	}
	return sb.String()
}

func (decl FuncDecl) MakeReturnsDecl() string {
	var sb strings.Builder
	for i, v := range decl.Out {
		sb.WriteString("var ret")
		sb.WriteString(strconv.Itoa(i + 1))
		sb.WriteString(" ")
		sb.WriteString(v.FieldType)
		sb.WriteString("\n")
	}
	return sb.String()
}

func (decl FuncDecl) MakeReturnsRcv() string {
	var sb strings.Builder
	for i, _ := range decl.Out {
		if i != 0 {
			sb.WriteString(", ")
		}
		sb.WriteString("ret")
		sb.WriteString(strconv.Itoa(i + 1))
	}
	return sb.String()
}

func (decl FuncDecl) MakeCbFuncType() string {
	var sb strings.Builder
	sb.WriteString("func(")
	for i, v := range decl.Out {
		if i != 0 {
			sb.WriteString(", ")
		}
		if v.FieldName != "" {
			sb.WriteString(v.FieldName)
			sb.WriteString(" ")
		}
		sb.WriteString(v.FieldType)
	}
	sb.WriteString(")")
	return sb.String()
}
