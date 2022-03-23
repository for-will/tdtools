package main

import (
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	game := &GameSql{}
	//game.BuildLootMission()
	//game.BuildCrystal()
	game.BuildSeasonTask()
	game.BuildSeasonPlayer()
	game.BuildSeasonReward()
}
