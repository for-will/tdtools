package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/myPuffer/gotosql"
	"reflect"
	"strings"
	"testing"
	"time"
)

var _db *sql.DB

func init() {
	_db = openDb()
}

func openDb() *sql.DB {
	db, err := sql.Open("mysql", "puffer:puffer123@tcp(127.0.0.1:3306)/gotosql?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	return db
}



func TestBana(t *testing.T) {

	GenModelAutoFile("store.go", "db",
		&AutoModel{
			Model: &_StorePurchase{},
			UserAdd: func(model *db.TableModel, sb *strings.Builder) {
				sb.WriteString("\n")
				sb.WriteString(model.BuildUpdateFunc("UpdateStorePurchased", "Purchased", "FreshTime"))
			},
		},
		&AutoModel{
			Model: &_OpeningActivity{},
			UserAdd: func(model *db.TableModel, sb *strings.Builder) {
				sb.WriteString("\n")
				sb.WriteString(model.BuildUpdateFunc("UpdateOpeningActivity", "Progress", "State"))
			},
		},
		&AutoModel{
			Model: &_ActivityTreasureBox{},
			UserAdd: func(model *db.TableModel, sb *strings.Builder) {
				sb.WriteString("\n")
				sb.WriteString(model.BuildUpdateFunc("UpdateTreasureBox", "Status"))
			},
		},
	)
}

func TestTimeType(t *testing.T) {

	typ := reflect.TypeOf(&time.Time{}).Elem()
	t.Log(typ.Kind() == reflect.Struct, typ.String())

	//tNow := time.Now()
	//t.Log(tNow.Unix(), tNow)
	//t.Log(tNow.UTC().Unix(), tNow.UTC())
}
