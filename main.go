package main

import (
	"encoding/base64"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"log"
	"math/rand"
	"os"
	"robot/GameMsg"
	"robot/client"
	"sync"
	"time"
)

func main() {

	rand.Seed(time.Now().UnixNano())

	if len(os.Args) == 2 {
		runAutoTestRobot(os.Args[1])
	} else {
		newRobot()
	}
	//Benchmark()
}

type operates struct {
	Server       string  `json:"server"`
	UserAccount  string  `json:"user_account"`
	UserPassword string  `json:"user_password"`
	ExploreArea  int32   `json:"explore_area"`
	ExploreType  int32   `json:"explore_type"`
	ExploreCnt   int32   `json:"explore_cnt"`
	Stages       []int32 `json:"stages"`
	KillNum      int32   `json:"kill_num"`
	Stars        int32   `json:"stars"`
}

func runAutoTestRobot(cmd string) {

	var op = &operates{}

	js, _ := base64.StdEncoding.DecodeString(cmd)
	if err := jsoniter.Unmarshal(js, op); err != nil {
		log.Fatal(err)
	}

	handlers := client.DefaultMsgHandler
	handlers[GameMsg.MsgId_S2C_SyncHeroValidTalentPage] = func(r *client.Robot, msg *GameMsg.SyncHeroValidTalentPage) {
		if r.PendingExplore > 0 {
			r.ExploreArea(r.ExploreAreaId, r.ExploreTimes)
		}

		if len(r.OverStages) > 0 {
			r.ReqOverStage(r.OverStages[0], r.Stars, r.KillNum)
		}
	}
	handlers[GameMsg.MsgId_S2C_ExploreRs] = client.TestOnExploreRs
	handlers[GameMsg.MsgId_S2C_OverStageRs] = client.OnRobotAutoOverStage

	r := &client.Robot{
		ServerAddr:     op.Server,
		MsgHandler:     handlers,
		Account:        op.UserAccount,
		Password:       op.UserPassword,
		ExploreAreaId:  op.ExploreArea,
		ExploreTimes:   op.ExploreType,
		PendingExplore: op.ExploreCnt,
		OverStages:     op.Stages,
		KillNum:        op.KillNum,
		Stars:          op.Stars,
	}

	r.Start()
	client.Log.Sync()
}

func newRobot() {
	r := &client.Robot{
		ServerAddr: client.ServerAddr,
		MsgHandler: client.DefaultMsgHandler,
		Account:    client.RobotAccount,
		Password:   client.RobotPassword,
	}

	r.Start()
	client.Log.Sync()
}

func Benchmark() {
	wg := &sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		r := &client.Robot{
			MsgHandler: client.DefaultMsgHandler,
			Account:    fmt.Sprintf("Test%d", i),
			Password:   "123456",
		}

		wg.Add(1)
		go func() {
			r.Start()
			wg.Done()
		}()
	}
	wg.Wait()
}

func f() (r int) {
	t := 5
	defer func() {
		r = r + t
	}()
	return 1
}
