package main

import (
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/tools/go/packages"
	"log"
)

func main() {

	var reloadPackage = func() *packages.Package {
		return loadPackage("D:/work/P/Server/LeafServer/src/server/db/",
			"lootmission.go", "crystal.go")
	}
	pkg := reloadPackage()

	var models []*Model
	for _, syntax := range pkg.Syntax {
		models = append(models, parseModel(syntax)...)
	}

	var getModel = func(name string) *Model {
		for _, m := range models {
			if m.Name == name {
				return m
			}
		}
		log.Fatalf("getModel %s nil", name)
		return nil
	}

	var genSql = func(name string, gen func(model *Model)) {
		if m := getModel(name); m != nil {
			gen(m)
		}
	}

	// LootMission
	genSql("LootMission", func(model *Model) {
		LoadLootMissions := model.DbSelect().
			Where(model.FieldEqualCond("PlayerSn")).
			GenFixedQueryFunc("LoadLootMissions")
		editFunction(pkg, "LoadLootMissions", LoadLootMissions)
	})

	// Crystal
	genSql("Crystal", func(model *Model) {
		GetPlayerCrystals := model.DbSelect().
			Where(model.FieldEqualCond("PlayerId")).
			GenFixedQueryFunc("GetPlayerCrystals")
		editFunction(reloadPackage(), "GetPlayerCrystals", GetPlayerCrystals)

		BatchInsertCrystal := model.GenBatchInsertFunc()
		editFunction(reloadPackage(), "BatchInsertCrystal", BatchInsertCrystal)

		CreateCrystal := model.GenCreateFunc()
		editFunction(reloadPackage(), "CreateCrystal", CreateCrystal)

		UpdateCrystal := model.GenUpdateFunc("UpdateCrystal", "Locked", "Lv", "Expr")
		editFunction(reloadPackage(), "UpdateCrystal", UpdateCrystal)

		DeleteCrystals := model.Where(model.FieldInCond("Id")).GenDeleteFunc("DeleteCrystals")
		editFunction(reloadPackage(), "DeleteCrystals", DeleteCrystals)
		//log.Println(DeleteCrystals)
	})
}
