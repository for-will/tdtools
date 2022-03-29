package client

import (
	"github.com/myPuffer/gotosql"
	"google.golang.org/protobuf/proto"
	"market/GameMsg"
	"reflect"
	"testing"
)

func TestResponseMessage(t *testing.T) {
	var msg interface{} = &GameMsg.EquipCrystalRs{ReturnCode: GameMsg.ReturnCode_EquipmentSysNotOpen}

	if rsp, ok := msg.(ResponseMessage); ok {
		t.Log(rsp.GetReturnCode())
	}
}

func TestReturnCode(t *testing.T) {
	msg := &GameMsg.EquipCrystalRs{ReturnCode: GameMsg.ReturnCode_EquipmentSysNotOpen}

	typ := reflect.TypeOf(msg).Elem()
	val := reflect.ValueOf(msg).Elem()
	for i := 0; i < typ.NumField(); i++ {
		if typ.Field(i).Type == reflect.TypeOf(new(GameMsg.ReturnCode)) {
			code := val.Field(i).Interface().(*GameMsg.ReturnCode)
			if *code != GameMsg.ReturnCode_OK {
				t.Log(code)
			}
		}
	}
}

func TestEmptyMsg(t *testing.T) {
	hb := &GameMsg.HeartBeat{}
	b, err := proto.Marshal(hb)
	t.Log(len(b), err)
}

func TestGotosql(t *testing.T) {
	db.GenModelAutoFile("sql_auto.go", "", &GameMsg.SyncPlayerBase{})
}
