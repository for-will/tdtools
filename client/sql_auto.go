//
// Code generated auto. DO NOT EDIT.

package client

import (
	"database/sql"
	"fmt"
)

func NewTableSyncPlayerBase(db *sql.DB) {
    sqlList := []string{
		`drop table if exists sync_player_base;`,
		`create table sync_player_base
(
    size_cache  int         not null,
    lv          int         not null,
    exp         int         not null,
    gold        int         not null,
    honor       int         not null,
    achievement int         not null,
    diamond     int         not null,
    name        varchar(64) not null,
    icon        int         not null
);`,
}

		for _, s := range sqlList {
		fmt.Println("Exec Sql:", s)
		_, err := db.Exec(s)
		if err != nil {
			fmt.Printf("db error: %+v", err)
		}
	}
}