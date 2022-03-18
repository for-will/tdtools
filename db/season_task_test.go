package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/myPuffer/gotosql"
	"strings"
	"testing"
	"time"
)

type _SeasonTask struct {
	Id       int32 `db:"primary_key"`
	PlayerSn int32 `db:"unique:udx_season_player_task"`
	TaskId   int32 `db:"unique:udx_season_player_task"`
	Status   int32
	Progress int32
	Looped   int32
}

type _SeasonPlayer struct {
	Id            int32 `db:"primary_key"`
	PlayerSn      int32 `db:"unique:udx_season_player"`
	SeasonId      int32
	Premium       int32
	SeasonExp     int32
	TodayExp      int32
	DayTimeOut    time.Time `db:"type:timestamp,default:current_timestamp()"`
	WeekTimeOut   time.Time `db:"type:timestamp,default:current_timestamp()"`
	SeasonTimeOut time.Time `db:"type:timestamp,default:current_timestamp()"`
	Settled       bool
}

type _SeasonReward struct {
	Id       int32 `db:"primary_key"`
	PlayerSn int32 `db:"unique:udx_season_reward_player_lv"`
	Lv       int32 `db:"unique:udx_season_reward_player_lv"`
	Base     bool
	Premium  bool
}

func TestSeasonTask(t *testing.T) {

	GenModelAutoFile("season_task_out.go", "db",
		&AutoModel{
			Model: &_SeasonTask{},
			UserAdd: func(model *db.TableModel, sb *strings.Builder) {
				sb.WriteString("\n")
				sb.WriteString(model.BuildUpdateFunc("UpdateSeasonTask", "State", "Progress", "Loop"))
			},
		},
		&AutoModel{
			Model: &_SeasonPlayer{},
			UserAdd: func(model *db.TableModel, sb *strings.Builder) {
				sb.WriteString("\n")
			},
		},
		&AutoModel{
			Model: &_SeasonReward{},
			UserAdd: func(model *db.TableModel, sb *strings.Builder) {
				sb.WriteString("\n")
			},
		},
	)
}
