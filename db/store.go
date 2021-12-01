//
// Code generated auto. DO NOT EDIT.

package db

import (
	"database/sql"
	"strings"
	"time"
)

type StorePurchase struct {
	Id        int32
	PlayerSn  int32
	GoodsId   int32
	Purchased int32
	FreshTime time.Time
}

func NewTblStorePurchase(db *sql.DB) {

	querySql := `create or replace table store_purchase
(
    id         int       not null auto_increment
        primary key,
    player_sn  int       not null,
    goods_id   int       not null,
    purchased  int       not null,
    fresh_time timestamp not null,
	unique index udx_internal_purchase_player_goods(player_sn, goods_id) using btree
)`
	LogSql(querySql)
	_, err := db.Exec(querySql)
	if err != nil {
		LogError("%+v", err)
	}
}

func SaveStorePurchase(db *sql.DB, obj *StorePurchase) error {

	LogSql("INSERT INTO store_purchase(player_sn, goods_id, purchased, fresh_time) VALUE ('%v', '%v', '%v', '%v')",
		obj.PlayerSn, obj.GoodsId, obj.Purchased, obj.FreshTime)
	result, err := db.Exec("INSERT INTO store_purchase(player_sn, goods_id, purchased, fresh_time) VALUE (?, ?, ?, ?)",
		obj.PlayerSn, obj.GoodsId, obj.Purchased, obj.FreshTime)

	if err != nil {
		LogError("db query error: %+v", err)
		return err
	}

	if id, err := result.LastInsertId(); err != nil {
		LogError("get last insert id error: %+v", err)
		return err
	} else {
		obj.Id = int32(id)
	}
	return nil
}

func FirstStorePurchase(db *sql.DB, out *StorePurchase, cond string, args ...interface{}) error {

	sfmt := strings.Replace(cond, "?", "'%+v'", -1)
	LogSql("SELECT id, player_sn, goods_id, purchased, fresh_time FROM store_purchase WHERE "+sfmt+" LIMIT 1", args...)
	rows, err := db.Query("SELECT id, player_sn, goods_id, purchased, fresh_time FROM store_purchase WHERE "+cond+" LIMIT 1", args...)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		LogError("db query error: %+v", err)
		return err
	}

	if rows.Next() {
		if err = rows.Scan(&out.Id, &out.PlayerSn, &out.GoodsId, &out.Purchased, &out.FreshTime); err != nil {
			LogError("db query Scan error: %+v", err)
			return err
		}
	} else {
		return sql.ErrNoRows
	}
	return nil
}

func FindStorePurchase(db *sql.DB, out *[]*StorePurchase, cond string, args ...interface{}) error {

	sfmt := strings.Replace(cond, "?", "'%+v'", -1)
	LogSql("SELECT id, player_sn, goods_id, purchased, fresh_time FROM store_purchase WHERE "+sfmt, args...)
	rows, err := db.Query("SELECT id, player_sn, goods_id, purchased, fresh_time FROM store_purchase WHERE "+cond, args...)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		LogError("db query error: %+v", err)
		return err
	}

	for rows.Next() {
		obj := &StorePurchase{}
		if err = rows.Scan(&obj.Id, &obj.PlayerSn, &obj.GoodsId, &obj.Purchased, &obj.FreshTime); err != nil {
			LogError("db query Scan error: %+v", err)
			return err
		}
		*out = append(*out, obj)
	}
	return nil
}

func UpdateStorePurchased(db *sql.DB, where_id int32, purchased int32, fresh_time time.Time) error {

	_, err := db.Exec("UPDATE store_purchase SET purchased = ?, fresh_time = ? WHERE id = ?",
		purchased, fresh_time, where_id)

	if err != nil {
		LogError("db query error: %+v", err)
		LogSql("UPDATE store_purchase SET purchased = '%v', fresh_time = '%v' WHERE id = '%v'",
			purchased, fresh_time, where_id)
		return err
	}

	return nil
}

type OpeningActivity struct {
	Id           int32
	PlayerSn     int32
	ActivityType int32
	ActivityId   int32
	State        int32
	Progress     int32
	StartAt      time.Time
}

func NewTblOpeningActivity(db *sql.DB) {

	querySql := `create or replace table opening_activity
(
    id            int       not null auto_increment
        primary key,
    player_sn     int       not null,
    activity_type int       not null,
    activity_id   int       not null,
    state         int       not null,
    progress      int       not null,
    start_at      timestamp not null,
	unique index udx_opening_activity_player_activity(player_sn, activity_type, activity_id) using btree
)`
	LogSql(querySql)
	_, err := db.Exec(querySql)
	if err != nil {
		LogError("%+v", err)
	}
}

func SaveOpeningActivity(db *sql.DB, obj *OpeningActivity) error {

	LogSql("INSERT INTO opening_activity(player_sn, activity_type, activity_id, state, progress, start_at) VALUE ('%v', '%v', '%v', '%v', '%v', '%v')",
		obj.PlayerSn, obj.ActivityType, obj.ActivityId, obj.State, obj.Progress, obj.StartAt)
	result, err := db.Exec("INSERT INTO opening_activity(player_sn, activity_type, activity_id, state, progress, start_at) VALUE (?, ?, ?, ?, ?, ?)",
		obj.PlayerSn, obj.ActivityType, obj.ActivityId, obj.State, obj.Progress, obj.StartAt)

	if err != nil {
		LogError("db query error: %+v", err)
		return err
	}

	if id, err := result.LastInsertId(); err != nil {
		LogError("get last insert id error: %+v", err)
		return err
	} else {
		obj.Id = int32(id)
	}
	return nil
}

func FirstOpeningActivity(db *sql.DB, out *OpeningActivity, cond string, args ...interface{}) error {

	sfmt := strings.Replace(cond, "?", "'%+v'", -1)
	LogSql("SELECT id, player_sn, activity_type, activity_id, state, progress, start_at FROM opening_activity WHERE "+sfmt+" LIMIT 1", args...)
	rows, err := db.Query("SELECT id, player_sn, activity_type, activity_id, state, progress, start_at FROM opening_activity WHERE "+cond+" LIMIT 1", args...)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		LogError("db query error: %+v", err)
		return err
	}

	if rows.Next() {
		if err = rows.Scan(&out.Id, &out.PlayerSn, &out.ActivityType, &out.ActivityId, &out.State, &out.Progress, &out.StartAt); err != nil {
			LogError("db query Scan error: %+v", err)
			return err
		}
	} else {
		return sql.ErrNoRows
	}
	return nil
}

func FindOpeningActivity(db *sql.DB, out *[]*OpeningActivity, cond string, args ...interface{}) error {

	sfmt := strings.Replace(cond, "?", "'%+v'", -1)
	LogSql("SELECT id, player_sn, activity_type, activity_id, state, progress, start_at FROM opening_activity WHERE "+sfmt, args...)
	rows, err := db.Query("SELECT id, player_sn, activity_type, activity_id, state, progress, start_at FROM opening_activity WHERE "+cond, args...)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		LogError("db query error: %+v", err)
		return err
	}

	for rows.Next() {
		obj := &OpeningActivity{}
		if err = rows.Scan(&obj.Id, &obj.PlayerSn, &obj.ActivityType, &obj.ActivityId, &obj.State, &obj.Progress, &obj.StartAt); err != nil {
			LogError("db query Scan error: %+v", err)
			return err
		}
		*out = append(*out, obj)
	}
	return nil
}

func UpdateOpeningActivity(db *sql.DB, where_id int32, progress int32, state int32) error {

	_, err := db.Exec("UPDATE opening_activity SET progress = ?, state = ? WHERE id = ?",
		progress, state, where_id)

	if err != nil {
		LogError("db query error: %+v", err)
		LogSql("UPDATE opening_activity SET progress = '%v', state = '%v' WHERE id = '%v'",
			progress, state, where_id)
		return err
	}

	return nil
}

type ActivityTreasureBox struct {
	Id       int32
	PlayerSn int32
	BoxId    int32
	Status   int32
}

func NewTblActivityTreasureBox(db *sql.DB) {

	querySql := `create or replace table activity_treasure_box
(
    id        int not null auto_increment
        primary key,
    player_sn int not null,
    box_id    int not null,
    status    int not null,
	unique index udx_activity_treasure_player_box(player_sn, box_id) using btree
)`
	LogSql(querySql)
	_, err := db.Exec(querySql)
	if err != nil {
		LogError("%+v", err)
	}
}

func SaveActivityTreasureBox(db *sql.DB, obj *ActivityTreasureBox) error {

	LogSql("INSERT INTO activity_treasure_box(player_sn, box_id, status) VALUE ('%v', '%v', '%v')",
		obj.PlayerSn, obj.BoxId, obj.Status)
	result, err := db.Exec("INSERT INTO activity_treasure_box(player_sn, box_id, status) VALUE (?, ?, ?)",
		obj.PlayerSn, obj.BoxId, obj.Status)

	if err != nil {
		LogError("db query error: %+v", err)
		return err
	}

	if id, err := result.LastInsertId(); err != nil {
		LogError("get last insert id error: %+v", err)
		return err
	} else {
		obj.Id = int32(id)
	}
	return nil
}

func FirstActivityTreasureBox(db *sql.DB, out *ActivityTreasureBox, cond string, args ...interface{}) error {

	sfmt := strings.Replace(cond, "?", "'%+v'", -1)
	LogSql("SELECT id, player_sn, box_id, status FROM activity_treasure_box WHERE "+sfmt+" LIMIT 1", args...)
	rows, err := db.Query("SELECT id, player_sn, box_id, status FROM activity_treasure_box WHERE "+cond+" LIMIT 1", args...)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		LogError("db query error: %+v", err)
		return err
	}

	if rows.Next() {
		if err = rows.Scan(&out.Id, &out.PlayerSn, &out.BoxId, &out.Status); err != nil {
			LogError("db query Scan error: %+v", err)
			return err
		}
	} else {
		return sql.ErrNoRows
	}
	return nil
}

func FindActivityTreasureBox(db *sql.DB, out *[]*ActivityTreasureBox, cond string, args ...interface{}) error {

	sfmt := strings.Replace(cond, "?", "'%+v'", -1)
	LogSql("SELECT id, player_sn, box_id, status FROM activity_treasure_box WHERE "+sfmt, args...)
	rows, err := db.Query("SELECT id, player_sn, box_id, status FROM activity_treasure_box WHERE "+cond, args...)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		LogError("db query error: %+v", err)
		return err
	}

	for rows.Next() {
		obj := &ActivityTreasureBox{}
		if err = rows.Scan(&obj.Id, &obj.PlayerSn, &obj.BoxId, &obj.Status); err != nil {
			LogError("db query Scan error: %+v", err)
			return err
		}
		*out = append(*out, obj)
	}
	return nil
}

func UpdateTreasureBox(db *sql.DB, where_id int32, status int32) error {

	_, err := db.Exec("UPDATE activity_treasure_box SET status = ? WHERE id = ?",
		status, where_id)

	if err != nil {
		LogError("db query error: %+v", err)
		LogSql("UPDATE activity_treasure_box SET status = '%v' WHERE id = '%v'",
			status, where_id)
		return err
	}

	return nil
}
