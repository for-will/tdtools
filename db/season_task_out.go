//
// Code generated auto. DO NOT EDIT.

package db

import (
	"database/sql"
	"time"
)

type SeasonTask struct {
	Id       int32
	PlayerSn int32
	TaskId   int32
	Status   int32
	Progress int32
	Looped   int32
}

func NewTblSeasonTask(db *sql.DB) {

	querySql := `create or replace table season_task
(
    id        int not null auto_increment
        primary key,
    player_sn int not null,
    task_id   int not null,
    status    int not null,
    progress  int not null,
    looped    int not null,
	unique index udx_season_player_task(player_sn, task_id) using btree
)`
	LogSql(querySql)
	_, err := db.Exec(querySql)
	if err != nil {
		LogError("%+v", err)
	}
}

func UpdateSeasonTask(db *sql.DB, where_id int32, progress int32) error {

	_, err := db.Exec("UPDATE season_task SET progress = ? WHERE id = ?",
		progress, where_id)

	if err != nil {
		LogError("db query error: %+v", err)
		LogSql("UPDATE season_task SET progress = '%v' WHERE id = '%v'",
			progress, where_id)
		return err
	}

	return nil
}

type SeasonPlayer struct {
	Id            int32
	PlayerSn      int32
	SeasonId      int32
	Premium       int32
	SeasonExp     int32
	TodayExp      int32
	DayTimeOut    time.Time
	WeekTimeOut   time.Time
	SeasonTimeOut time.Time
}

func NewTblSeasonPlayer(db *sql.DB) {

	querySql := `create or replace table season_player
(
    id              int       not null auto_increment
        primary key,
    player_sn       int       not null,
    season_id       int       not null,
    premium         int       not null,
    season_exp      int       not null,
    today_exp       int       not null,
    day_time_out    timestamp not null,
    week_time_out   timestamp not null,
    season_time_out timestamp not null,
	unique index udx_season_player(player_sn) using btree
)`
	LogSql(querySql)
	_, err := db.Exec(querySql)
	if err != nil {
		LogError("%+v", err)
	}
}
