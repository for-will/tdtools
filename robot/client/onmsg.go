package client

import (
	"go.uber.org/zap"
	"market/GameMsg"
	"os"
)

var DefaultMsgHandler = map[GameMsg.MsgId]interface{}{
	NetworkConnected:                          OnConnected,
	GameMsg.MsgId_S2C_SyncMainlineTask:        OnSyncMainlineTaskRs,
	GameMsg.MsgId_S2C_AccountCheckRs:          OnAccountCheckRs,
	GameMsg.MsgId_S2C_SyncPlayer:              OnSyncPlayer,
	GameMsg.MsgId_S2C_EquipBackPackRs:       OnCrystalBackPackRs,
	GameMsg.MsgId_S2C_SyncPlayerTalentList:    OnSyncPlayerTalentList,
	GameMsg.MsgId_S2C_HeroTalentInfoRs:        OnHeroTalentInfoRs,
	GameMsg.MsgId_S2C_StoreInfoRs:             OnStoreInfoRs,
	GameMsg.MsgId_S2C_ExploreRs:               OnExploreRs,
	GameMsg.MsgId_S2C_ShowWebViewRs:           OnShowWebViewRs,
	GameMsg.MsgId_S2C_SyncHeroValidTalentPage: OnSyncHeroValidTalentPage,
}

func OnConnected(r *Robot) {
	r.Login(r.Account, r.Password)
	//r.OpeningActivitiesReq()
	//r.Login("11", "11")
}

func OnAccountCheckRs(r *Robot, msg *GameMsg.AccountCheckRs) {
	//r.Explore()
	//r.OpeningActivitiesReq()
}

func OnSyncPlayer(r *Robot, msg *GameMsg.SyncPlayer) {
	r.SyncPlayer = msg
	//r.ValidCardSn = msg.Cards[2].GetSn()
	//r.ValidCardSn = 1627957702
	for _, c := range msg.Cards {
		if c.GetId() == UseCardId {
			r.ValidCardSn = c.GetSn()
			break
		}
	}
	//info := &GameMsg.SyncPlayer{}
	//info.CrystalList = append(info.CrystalList, msg.CrystalList...)
	//msg.CrystalList
	//data, _ := proto.Marshal(info)

	//Log.Debug("OnSyncPlayer",
	//	zap.Int("crystal_cnt", len(info.CrystalList)),
	//	zap.Int("data_size", len(data)))

	//  {"crystal_cnt": 1051, "data_size": 63581}
	//	{"crystal_cnt": 1051, "data_size": 39938}
	//  {"crystal_cnt": 1151, "data_size": 43738}
	//  {"crystal_cnt": 1251, "data_size": 47538}
	//  {"crystal_cnt": 1351, "data_size": 51338}
	//  {"crystal_cnt": 1451, "data_size": 55138}
	//  {"crystal_cnt": 1551, "data_size": 58938}
	//  {"crystal_cnt": 1651, "data_size": 62738}
}

func OnSyncMainlineTaskRs(r *Robot, msg *GameMsg.SyncMainTask) {
	//r.UnlockCard()
	//r.Explore()
	//r.GetCrystalBackpack()
	//r.UpgradePlayerTalent()
	//r.HeroTalentInfo()
	//fmt.Println("OnSyncMainlineTaskRs", JsonString(r), JsonString(msg))
}

func OnCrystalBackPackRs(r *Robot) {
	r.EquipCrystal()
}

func OnSyncPlayerTalentList(r *Robot, msg *GameMsg.SyncPlayerTalentList) {

}

func OnStoreInfoRs(r *Robot, msg *GameMsg.StoreInfoRs) {

	//for _, s := range msg.Stores {
	//	fmt.Println(time.Unix(int64(s.NextFreshTime), 0))
	//}
	//Log.Debug("NextFreshTime")
	r.StorePurchaseReq()
	r.StorePurchaseReq()
}

func OnHeroTalentInfoRs(r *Robot, msg *GameMsg.HeroTalentInfoRs) {
	//r.UnlockHeroTalentPage()
	//r.ModifyHeroTalentPageName()
	//r.SwitchHeroTalentPage()
	//r.UpgradeHeroTalent()
	//r.ResetHeroTalentPage()
	//r.ModifyNickname()
}

func OnLootMissionListRs(r *Robot, msg *GameMsg.LootMissionListRs) {
	zap.L().Debug("OnLootMissionListRs", zap.Any("rsp", msg))
}

func OnExploreRs(r *Robot, msg *GameMsg.ExploreRs) {
	zap.L().Debug("OnExploreRs", zap.Any("rsp", msg))

	for _, card := range msg.Cards {
		r.ValidCardSn = card.Sn
		r.HeroTalentInfo()
	}
}

func TestOnExploreRs(r *Robot, msg *GameMsg.ExploreRs) {

	r.PendingExplore--

	if r.PendingExplore > 0 {
		r.ExploreArea(r.ExploreAreaId, r.ExploreTimes)
	} else {
		os.Exit(0)
	}
}

func OnRobotAutoOverStage(r *Robot, msg *GameMsg.OverStageRs) {
	r.OverStages = r.OverStages[1:]
	if len(r.OverStages) > 0 {
		r.ReqOverStage(r.OverStages[0], r.Stars, r.KillNum)
	} else {
		Log.Info("robot exit 0.")
		os.Exit(0)
	}
}

func OnShowWebViewRs(r *Robot, msg *GameMsg.ShowWebViewRs) {

}

func OnSyncHeroValidTalentPage(r *Robot, msg *GameMsg.SyncHeroValidTalentPage) {

	OnLoginComplete(r)
}

func OnLoginComplete(r *Robot) {
	//r.QuestionnaireReq()
	//r.RewardLootMission()
	//r.OverStage()
	//r.UnlockCard()
	//r.GetLootWall()
	//r.StorePurchaseReq()
	//r.RewardLootMission()
	//r.EquipCrystal()
	//r.PlaceLoot()
	//r.OALoginRewardReq()
	//r.OATaskRewardReq()
	//r.Explore()
	//r.UnlockHeroTalentPage()
	//r.LootMissionList()
	//r.UnlockHeroTalentPage()
	//r.SeasonTaskRewardReq()
	//r.SeasonLvRewardReq()
	//r.DailySignReq()
	//r.SeasonReq()
	r.ReportReq()
}
