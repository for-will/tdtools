package client

import (
	"fmt"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"market/GameMsg"
	"reflect"
	"regexp"
	"testing"
	"time"
)

func Test_fetchReturnCode(t *testing.T) {
	msg := &GameMsg.ModifyPlayerNameRs{ReturnCode: GameMsg.ReturnCode_ModifyNameUsed}
	fetchReturnCode(msg)
	fmt.Println(reflect.TypeOf(new(GameMsg.ReturnCode)).Elem())
}

func TestRegex(t *testing.T) {
	exp, _ := regexp.Compile(`[\s\=\(\)（）*"']`)
	ok := exp.MatchString("ab'1c")
	t.Log(ok)
}

func TestTimeJson(t *testing.T) {
	t.Log(JsonString(time.Now()))

	//_ = main.Generator{}
	var tm time.Time
	//jsoniter.UnmarshalFromString("\"2021-11-03T14:08:57+08:00\"", &tm)
	//time.Parse()
	tm.UnmarshalJSON([]byte("\"2021-11-03T00:00:00+08:00\""))
	t.Log(JsonString(tm))
}

func TestInt32Sz(t *testing.T) {
	msg := &GameMsg.TestInt32Size{
		A1: 1,
		A2: -1,
	}
	data, _ := proto.Marshal(msg)
	jsb, _ := protojson.Marshal(msg)
	t.Log(len(data), string(jsb))
}
