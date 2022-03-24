package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"testing"
)

func Test_loadPackage(t *testing.T) {
	//pkg := loadPackage("D:/work/P/Server/LeafServer/src/server/db/lootmission.go")
	//pkg := loadPackage("D:\\work\\P\\robot\\db")
	pkg := loadPackage("D:/work/P/Server/LeafServer/src/server/db/",
		"lootmission.go", "crystal.go")
	t.Log("Syntax Size:", len(pkg.Syntax))
	for _, syntax := range pkg.Syntax {
		decls := parseModel(syntax)
		for _, decl := range decls {
			if decl.Name == "LootMission" {
				LoadLootMissions := decl.DbSelect().Where("PlayerSn").
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

func TestMysqlQuery(t *testing.T) {

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

const dataSourceName = "game:game123@tcp(127.0.0.1:3306)/game?charset=utf8mb4&parseTime=True&loc=Local"
