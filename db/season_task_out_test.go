package db

import (
	"database/sql"
	"testing"
)

func openGameDB() *sql.DB {
	db, err := sql.Open("mysql", "game:game123@tcp(127.0.0.1:3306)/game?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	return db
}

func TestNewTblSeasonPlayer(t *testing.T) {
	NewTblSeasonPlayer(openGameDB())
}
