package internal

import (
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"io/ioutil"
)

type TaskCfg struct {
	ID        uint32 `gorm:"primarykey"`
	TaskId    int32  `json:"Id"`
	Front     int32  `json:"Front"`
	Type      int32  `json:"Type"`
	Group     int32  `json:"Group"`
	Sort      int32  `json:"Sort"`
	Name      int32  `json:"Name"`
	Comment   string `json:"Comment" gorm:"type:varchar(100)"`
	Describe  int32  `json:"Describe"`
	Condition int32  `json:"Condition"`
	TaskCond  string `json:"TaskCond" gorm:"type:varchar(100)"`
	ParamX    int32  `json:"ParamX"`
	ParamY    int32  `json:"ParamY"`
	Count     int32  `json:"ParamZ"`
}

func CreateTaskCfg() {
	db.Exec("DROP TABLE IF EXISTS task_cfg")
	db.AutoMigrate(&TaskCfg{})
	//db.Delete(&TaskCfg{})

	fileData, err := ioutil.ReadFile("bin/gamedata/TaskNewCfg.json")
	if err != nil {
		Logger.Fatal("CreateTaskCfg", zap.Error(err))
	}

	var List []*TaskCfg
	jsoniter.Unmarshal(fileData, &List)
	db.CreateInBatches(List, 100)

}
