package main

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {

	game := &GameSql{}
	game.BuildLootMission()
	game.BuildSeasonTask()
	game.BuildSeasonPlayer()
	game.BuildSeasonReward()
	game.BuildSignIn()
	game.BuildEquip()
	game.BuildTower()
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
