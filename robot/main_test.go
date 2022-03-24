package main

import (
	"encoding/base64"
	jsoniter "github.com/json-iterator/go"
	"testing"
)

func Test_runAutoTestRobot(t *testing.T) {

	var op = &operates{
		Server:       "172.16.1.218:3563",
		UserAccount:  "x11",
		UserPassword: "",
		ExploreArea:  5,
		ExploreType:  1,
		ExploreCnt:   3,
	}

	jsons, _ := jsoniter.MarshalToString(op)
	t.Log(jsons)

	t.Log(base64.StdEncoding.EncodeToString([]byte(jsons)))
}
