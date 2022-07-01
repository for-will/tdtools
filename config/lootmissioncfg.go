package config

import (
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"market/GameMsg"
	"unsafe"
)

type LootMission struct {
	Id            int                    `json:"Id"`
	FrontId       int                    `json:"FrontId"`
	Open          int                    `json:"Open"`
	Lv            int                    `json:"Lv"`
	Name          int                    `json:"Name"`
	Describe      int                    `json:"Describe"`
	Type          int                    `json:"Type"`
	Grouping      int                    `json:"Grouping"`
	Condition     GameMsg.TASK_CONDITION `json:"Condition"`
	ConditionDesc int                    `json:"ConditionDesc"`
	Param         int                    `json:"Param"`
	RewardCfg     int                    `json:"RewardCfg"`
	Exp           int                    `json:"Exp"`
	Sort          int                    `json:"Sort"`
	TheSpoils     int                    `json:"Thespoils"`
}

var LootMissionCfg []*LootMission

func LoadConfig() []*LootMission {

	jsoniter.RegisterTypeEncoderFunc("GameMsg.TASK_CONDITION", func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
		stream.WriteString((*GameMsg.TASK_CONDITION)(ptr).String())
	}, nil)

	jsoniter.RegisterTypeDecoderFunc("GameMsg.TASK_CONDITION", func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		switch iter.WhatIsNext() {
		case jsoniter.NumberValue:
			*(*GameMsg.TASK_CONDITION)(ptr) = GameMsg.TASK_CONDITION(iter.ReadInt32())
		case jsoniter.StringValue:
			*(*GameMsg.TASK_CONDITION)(ptr) = GameMsg.TASK_CONDITION(GameMsg.TASK_CONDITION_value[iter.ReadString()])
		default:
			*(*interface{})(ptr) = iter.Read()
		}
	})

	b, err := ioutil.ReadFile("../bin/LootMissionsCfg.json")
	if err != nil {
		panic(err)
	}
	err = jsoniter.Unmarshal(b, &LootMissionCfg)
	if err != nil {
		panic(err)
	}

	//s, _ := jsoniter.MarshalIndent(LootMissionCfg, "", "    ")
	//fmt.Printf("%+v", string(s))
	return LootMissionCfg
}
