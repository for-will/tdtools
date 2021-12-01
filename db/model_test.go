package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/myPuffer/gotosql"
	"io/ioutil"
	"os/exec"
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

func TestCreateTable(t *testing.T) {

	//db2.GenModelAutoFile("model_sql_auto.go", &TestTableTask{})
	NewTableTestTableTask(_db)
}

func GenModelAutoFile(file string, pkg string, models ...interface{}) {
	var sb = &strings.Builder{}
	sb.WriteString("//\n")
	sb.WriteString("// Code generated auto. DO NOT EDIT.\n\n")
	sb.WriteString("package " + pkg)
	sb.WriteString(`

import (
	"database/sql"
	"strings"
)`)

	for _, m := range models {

		var model *db.TableModel
		var UserAddFunc func(model *db.TableModel, sb *strings.Builder)
		if a, ok := m.(*DbAutoModel); ok {
			model = db.Model(a.Model)
			UserAddFunc = a.UserAdd
		} else {
			model = db.Model(m)
		}

		model.ModName = strings.TrimLeft(model.ModName, "_")
		model.TblName = strings.TrimLeft(model.TblName, "_")
		//model := db.Model(m)
		//model.LogSql = "log.Debug"
		//model.LogError = "log.Error"

		sb.WriteString("\n")
		sb.WriteString(model.TypeStruct())
		sb.WriteString("\n")
		sb.WriteString(model.BuildCreateTableFunc())
		sb.WriteString("\n")
		sb.WriteString(model.BuildSaveFunc())
		sb.WriteString("\n")
		sb.WriteString(model.BuildFindOneFunc())
		sb.WriteString("\n")
		sb.WriteString(model.BuildFindFunc())

		if UserAddFunc != nil {
			UserAddFunc(model, sb)
		}
	}
	ioutil.WriteFile(file, []byte(sb.String()), 0664)

	info, err := exec.Command("go", "fmt", file).Output()
	fmt.Println(string(info), err)
}

type DbAutoModel struct {
	Model   interface{}
	UserAdd func(model *db.TableModel, sb *strings.Builder)
}

func TestBana(t *testing.T) {

	GenModelAutoFile("store.go", "db",
		&DbAutoModel{
			Model: &_StorePurchase{},
			UserAdd: func(model *db.TableModel, sb *strings.Builder) {
				sb.WriteString("\n")
				sb.WriteString(model.BuildUpdateFunc("UpdateStorePurchased", "Purchased", "FreshTime"))
			},
		},
		&DbAutoModel{
			Model: &_OpeningActivity{},
			UserAdd: func(model *db.TableModel, sb *strings.Builder) {
				sb.WriteString("\n")
				sb.WriteString(model.BuildUpdateFunc("UpdateOpeningActivity", "Progress", "State"))
			},
		},
		&DbAutoModel{
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
