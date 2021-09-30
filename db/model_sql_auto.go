//
// Code generated auto. DO NOT EDIT.

package db

import (
	"database/sql"
	"fmt"
)

func NewTableTestTableTask(db *sql.DB) {
    sqlList := []string{
		`drop table if exists test_table_task;`,
		`create table test_table_task
(
    id          int       not null auto_increment
        primary key,
    player_sn   int       not null,
    mission     int       not null,
    state       tinyint   not null,
    progress    int       not null,
    rewarded_at timestamp not null,
    ps_x        int       not null,
    ps_y        int       not null
);`,
		`create unique index idx_player_mission on test_table_task (player_sn, mission);`,
		`create unique index uni_test_table_task_state on test_table_task (state);`,
		`create index idx_test_table_task_id on test_table_task (id);`,
}

		for _, s := range sqlList {
		fmt.Println("Exec Sql:", s)
		_, err := db.Exec(s)
		if err != nil {
			fmt.Printf("db error: %+v", err)
		}
	}
}