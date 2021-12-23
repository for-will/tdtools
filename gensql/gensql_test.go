package main

import (
	"io/ioutil"
	"testing"
)

func Test_loadPackage(t *testing.T) {
	//pkg := loadPackage("D:/work/P/Server/LeafServer/src/server/db/lootmission.go")
	//pkg := loadPackage("D:\\work\\P\\robot\\db")
	pkg := loadPackage("D:\\work\\P\\Server\\LeafServer\\src\\server\\db\\",
		"lootmission.go", "crystal.go")
	t.Log("Syntax Size:", len(pkg.Syntax))
	for _, syntax := range pkg.Syntax {
		decls := parseModel(syntax)
		for _, decl := range decls {
			if decl.Name == "LootMission" {
				LoadLootMissions := decl.DbSelect().
					Where(decl.FieldEqualCond("PlayerSn")).
					GenFixedQueryFunc("LoadLootMissions")
				t.Log(LoadLootMissions)
				editFunction(pkg, "LoadLootMissions", LoadLootMissions)

				//t.Log(decl.GenCreateFunc())
				//
				//t.Log(decl.GenBatchInsertFunc())

			}
		}
	}
}

func Test_extractFunc(t *testing.T) {
	pkg := loadPackage("D:/work/P/Server/LeafServer/src/server/db/lootmission.go")
	for _, syntax := range pkg.Syntax {
		funcDecl := extractFunc(syntax, "LoadLootMissions")
		if funcDecl != nil {
			a, b := pkg.Fset.Position(funcDecl.Pos()), pkg.Fset.Position(funcDecl.End())
			t.Log(a, b)
			text, _ := ioutil.ReadFile(a.Filename)
			t.Log(string(text[a.Offset:b.Offset]))
		}
	}
}
