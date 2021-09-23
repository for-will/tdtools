package client

import (
	"encoding/binary"
	jsoniter "github.com/json-iterator/go"
	"reflect"
	"robot/GameMsg"
)

var (
	messageId                    = map[reflect.Type]GameMsg.MsgId{}
	messageType                  = map[GameMsg.MsgId]reflect.Type{}
	ByteOrder   binary.ByteOrder = binary.LittleEndian
)

func init() {
	var msgId = map[GameMsg.MsgId]interface{}{
		GameMsg.MsgId_C2S_AccountCheck:               &GameMsg.AccountCheck{},
		GameMsg.MsgId_S2C_AccountCheckRs:             &GameMsg.AccountCheckRs{},
		GameMsg.MsgId_S2C_TaskInfo:                   &GameMsg.TaskInfo{},
		GameMsg.MsgId_C2S_HeartBeat:                  &GameMsg.HeartBeat{},
		GameMsg.MsgId_S2C_HeartBeatRs:                &GameMsg.HeartBeatRs{},
		GameMsg.MsgId_C2S_Explore:                    &GameMsg.Explore{},
		GameMsg.MsgId_S2C_ExploreRs:                  &GameMsg.ExploreRs{},
		GameMsg.MsgId_S2C_SyncPlayer:                 &GameMsg.SyncPlayer{},
		GameMsg.MsgId_S2C_SyncMainlineTask:           &GameMsg.SyncMainTask{},
		GameMsg.MsgId_S2C_UpdateInfo:                 &GameMsg.UpdateInfo{},
		GameMsg.MsgId_C2S_CardUnLock:                 &GameMsg.CardUnLock{},
		GameMsg.MsgId_S2C_CardUnLockRs:               &GameMsg.CardUnLockRs{},
		GameMsg.MsgId_C2S_CrystalBackPack:            &GameMsg.CrystalBackPack{},
		GameMsg.MsgId_S2C_CrystalBackPackRs:          &GameMsg.CrystalBackPackRs{},
		GameMsg.MsgId_C2S_EquipCrystal:               &GameMsg.EquipCrystal{},
		GameMsg.MsgId_S2C_EquipCrystalRs:             &GameMsg.EquipCrystalRs{},
		GameMsg.MsgId_S2C_MainlineTaskInfo:           &GameMsg.TaskInfo{},
		GameMsg.MsgId_C2S_UpgradePlayerTalent:        &GameMsg.UpgradePlayerTalent{},
		GameMsg.MsgId_S2C_UpgradePlayerTalentRs:      &GameMsg.UpgradePlayerTalentRs{},
		GameMsg.MsgId_S2C_SyncPlayerTalentList:       &GameMsg.SyncPlayerTalentList{},
		GameMsg.MsgId_C2S_HeroTalentInfo:             &GameMsg.HeroTalentInfo{},
		GameMsg.MsgId_S2C_HeroTalentInfoRs:           &GameMsg.HeroTalentInfoRs{},
		GameMsg.MsgId_C2S_ResetHeroTalentPage:        &GameMsg.ResetHeroTalentPage{},
		GameMsg.MsgId_S2C_ResetHeroTalentPageRs:      &GameMsg.ResetHeroTalentPageRs{},
		GameMsg.MsgId_C2S_UnlockHeroTalentPage:       &GameMsg.UnlockHeroTalentPage{},
		GameMsg.MsgId_S2C_UnlockHeroTalentPageRs:     &GameMsg.UnlockHeroTalentPageRs{},
		GameMsg.MsgId_C2S_ModifyHeroTalentPageName:   &GameMsg.ModifyHeroTalentPageName{},
		GameMsg.MsgId_S2C_ModifyHeroTalentPageNameRs: &GameMsg.ModifyHeroTalentPageNameRs{},
		GameMsg.MsgId_C2S_SwitchHeroTalentPage:       &GameMsg.SwitchHeroTalentPage{},
		GameMsg.MsgId_S2C_SwitchHeroTalentPageRs:     &GameMsg.SwitchHeroTalentPageRs{},
		GameMsg.MsgId_C2S_UpgradeHeroTalent:          &GameMsg.UpgradeHeroTalent{},
		GameMsg.MsgId_S2C_UpgradeHeroTalentRs:        &GameMsg.UpgradeHeroTalentRs{},
		GameMsg.MsgId_S2C_SyncHeroValidTalentPage:    &GameMsg.SyncHeroValidTalentPage{},
		GameMsg.MsgId_S2C_ChatPrivateListRs:          &GameMsg.ChatPrivateListRs{},
		GameMsg.MsgId_C2S_LogInstallReq:              &GameMsg.LogInstallReq{},
		GameMsg.MsgId_S2C_LogInstallRsp:              &GameMsg.LogInstallRsp{},
		GameMsg.MsgId_C2S_LootMissionList:            &GameMsg.LootMissionList{},
		GameMsg.MsgId_S2C_LootMissionListRs:          &GameMsg.LootMissionListRs{},
		GameMsg.MsgId_C2S_RewardLootMission:          &GameMsg.RewardLootMission{},
		GameMsg.MsgId_S2C_RewardLootMissionRs:        &GameMsg.RewardLootMissionRs{},
		GameMsg.MsgId_C2S_GetLootWall:                &GameMsg.GetLootWall{},
		GameMsg.MsgId_S2C_GetLootWallRs:              &GameMsg.GetLootWallRs{},
		GameMsg.MsgId_C2S_PlaceLoot:                  &GameMsg.PlaceLoot{},
		GameMsg.MsgId_S2C_PlaceLootRs:                &GameMsg.PlaceLootRs{},
		GameMsg.MsgId_C2S_ClearLootWall:              &GameMsg.ClearLootWall{},
		GameMsg.MsgId_S2C_ClearLootWallRs:            &GameMsg.ClearLootWallRs{},
		GameMsg.MsgId_C2S_ModifyPlayerIcon:           &GameMsg.ModifyPlayerIcon{},
		GameMsg.MsgId_S2C_ModifyPlayerIconRs:         &GameMsg.ModifyPlayerIconRs{},
		GameMsg.MsgId_C2S_ModifyPlayerName:           &GameMsg.ModifyPlayerName{},
		GameMsg.MsgId_S2C_ModifyPlayerNameRs:         &GameMsg.ModifyPlayerNameRs{},
		GameMsg.MsgId_S2C_PlayerOffline:              &GameMsg.PlayerOffline{},
	}

	for id, v := range msgId {
		typ := reflect.TypeOf(v)
		messageType[id] = typ
		messageId[typ] = id
	}
}

func JsonString(v interface{}) string {
	s, _ := jsoniter.MarshalToString(v)
	return s
}

func NewRetCode(code GameMsg.ReturnCode) *GameMsg.ReturnCode {
	p := new(GameMsg.ReturnCode)
	*p = code
	return p
}
