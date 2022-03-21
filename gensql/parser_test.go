package main

import (
	"reflect"
	"testing"
)

func Test_extractTag(t *testing.T) {
	tag := reflect.StructTag(`db:"unique:udx_season_reward_player_lv,"`)

	//re := regexp.MustCompile(`db:"([\w:,]+)`)
	//matches := re.FindAllStringSubmatch(tag, -1)
	//t.Logf("%+v", matches[0][1])
	t.Log(tag.Get("db"))
}
