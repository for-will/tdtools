package main

import "testing"

func TestGameSql_BuildSignIn(t *testing.T) {
	game := &GameSql{}
	game.BuildSignIn()
}
