package client

import (
	"github.com/myPuffer/gotosql"
	"google.golang.org/protobuf/proto"
	"reflect"
	"robot/GameMsg"
	"testing"
)

func TestReturnCode(t *testing.T) {
	msg := &GameMsg.EquipCrystalRs{ReturnCode: NewRetCode(GameMsg.ReturnCode_EquipmentSysNotOpen)}

	typ := reflect.TypeOf(msg).Elem()
	val := reflect.ValueOf(msg).Elem()
	for i := 0; i < typ.NumField(); i++ {
		//t.Log(typ.Field(i).Type, typ.Field(i).Type == reflect.TypeOf(new(GameMsg.ReturnCode)))
		if typ.Field(i).Type == reflect.TypeOf(new(GameMsg.ReturnCode)) {
			code := val.Field(i).Interface().(*GameMsg.ReturnCode)
			if *code != GameMsg.ReturnCode_OK {
				t.Log(code)
			}
		}

		//t.Log(val.Field(i).Interface().(GameMsg.ReturnCode))
		//
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
