package client

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"reflect"
	"robot/GameMsg"
	"time"
)

type Robot struct {
	*Client
	MsgHandler  map[GameMsg.MsgId]interface{}
	ValidCardSn int32
	SyncPlayer  *GameMsg.SyncPlayer
	Account     string
	Password    string
}

const (
	NetworkConnected GameMsg.MsgId = -1 // 连接成功
)

func (r *Robot) Start() {

	r.Client = &Client{
		msgHandler: func(id GameMsg.MsgId, message proto.Message) {
			r.OnMessage(id, message)
		},
	}
	r.Init()
	r.Run()
}

func (r *Robot) OnMessage(id GameMsg.MsgId, msg interface{}) {

	h, ok := r.MsgHandler[id]
	if !ok {
		return
	}

	typ := reflect.TypeOf(h)
	val := reflect.ValueOf(h)

	var in []reflect.Value
	for i := 0; i < typ.NumIn(); i++ {
		switch typ.In(i) {
		case reflect.TypeOf(r):
			in = append(in, reflect.ValueOf(r))

		case reflect.TypeOf(msg):
			in = append(in, reflect.ValueOf(msg))

		default:
			fmt.Printf("can't provie func %d parameter: %v", i, typ.In(i))
			return
		}
	}
	val.Call(in)
}

func (r *Robot) Login(account string, password string) {
	r.SendMsg(&GameMsg.AccountCheck{
		Account:  NewString(account),
		Password: NewString(password),
	})
}

func (r *Robot) Explore() {
	var area int32 = 5
	r.SendMsg(&GameMsg.Explore{
		Area:  NewInt32(area),
		Times: GameMsg.ExploreTimes_Ten,
	})
}

func (r *Robot) UnlockCard() {

	//var cardId int32 = 1005
	r.SendMsg(&GameMsg.CardUnLock{CardId: NewInt32(1005)})
}

func (r *Robot) HeroQualityUp() {
	r.SendMsg(&GameMsg.HeroQualityUp{HeroSn: NewInt32(1632454723)})
}

func (r *Robot) GetCrystalBackpack() {

	r.SendMsg(&GameMsg.CrystalBackPack{})
}

func (r *Robot) EquipCrystal() {
	r.SendMsg(&GameMsg.EquipCrystal{
		Sn:     NewInt32(2),
		HeroSn: NewInt32(1624860736),
		Slot:   NewInt32(0),
	})
}

func (r *Robot) UpgradePlayerTalent() {
	r.SendMsg(&GameMsg.UpgradePlayerTalent{Id: NewInt32(1)})
}

func (r *Robot) HeroTalentInfo() {
	r.SendMsg(&GameMsg.HeroTalentInfo{CardSn: NewInt32(r.ValidCardSn)})
}

func (r *Robot) UnlockHeroTalentPage() {
	r.SendMsg(&GameMsg.UnlockHeroTalentPage{
		CardSn: NewInt32(r.ValidCardSn),
		Page:   NewInt32(3),
	})
}

func (r *Robot) ModifyHeroTalentPageName() {
	r.SendMsg(&GameMsg.ModifyHeroTalentPageName{
		CardSn: NewInt32(r.ValidCardSn),
		Page:   NewInt32(2),
		Name:   NewString(ModifyHeroTalentPageName),
	})
}

func (r *Robot) SwitchHeroTalentPage() {
	r.SendMsg(&GameMsg.SwitchHeroTalentPage{
		CardSn: NewInt32(r.ValidCardSn),
		Page:   NewInt32(2),
	})
}

func (r *Robot) UpgradeHeroTalent() {
	r.SendMsg(&GameMsg.UpgradeHeroTalent{
		CardSn:   NewInt32(r.ValidCardSn),
		Page:     NewInt32(1),
		TalentId: NewInt32(UpgradeHeroTalent),
	})
}

func (r *Robot) ResetHeroTalentPage() {
	r.SendMsg(&GameMsg.ResetHeroTalentPage{
		CardSn: NewInt32(r.ValidCardSn),
		Page:   NewInt32(1),
	})
}

func (r *Robot) LogInstall() {
	r.SendMsg(&GameMsg.LogInstallReq{
		InstallTime: NewInt64(time.Now().Unix()),
		Ip:          NewString("127.0.0.1"),
		DeviceModel: NewString("iphone"),
		OsName:      NewString("ios6"),
		OsVer:       NewString("1.2.3"),
		MacAddr:     NewString("[::1]"),
		Udid:        NewString("123456"),
		AppChannel:  NewString("appstore"),
		AppVer:      NewString("2.3.1"),
	})
}

func (r *Robot) LootMissionList() {
	r.SendMsg(&GameMsg.LootMissionList{})
}

func (r *Robot) RewardLootMission() {
	r.SendMsg(&GameMsg.RewardLootMission{Id: NewInt32(100103)})
}

func (r *Robot) GetLootWall() {
	r.SendMsg(&GameMsg.GetLootWall{})
}

func (r *Robot) ClearLootWall() {
	r.SendMsg(&GameMsg.ClearLootWall{})
}

func (r *Robot) PlaceLoot() {
	req := &GameMsg.PlaceLoot{}
	req.List = append(req.List, &GameMsg.LootItem{
		LootMissionId: NewInt32(100103),
		PsX:           NewInt32(3),
		PsY:           NewInt32(40),
	})
	req.List = append(req.List, &GameMsg.LootItem{
		LootMissionId: NewInt32(100102),
		PsX:           NewInt32(30),
		PsY:           NewInt32(40),
	})
	r.SendMsg(req)
}

func (r *Robot) ModifyNickname() {
	r.SendMsg(&GameMsg.ModifyPlayerName{
		Name: NewString("x-1"),
	})
}

func (r *Robot) ModifyHeadImage() {
	r.SendMsg(&GameMsg.ModifyPlayerIcon{
		Icon: NewInt32(12),
	})
}

func (r *Robot) InitPlayerName() {
	r.SendMsg(&GameMsg.InitPlayerName{
		Name: "x3",
	})
}

func (r *Robot) OverStage() {
	r.SendMsg(&GameMsg.OverStage{
		StageId: 1010001,
		IsWin:   true,
		Param:   0,
		KillNum: 11,
		//EnemyList: nil,
	})
}

func (r *Robot) StoreInfoReq() {
	r.SendMsg(&GameMsg.StoreInfoReq{})
}

func (r *Robot) StorePurchaseReq() {
	r.SendMsg(&GameMsg.StorePurchaseReq{
		Id:  6,
		Cnt: 2,
	})
}

func (r *Robot) OpeningActivitiesReq() {
	r.SendMsg(&GameMsg.OpeningActivitiesReq{})
}

func (r *Robot) OALoginRewardReq() {
	r.SendMsg(&GameMsg.OALoginRewardReq{
		ActivityId: NewInt32(1),
	})
}

func (r *Robot) OATaskRewardReq() {
	r.SendMsg(&GameMsg.OATaskRewardReq{
		Id: 13,
	})
}

func OnConnected(r *Robot) {
	r.Login(r.Account, r.Password)
	//r.Login("11", "11")
}

func OnAccountCheckRs(r *Robot, msg *GameMsg.AccountCheckRs) {
	//r.Explore()
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
	info := &GameMsg.SyncPlayer{}
	info.CrystalList = append(info.CrystalList, msg.CrystalList...)
	//msg.CrystalList
	data, _ := proto.Marshal(info)

	Log.Debug("OnSyncPlayer",
		zap.Int("crystal_cnt", len(info.CrystalList)),
		zap.Int("data_size", len(data)))
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

func OnCrystalBackPackRs(r *Robot, msg *GameMsg.CrystalBackPackRs) {
	r.EquipCrystal()
}

func OnSyncPlayerTalentList(r *Robot, msg *GameMsg.SyncPlayerTalentList) {
	//r.HeroTalentInfo()
	//r.UpgradePlayerTalent()
	//r.Explore()
	//r.LogInstall()
	//r.UnlockCard()
	//r.LootMissionList()
	//r.PlaceLoot()
	//r.GetLootWall()
	//r.ClearLootWall()
	//r.GetLootWall()
	r.ModifyNickname()
	//r.ModifyHeadImage()
	//r.HeroQualityUp()
	//r.InitPlayerName()
	//fmt.Println(r.Account, JsonString(msg))
	//r.OverStage()
	//r.StoreInfoReq()
	//r.OpeningActivitiesReq()
	//r.OALoginRewardReq()
	//r.OATaskRewardReq()
	//r.OpeningActivitiesReq()
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

//
//func NewInt32(v int32) *int32 {
//	return &v
//}
//
//func NewInt64(v int64) *int64 {
//	return &v
//}
//
//func NewString(v string) *string {
//	return &v
//}

func NewInt32(v int32) int32 {
	return v
}

func NewInt64(v int64) int64 {
	return v
}

func NewString(v string) string {
	return v
}

var DefaultMsgHandler = map[GameMsg.MsgId]interface{}{
	NetworkConnected:                       OnConnected,
	GameMsg.MsgId_S2C_SyncMainlineTask:     OnSyncMainlineTaskRs,
	GameMsg.MsgId_S2C_AccountCheckRs:       OnAccountCheckRs,
	GameMsg.MsgId_S2C_SyncPlayer:           OnSyncPlayer,
	GameMsg.MsgId_S2C_CrystalBackPackRs:    OnCrystalBackPackRs,
	GameMsg.MsgId_S2C_SyncPlayerTalentList: OnSyncPlayerTalentList,
	GameMsg.MsgId_S2C_HeroTalentInfoRs:     OnHeroTalentInfoRs,
	GameMsg.MsgId_S2C_StoreInfoRs:          OnStoreInfoRs,
	GameMsg.MsgId_S2C_ExploreRs:            OnExploreRs,
}
