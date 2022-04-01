package db

import (
	"fmt"
	db "github.com/myPuffer/gotosql"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"
)

//CREATE USER 'game'@'%' IDENTIFIED BY 'game123';
//create database game default charset utf8 collate utf8_general_ci;
//grant all privileges on game.* to 'game'@'%' identified by 'game123';
//flush privileges;

var LogSql = db.LogSql
var LogError = db.LogError

type AutoModel struct {
	Model   interface{}
	UserAdd func(model *db.TableModel, sb *strings.Builder)
}

func GenModelAutoFile(file string, pkg string, models ...interface{}) {
	var sb = &strings.Builder{}
	sb.WriteString("//\n")
	sb.WriteString("// Code generated auto. DO NOT EDIT.\n\n")
	sb.WriteString("package " + pkg)
	sb.WriteString(`

import (
	"database/sql"
)`)

	for _, m := range models {

		var model *db.TableModel
		var UserAddFunc func(model *db.TableModel, sb *strings.Builder)
		if a, ok := m.(*AutoModel); ok {
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
		//sb.WriteString(model.BuildSaveFunc())
		//sb.WriteString("\n")
		//sb.WriteString(model.BuildFindOneFunc())
		//sb.WriteString("\n")
		//sb.WriteString(model.BuildFindFunc())

		if UserAddFunc != nil {
			UserAddFunc(model, sb)
		}
	}
	ioutil.WriteFile(file, []byte(sb.String()), 0664)

	info, err := exec.Command("go", "fmt", file).Output()
	fmt.Println(string(info), err)
}

type RewardTask struct {
	ReturnCode  *int32 `db:"index"`
	Id          *int32 `db:"primary_key"`
	ExpAdd      *int32
	Gold        *int32
	Honor       *int32
	Achievement *int32
	Diamond     *int32
}

type TestTableTask struct {
	Id         int32 `db:"name:sn,type:int,index,primary_key"`
	PlayerSn   int32 `db:"unique:idx_player_mission"`
	Mission    int32 `db:"mission,unique:idx_player_mission"`
	State      int8  `db:"unique"`
	Progress   int32
	RewardedAt time.Time `db:"type:timestamp"`
	PsX        int32
	PsY        int32
}

type _StorePurchase struct {
	Id        int32 `db:"primary_key"`
	PlayerSn  int32 `db:"unique:udx_internal_purchase_player_goods"`
	GoodsId   int32 `db:"unique:udx_internal_purchase_player_goods"`
	Purchased int32
	FreshTime time.Time `db:"type:timestamp"`
}

type _OpeningActivity struct {
	Id           int32 `db:"primary_key"`
	PlayerSn     int32 `db:"unique:udx_opening_activity_player_activity"`
	ActivityType int32 `db:"unique:udx_opening_activity_player_activity"`
	ActivityId   int32 `db:"unique:udx_opening_activity_player_activity"`
	State        int32
	Progress     int32
	StartAt      time.Time `db:"type:timestamp"`
}

type _ActivityTreasureBox struct {
	Id       int32 `db:"primary_key"`
	PlayerSn int32 `db:"unique:udx_activity_treasure_player_box"`
	BoxId    int32 `db:"unique:udx_activity_treasure_player_box"`
	Status   int32
}
