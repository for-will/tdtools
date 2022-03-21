package main

import (
	"golang.org/x/tools/go/packages"
	"log"
	"testing"
)

func TestModel_DbCreateTbl(t *testing.T) {
	var reloadPackage = func() *packages.Package {
		return loadPackage("D:/work/P/Server/LeafServer/src/server/db/",
			"season_reward.go", "season_player.go")
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

	var genFunc = func(name string, gen func(model *Model)) {
		if m := getModel(name); m != nil {
			gen(m)
		}
	}

	// SeasonPlayer
	genFunc("SeasonReward", func(model *Model) {
		Sql := model.DbCreateTbl()
		t.Log(Sql)

		//for _, field := range model.Fields {
		//	t.Log(field)
		//}
	})

	// SeasonPlayer
	genFunc("SeasonPlayer", func(model *Model) {
		Sql := model.DbCreateTbl()
		t.Log(Sql)

		//for _, field := range model.Fields {
		//	t.Log(field)
		//}
	})
}
