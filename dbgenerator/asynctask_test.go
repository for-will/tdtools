package main_test

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
	"text/template"
)

func TestParse(t *testing.T) {
	funcs := parseFuncDecls()

	for _, decl := range funcs {
		t.Log(GenAsyncTaskWrap(decl))
	}
}

func TestGenAsync(t *testing.T) {
	GenAsync()
}

func GenAsync() {
	var sb strings.Builder
	sb.WriteString(`package db

import (
	"github.com/name5566/leaf/log"
)
`)
	funcs := parseFuncDecls()

	for _, decl := range funcs {
		sb.WriteString(GenAsyncTaskWrap(decl))
	}

	out, err := format.Source([]byte(sb.String()))
	if err != nil {
		//panic(err)
		out = []byte(sb.String())
	}
	ioutil.WriteFile("task_async.go", out, 0664)
}

func GenAsyncTaskWrap(f *FuncDecl) string {
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

func parseFuncDecls() []*FuncDecl {
	fset := token.NewFileSet()
	src, err := ioutil.ReadFile("player_base.go")
	if err != nil {
		panic(err)
	}
	f, err := parser.ParseFile(fset, "player_base.go", src, parser.AllErrors)
	//f, err := parser.ParseDir(fset, ".", func(info fs.FileInfo) bool {
	//	return true
	//}, parser.AllErrors)
	if err != nil {
		panic(err)
	}

	var funcs []*FuncDecl
	for _, decl := range f.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			//t.Log(funcDecl.Name)
			fd := &FuncDecl{
				FuncName: funcDecl.Name.String(),
			}
			for _, field := range funcDecl.Type.Params.List {
				for _, name := range field.Names {
					objDecl := name.Obj.Decl.(*ast.Field)
					fieldType := src[objDecl.Type.Pos()-1 : objDecl.Type.End()-1]
					//t.Log(name, ":", string(fieldType))
					fd.In = append(fd.In, &FuncField{
						FieldName: name.String(),
						FieldType: string(fieldType),
					})
				}
			}
			if funcDecl.Type.Results != nil {

				for _, field := range funcDecl.Type.Results.List {
					for _, name := range field.Names {
						objDecl := name.Obj.Decl.(*ast.Field)
						fieldType := src[objDecl.Type.Pos()-1 : objDecl.Type.End()-1]
						fd.Out = append(fd.In, &FuncField{
							FieldName: name.String(),
							FieldType: string(fieldType),
						})
					}
					if field.Names == nil {
						fieldType := src[field.Type.Pos()-1 : field.Type.End()-1]
						fd.Out = append(fd.Out, &FuncField{
							FieldType: string(fieldType),
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
