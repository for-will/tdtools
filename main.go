package main

import (
	"robot/GameMsg"
	"robot/client"
)

func main() {

	//client.UpdateDb()
	//fmt.Println(f())
	//
	newRobot()
}

func newRobot() {
	r := &client.Robot{
		MsgHandler: map[GameMsg.MsgId]interface{}{
			client.NetworkConnected:                client.OnConnected,
			GameMsg.MsgId_S2C_SyncMainlineTask:     client.OnSyncMainlineTaskRs,
			GameMsg.MsgId_S2C_AccountCheckRs:       client.OnAccountCheckRs,
			GameMsg.MsgId_S2C_SyncPlayer:           client.OnSyncPlayer,
			GameMsg.MsgId_S2C_CrystalBackPackRs:    client.OnCrystalBackPackRs,
			GameMsg.MsgId_S2C_SyncPlayerTalentList: client.OnSyncPlayerTalentList,
			GameMsg.MsgId_S2C_HeroTalentInfoRs:     client.OnHeroTalentInfoRs,
		},
	}

	r.Start()
	client.Log.Sync()
}

func f() (r int) {
	t := 5
	defer func() {
		r = r + t
	}()
	return 1
}
