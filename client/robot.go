package client

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"math/rand"
	"reflect"
	"robot/GameMsg"
	"robot/js"
	"strconv"
	"time"
)

type Robot struct {
	*Client
	ServerAddr  string
	MsgHandler  map[GameMsg.MsgId]interface{}
	ValidCardSn int32
	SyncPlayer  *GameMsg.SyncPlayer
	Account     string
	Password    string

	PendingExplore int32
	ExploreAreaId  int32
	ExploreTimes   int32

	OverStages []int32
	KillNum    int32
	Stars      int32
}

const (
	NetworkConnected GameMsg.MsgId = -1 // 连接成功
)

func (r *Robot) Start() {

	r.Client = &Client{
		ServerTCP: r.ServerAddr,
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

	r.SendMsg(&GameMsg.Explore{
		Area:  6,
		Times: GameMsg.ExploreTimes_Ten,
	})
}

func (r *Robot) ExploreArea(area int32, times int32) {
	r.SendMsg(&GameMsg.Explore{
		Area:  area,
		Times: GameMsg.ExploreTimes(times),
	})
}

func (r *Robot) UnlockCard() {

	//var cardId int32 = 1005
	r.SendMsg(&GameMsg.CardUnLock{CardId: NewInt32(1022)})
}

func (r *Robot) HeroQualityUp() {
	r.SendMsg(&GameMsg.HeroQualityUp{HeroSn: 1639376076})
}

func (r *Robot) GetCrystalBackpack() {

	r.SendMsg(&GameMsg.CrystalBackPack{})
}

func (r *Robot) EquipCrystal() {
	r.SendMsg(&GameMsg.EquipCrystal{
		Sn:     NewInt32(11239),
		HeroSn: NewInt32(1639376076),
		Slot:   NewInt32(1),
	})
}

func (r *Robot) UpgradePlayerTalent() {
	r.SendMsg(&GameMsg.UpgradePlayerTalent{Id: 5})
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
	r.SendMsg(&GameMsg.RewardLootMission{Id: NewInt32(30201)})
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
		LootMissionId: NewInt32(30201),
		PsX:           NewInt32(3),
		PsY:           NewInt32(40),
	})
	req.List = append(req.List, &GameMsg.LootItem{
		LootMissionId: NewInt32(30103),
		PsX:           NewInt32(30),
		PsY:           NewInt32(40),
	})
	r.SendMsg(req)
}

func (r *Robot) ModifyNickname() {
	r.SendMsg(&GameMsg.ModifyPlayerName{
		Name: NewString("aa" + strconv.Itoa(rand.Intn(100000))),
	})
}

func (r *Robot) ModifyHeadImage() {
	r.SendMsg(&GameMsg.ModifyPlayerIcon{
		Icon: NewInt32(12),
	})
}

func (r *Robot) InitPlayerName() {
	r.SendMsg(&GameMsg.InitPlayerName{
		Name: "N1638240472",
	})
}

func (r *Robot) OverStage() {
	r.SendMsg(&GameMsg.OverStage{
		StageId: 128,
		//IsWin:   true,
		//Param:   3,
		//KillNum: 10,
		//EnemyList: nil,
	})
}

func (r *Robot) ReqOverStage(stageId int32, stars int32, killNum int32) {
	r.SendMsg(&GameMsg.OverStage{
		StageId: stageId,
		IsWin:   true,
		Param:   stars,
		KillNum: killNum,
	})
}

func (r *Robot) StoreInfoReq() {
	r.SendMsg(&GameMsg.StoreInfoReq{})
}

func (r *Robot) StorePurchaseReq() {
	r.SendMsg(&GameMsg.StorePurchaseReq{
		Id:  5,
		Cnt: 1,
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

func (r *Robot) GetRewardStage() {
	r.SendMsg(&GameMsg.GetRewardStage{
		Id:        3,
		ChapterId: 1,
	})
}

func (r *Robot) GetTaskReward() {
	r.SendMsg(&GameMsg.GetTaskReward{
		Id: []int32{10101},
	})
}

func (r *Robot) UseItem() {
	r.SendMsg(&GameMsg.UseItem{
		Id:  1001,
		Cnt: 1,
	})
}

func (r *Robot) QuestionnaireReq() {
	r.SendMsg(&GameMsg.QuestionnaireReq{
		Action: "",
		Data: js.MinifyJson(&struct {
			Id string `json:"surveyId"`
		}{
			Id: "cev4f3",
		}),
	})
}

func (r *Robot) SeasonReq() {
	r.SendMsg(&GameMsg.SeasonReq{})
}
