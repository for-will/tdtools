package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/packages"
	"log"
)

type GameSql struct {
	Models   []*Model
	OnModel  *Model
	PkgDir   string
	SrcFiles []string
}

const (
	GameDbDir              = "D:/work/P/Server/LeafServer/src/server/db/"
	LootMissionGoFileName  = "lootmission.go"
	CrystalGoFileName      = "crystal.go"
	SeasonTaskGoFileName   = "season_task.go"
	SeasonPlayerGoFileName = "season_player.go"
	SeasonRewardGoFileName = "season_reward.go"
	SignInGoFileName       = "signin.go"
	EquipGoFileName        = "equip.go"
	TowerGoFileName        = "tower.go"
)

type FileSyntax struct {
	Fset   *token.FileSet
	Syntax *ast.File
}

func (g *GameSql) ReloadPackage() *packages.Package {
	return loadPackage(g.PkgDir, g.SrcFiles...)
}

func (g *GameSql) ReloadFile() *FileSyntax {
	//var files []string
	//for _, file := range g.SrcFiles {
	//	files = append(files, g.PkgDir+file)
	//}
	return parseFile(g.PkgDir + g.SrcFiles[0])
}

func (g *GameSql) Init(ModelName string, PkgDir string, Files ...string) {
	g.PkgDir = PkgDir
	g.SrcFiles = Files

	//pkg := g.ReloadPackage()
	//for _, syntax := range pkg.Syntax {
	//	g.Models = append(g.Models, parseModel(syntax)...)
	//}
	file := g.ReloadFile()
	g.Models = append(g.Models, parseModel(file.Syntax)...)
	g.OnModel = g.GetModel(ModelName)
}

func (g *GameSql) GetModel(name string) *Model {
	for _, m := range g.Models {
		if m.Name == name {
			return m
		}
	}
	log.Fatalf("GetModel %s nil", name)
	return nil
}

func (g *GameSql) GenerateMethod(MethodName string, Generator func(model *Model, MethodName string) string) {
	log.Output(2, fmt.Sprintf("GenerateMethod: %s", MethodName))
	MethodSrc := Generator(g.OnModel, MethodName)
	//editFunction(g.ReloadPackage(), MethodName, MethodSrc)
	f := g.ReloadFile()
	replaceFunction(f.Fset, f.Syntax, MethodName, MethodSrc)
}

func (g *GameSql) BuildLootMission() {
	g.Init("LootMission", GameDbDir, LootMissionGoFileName)

	g.GenerateMethod("LoadLootMissions", func(model *Model, MethodName string) string {
		return model.DbSelect().Where("PlayerSn").GenFixedQueryFunc(MethodName)
	})
}

//func (g *GameSql) BuildCrystal() {
//	g.Init("Crystal", GameDbDir, CrystalGoFileName)
//
//	g.GenerateMethod("GetPlayerCrystals", func(model *Model, MethodName string) string {
//		return model.DbSelect().Where("PlayerId").
//			GenFixedQueryFunc(MethodName)
//	})
//
//	g.GenerateMethod("BatchInsertCrystal", func(model *Model, MethodName string) string {
//		return model.GenBatchInsertFunc()
//	})
//
//	g.GenerateMethod("CreateCrystal", func(model *Model, MethodName string) string {
//		return model.GenCreateFunc()
//	})
//
//	g.GenerateMethod("UpdateCrystal", func(model *Model, MethodName string) string {
//		return model.GenUpdateFunc(MethodName, "Locked", "Lv", "Expr")
//	})
//
//	g.GenerateMethod("DeleteCrystals", func(model *Model, MethodName string) string {
//		return model.Where("Id IN (?)").GenDeleteFunc(MethodName)
//	})
//}

func (g *GameSql) BuildSeasonPlayer() {
	g.Init("SeasonPlayer", GameDbDir, SeasonPlayerGoFileName)

	g.GenerateMethod("NewTblSeasonPlayer", func(model *Model, MethodName string) string {
		return model.GenNewTblFunc()
	})

	g.GenerateMethod("LoadSeasonPlayer", func(model *Model, MethodName string) string {
		return model.DbSelect().Where("PlayerSn").
			GenFixedQueryFunc(MethodName)
	})

	g.GenerateMethod("CreateSeasonPlayer", func(model *Model, MethodName string) string {
		return model.GenCreateFunc()
	})

	g.GenerateMethod("UpdateSeasonPlayerFields", func(model *Model, MethodName string) string {
		return model.GenUpdateFunc(MethodName, "SeasonId", "Premium", "SeasonExp",
			"TodayExp", "DayTimeOut", "WeekTimeOut", "SeasonTimeOut", "Settled")
	})
}

func (g *GameSql) BuildSeasonTask() {

	g.Init("SeasonTask", GameDbDir, SeasonTaskGoFileName)

	g.GenerateMethod("NewTblSeasonTask", func(model *Model, MethodName string) string {
		return model.GenNewTblFunc()
	})

	g.GenerateMethod("LoadSeasonTasks", func(model *Model, MethodName string) string {
		return model.DbSelect().Where("PlayerSn").GenFixedQueryFunc(MethodName)
	})

	g.GenerateMethod("BatchInsertSeasonTask", func(model *Model, MethodName string) string {
		return model.GenBatchInsertFunc()
	})

	g.GenerateMethod("UpdateSeasonTaskProgress", func(model *Model, MethodName string) string {
		return model.GenUpdateFunc(MethodName, "Progress", "Status", "Looped", "Rewarded")
	})

	g.GenerateMethod("ResetSeasonTask", func(model *Model, MethodName string) string {
		return model.Where("PlayerSn").
			GenUpdateFunc(MethodName, "Progress", "Status", "Looped", "Rewarded")
	})

	g.GenerateMethod("ResetPlayerSeasonTasks", func(model *Model, MethodName string) string {
		return model.Where("Id IN (?)").
			GenBatchUpdateFunc(MethodName, "Progress", "Status", "Looped", "Rewarded")
	})
}

func (g GameSql) BuildSeasonReward() {
	g.Init("SeasonReward", GameDbDir, SeasonRewardGoFileName)

	g.GenerateMethod("LoadSeasonReward", func(model *Model, MethodName string) string {
		return model.DbSelect().Where("PlayerSn").
			GenFixedQueryFunc(MethodName)
	})

	g.GenerateMethod("BatchInsertSeasonReward", func(model *Model, MethodName string) string {
		return model.GenBatchInsertFunc()
	})

	g.GenerateMethod("UpdateSeasonReward", func(model *Model, MethodName string) string {
		return model.GenUpdateFunc(MethodName, "Base", "Premium")
	})
}

func (g GameSql) BuildSignIn() {
	g.Init("DailySignIn", GameDbDir, SignInGoFileName)

	g.GenerateMethod("NewTblDailySignIn", func(model *Model, MethodName string) string {
		return model.GenNewTblFunc()
	})

	g.GenerateMethod("LoadPlayerDailySign", func(model *Model, MethodName string) string {
		return model.DbSelect().Where("PlayerSn").
			GenFixedQueryFunc(MethodName)
	})

	g.GenerateMethod("CreateDailySignIn", func(model *Model, MethodName string) string {
		return model.GenCreateFunc()
	})

	g.GenerateMethod("UpdateDailySignIn", func(model *Model, MethodName string) string {
		return model.GenUpdateFunc(MethodName, "Signed", "NextDay", "NextMonth")
	})
}

func (g GameSql) BuildEquip() {
	g.Init("Equip", GameDbDir, EquipGoFileName)

	g.GenerateMethod("NewTblEquip", func(model *Model, MethodName string) string {
		return model.GenNewTblFunc()
	})

	g.GenerateMethod("GetPlayerEquips", func(model *Model, MethodName string) string {
		return model.DbSelect().Where("PlayerSn").
			GenFixedQueryFunc(MethodName)
	})

	g.GenerateMethod("CreateEquip", func(model *Model, MethodName string) string {
		return model.GenCreateFunc()
	})

	g.GenerateMethod("BatchInsertEquip", func(model *Model, MethodName string) string {
		return model.GenBatchInsertFunc()
	})

	g.GenerateMethod("UpdateEquip", func(model *Model, MethodName string) string {
		return model.GenUpdateFunc(MethodName,
			"Lv",
			"Expr",
			"Locked",
			"Attr",
			"MinorCnt",
			"Attr1", "Lv1",
			"Attr2", "Lv2",
			"Attr3", "Lv3",
			"Attr4", "Lv4",
			"Hero")
	})

	g.GenerateMethod("DeleteEquips", func(model *Model, MethodName string) string {
		return model.Where("Id IN (?)").GenDeleteFunc(MethodName)
	})
}

func (g GameSql) BuildTower() {
	g.Init("Tower", GameDbDir, TowerGoFileName)

	g.GenerateMethod("NewTblTower", func(model *Model, MethodName string) string {
		return model.GenNewTblFunc()
	})

	g.GenerateMethod("GetPlayerTowers", func(model *Model, MethodName string) string {
		return model.DbSelect().Where("PlayerSn").
			GenFixedQueryFunc(MethodName)
	})

	g.GenerateMethod("CreateTower", func(model *Model, MethodName string) string {
		return model.GenCreateFunc()
	})

	g.GenerateMethod("UpdateTower", func(model *Model, MethodName string) string {
		return model.GenUpdateFunc(MethodName, "Lv", "InUse")
	})
}
