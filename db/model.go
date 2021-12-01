package db

import (
	db "github.com/myPuffer/gotosql"
	"time"
)

var LogSql = db.LogSql
var LogError = db.LogError

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
