package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

var db *sql.DB

func init() {
	db = openDb()
}

func openDb() *sql.DB {
	db, err := sql.Open("mysql", "puffer:puffer123@tcp(127.0.0.1:3306)/gotosql?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	return db
}

func TestCreateTable(t *testing.T) {

	//db2.GenModelAutoFile("model_sql_auto.go", &TestTableTask{})
	NewTableTestTableTask(db)
}
