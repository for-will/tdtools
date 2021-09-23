package client

import (
	"fmt"
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
		Account:  &account,
		Password: &password,
	})
}

func (r *Robot) Explore() {
	var area int32 = 5
	var times GameMsg.ExploreTimes = GameMsg.ExploreTimes_Ten
	r.SendMsg(&GameMsg.Explore{
		Area:  &area,
		Times: &times,
	})
}

func (r *Robot) UnlockCard() {

	var cardId int32 = 1004
	r.SendMsg(&GameMsg.CardUnLock{CardId: &cardId})
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
		Name: NewString("vanish2"),
	})
}

func (r *Robot) ModifyHeadImage() {
	r.SendMsg(&GameMsg.ModifyPlayerIcon{
		Icon: NewInt32(11),
	})
}

func OnConnected(r *Robot) {
	r.Login(RobotAccount, RobotPassword)
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
}

func OnHeroTalentInfoRs(r *Robot, msg *GameMsg.HeroTalentInfoRs) {
	//r.UnlockHeroTalentPage()
	//r.ModifyHeroTalentPageName()
	//r.SwitchHeroTalentPage()
	//r.UpgradeHeroTalent()
	//r.ResetHeroTalentPage()
	//r.ModifyNickname()
}

func NewInt32(v int32) *int32 {
	return &v
}

func NewInt64(v int64) *int64 {
	return &v
}

func NewString(v string) *string {
	return &v
}
