package main

import (
	"log"
	"robot/js"
	"testing"
)

func Test_parseFuncDecls(t *testing.T) {
	pkg := parsePackage([]string{"D:/work/P/Server/LeafServer/src/server/game/internal/lootmission.go"}, nil)
	//pkg := parsePackage([]string{"D:/work/P/Server/LeafServer/src/server/game/internal/"}, nil)

	for _, syntax := range pkg.Syntax {
		fds := extractHandlerDecls(syntax)
		t.Log(js.IdentJson(fds))
	}
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func TestGenHandlerWrap(t *testing.T) {
	out := GenHandlerWrap("D:/work/P/Server/LeafServer/src/server/game/internal/lootmission.go")
	println(out)
}
