package internal

import (
	"fmt"
	"go/ast"
	"go/format"
	"golang.org/x/tools/go/packages"
	"log"
	"market/GameMsg"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"text/template"
)

func GenHandlerWrap(dir string) string {

	var handlers []*FuncDecl
	pkg := parsePackage([]string{dir}, nil)
	for _, syntax := range pkg.Syntax {
		handlers = append(handlers, extractHandlerDecls(syntax)...)
	}
	if len(handlers) == 0 {
		log.Fatal("handlers empty")
	}

	var sb strings.Builder
	sb.WriteString("// Code generated by \"suituphandler\"; DO NOT EDIT.\n")
	sb.WriteString("\n")
	sb.WriteString("//go:generate suituphandler\n")
	sb.WriteString("\n")
	sb.WriteString(`package internal

import (
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"github.com/name5566/leaf/network/protobuf"
	"server/db"
	"server/msg/GameMsg"
)`)
	for _, h := range handlers {
		sb.WriteString(GenerateHandler(h))
	}

	sb.WriteString("\n\n")
	sb.WriteString(GenerateRegisterFunc(handlers))
	sb.WriteString("\n\n")
	sb.WriteString(GenBindMsgId())

	return sb.String()

}

func GenerateHandler(fi *FuncDecl) string {

	tplText := `

func {{.FunctionName}}(args []interface{}) {
	r := args[0].({{.RequestMsgName}})
	a := args[1].(gate.Agent)

	if a.UserData() == nil {
		log.Error("{{.FunctionName}} userdata is nil %s", a.RemoteAddr())
		return
	}

	player := a.UserData().(*db.Player)
	if player == nil {
		return
	}

	if rsp := {{.HandlerName}}({{.HandlerArgs}}); rsp != nil {
		a.WriteMsg(rsp)
	}
}`

	tpl := template.New("handler")
	tpl.Parse(tplText)

	var sb = &strings.Builder{}
	template.Must(tpl, tpl.Execute(sb, fi.GetHandlerInfo()))

	return sb.String()
}

func GenerateRegisterFunc(handlers []*FuncDecl) string {
	var sb strings.Builder
	sb.WriteString("func RegisterHandlers(){\n")
	for _, h := range handlers {
		info := h.GetHandlerInfo()
		//sb.WriteString(fmt.Sprintf("handler(&%s{}, %s)\n", info.RequestMsgName[1:], info.FunctionName))
		sb.WriteString(fmt.Sprintf("handler((%s)(nil), %s)\n", info.RequestMsgName, info.FunctionName))
	}
	sb.WriteString("}")
	return sb.String()
}

func GenBindMsgId() string {
	var msgIds []int32
	for id := range GameMsg.MsgId_name {
		msgIds = append(msgIds, id)
	}
	sort.Slice(msgIds, func(i, j int) bool {
		return msgIds[i] < msgIds[j]
	})

	var sb strings.Builder
	sb.WriteString("func BindMsgId(Processor *protobuf.Processor) {\n")
	for _, id := range msgIds {
		idName := GameMsg.MsgId_name[id]
		msgName := idName[4:]
		if msg, ok := SpecialMsgId[GameMsg.MsgId(id)]; ok {
			msgName = reflect.TypeOf(msg).Elem().Name()
		}

		s := fmt.Sprintf("\t//Processor.RegisterMsgAndID(&GameMsg.%s{}, uint32(GameMsg.MsgId_%s))\n",
			msgName, idName)
		sb.WriteString(s)
	}
	sb.WriteString("}")
	return sb.String()
}

func GenMsgIdMap() []byte {
	var msgIds []int32
	for id := range GameMsg.MsgId_name {
		msgIds = append(msgIds, id)
	}
	sort.Slice(msgIds, func(i, j int) bool {
		return msgIds[i] < msgIds[j]
	})

	var sb strings.Builder
	sb.WriteString(`package client

import (
	"market/GameMsg"
)

var MessageIdMap = map[GameMsg.MsgId]interface{}{
`)
	for _, id := range msgIds {
		idName := GameMsg.MsgId_name[id]
		msgName := idName[4:]
		if msg, ok := SpecialMsgId[GameMsg.MsgId(id)]; ok {
			msgName = reflect.TypeOf(msg).Elem().Name()
		}
		s := fmt.Sprintf("\tGameMsg.MsgId_%s:  (*GameMsg.%s)(nil),\n", idName, msgName)
		sb.WriteString(s)
	}
	sb.WriteString("}")
	out := sb.String()
	if formatted, err := format.Source([]byte(out)); err == nil {
		return formatted
	}
	return []byte(out)
}

type HandlerInfo struct {
	FunctionName   string
	RequestMsgName string
	HandlerName    string
	HandlerArgs    string
}

func parsePackage(patterns []string, tags []string) *packages.Package {
	cfg := &packages.Config{
		Mode:       packages.NeedFiles | packages.NeedSyntax | packages.NeedTypes | packages.NeedName,
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

func extractHandlerDecls(f *ast.File) []*FuncDecl {

	var funcs []*FuncDecl
	for _, decl := range f.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {

			fd := &FuncDecl{
				FuncName: funcDecl.Name.String(),
			}
			for _, field := range funcDecl.Type.Params.List {
				var pkgName, fieldType = typeExprName(field.Type)

				for _, name := range field.Names {
					fd.In = append(fd.In, &FuncField{
						FieldName: name.String(),
						FieldType: fieldType,
						PkgName:   pkgName,
					})
				}
			}
			if funcDecl.Type.Results != nil {

				for _, field := range funcDecl.Type.Results.List {
					pkgName, fieldType := typeExprName(field.Type)
					for _, name := range field.Names {
						fd.Out = append(fd.Out, &FuncField{
							FieldName: name.String(),
							FieldType: fieldType,
							PkgName:   pkgName,
						})
					}
					if field.Names == nil {
						fd.Out = append(fd.Out, &FuncField{
							FieldType: fieldType,
							PkgName:   pkgName,
						})
					}
				}
			}

			if fd.IsLegalHandler() {
				funcs = append(funcs, fd)
			}
		}
	}
	return funcs
}

type FuncField struct {
	FieldName string
	FieldType string
	PkgName   string
}

type FuncDecl struct {
	FuncName string
	In       []*FuncField
	Out      []*FuncField
}

func (fd *FuncDecl) IsLegalHandler() bool {
	if len(fd.Out) != 1 || fd.Out[0].PkgName != "GameMsg" || fd.Out[0].FieldType[0] != '*' {
		return false
	}

	var reqCnt int
	for _, field := range fd.In {
		if field.PkgName == "GameMsg" {
			if field.FieldType[0] != '*' {
				return false
			}
			reqCnt++
			continue
		}
		if field.FieldType == "gate.Agent" {
			continue
		}
		if field.FieldType == "*db.Player" {
			continue
		}
		return false
	}
	if reqCnt != 1 {
		return false
	}

	return true
}

func (fd *FuncDecl) RequestMsgName() string {
	for _, f := range fd.In {
		if f.PkgName == "GameMsg" {
			return f.FieldType
		}
	}
	return ""
}

func (fd *FuncDecl) GetHandlerInfo() *HandlerInfo {
	re, _ := regexp.Compile(`(.*)Handler$`)
	functionName := re.ReplaceAllString("handle"+strings.Title(fd.FuncName), "$1")
	info := &HandlerInfo{
		FunctionName:   functionName,
		RequestMsgName: fd.RequestMsgName(),
		HandlerName:    fd.FuncName,
		HandlerArgs:    fd.HandlerArgs(),
	}

	return info
}

func (fd *FuncDecl) HandlerArgs() string {
	var args []string
	for _, f := range fd.In {
		if f.PkgName == "GameMsg" {
			args = append(args, "r")
		} else if f.FieldType == "gate.Agent" {
			args = append(args, "a")
		} else if f.FieldType == "*db.Player" {
			args = append(args, "player")
		}

	}
	return strings.Join(args, ", ")
}

func typeExprName(Type ast.Expr) (pkg, typ string) {
	if ident, ok := Type.(*ast.Ident); ok && ident != nil {
		return "", ident.Name
	}
	if array, ok := Type.(*ast.ArrayType); ok && array != nil {
		_, name := typeExprName(array.Elt)
		return "", "[]" + name
	}
	if star, ok := Type.(*ast.StarExpr); ok && star != nil {
		pkgName, name := typeExprName(star.X)
		return pkgName, "*" + name
	}

	if star, ok := Type.(*ast.SelectorExpr); ok && star != nil {
		_, p := typeExprName(star.X)
		_, name := typeExprName(star.Sel)
		return p, p + "." + name
	}

	//log.Printf("typeExprName: %#v", Type)
	return "", "interface{}"
}

var SpecialMsgId = map[GameMsg.MsgId]interface{}{
	GameMsg.MsgId_C2S_QualityUp:         &GameMsg.HeroQualityUp{},
	GameMsg.MsgId_S2C_QualityUpRs:       &GameMsg.HeroQualityUpRs{},
	GameMsg.MsgId_S2C_SyncMainlineTask:  &GameMsg.SyncMainTask{},
	GameMsg.MsgId_C2S_Relogin:           &GameMsg.ReLogin{},
	GameMsg.MsgId_S2C_ReloginRs:         &GameMsg.ReLoginRs{},
	GameMsg.MsgId_C2C_MailList:          &GameMsg.MailInfoReq{},
	GameMsg.MsgId_S2C_MailListRs:        &GameMsg.MailInfoRes{},
	GameMsg.MsgId_C2C_MailDelete:        &GameMsg.MailDeleteReq{},
	GameMsg.MsgId_C2C_RankList:          &GameMsg.RankListReq{},
	GameMsg.MsgId_S2C_RankListRs:        &GameMsg.RankListRes{},
	GameMsg.MsgId_C2C_MailGetDesc:       &GameMsg.MailGetDescReq{},
	GameMsg.MsgId_C2S_VersionInfoReqID:  &GameMsg.C2S_VersionInfoReq{},
	GameMsg.MsgId_S2C_VersionInfoRspID:  &GameMsg.S2C_VersionInfoRsp{},
	GameMsg.MsgId_C2S_MailHeadListReqID: &GameMsg.C2S_MailHeadListReq{},
	GameMsg.MsgId_S2C_MailHeadListRspID: &GameMsg.S2C_MailHeadListRsp{},
	GameMsg.MsgId_C2S_MailBodyReqID:     &GameMsg.C2S_MailBodyReq{},
	GameMsg.MsgId_S2C_MailBodyRspID:     &GameMsg.S2C_MailBodyRsp{},
	GameMsg.MsgId_C2S_MailStateReqID:    &GameMsg.C2S_MailStateReq{},
	GameMsg.MsgId_S2C_MailStateRspID:    &GameMsg.S2C_MailStateRsp{},
	GameMsg.MsgId_C2S_MailBoxStateReqID: &GameMsg.C2S_MailBoxStateReq{},
	GameMsg.MsgId_S2C_MailBoxStateRspID: &GameMsg.S2C_MailBoxStateRsp{},
	GameMsg.MsgId_C2S_MailGetAwardReqID: &GameMsg.C2S_MailGetAwardReq{},
	GameMsg.MsgId_S2C_MailGetAwardRspID: &GameMsg.S2C_MailGetAwardRsp{},
	GameMsg.MsgId_C2C_MailGetAllAward:   &GameMsg.MailGetAllAwardReq{},
	GameMsg.MsgId_S2C_MailGetAllAwardRs: &GameMsg.MailGetAllAwardRsp{},
	GameMsg.MsgId_C2S_MailDeleteReqID:   &GameMsg.C2S_MailDeleteReq{},
	GameMsg.MsgId_S2C_MailDeleteRspID:   &GameMsg.S2C_MailDeleteRsp{},
	GameMsg.MsgId_S2C_HeroFightRs:       &GameMsg.HeroFight{},
	GameMsg.MsgId_S2C_NewMail:           &GameMsg.NewMailRs{},
}
