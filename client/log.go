package client

import (
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"log"
	"robot/GameMsg"
	"robot/js"
)

var Log *zap.Logger

func init() {

	//zap.Development()
	l, _ := zap.NewDevelopment(zap.Development())

	//l, _ := zap.NewProduction()
	//l.Info("hello info")
	//zap.New()
	Log = l
}

const (
	ColorRcvErr = "\u001B[31m"
	ColorRcv    = "\u001B[0;36m"
	ColorNotify = "\u001B[33m"
	ColorSnd    = "\u001B[7;32m"
)

func LogRcvMsg(id GameMsg.MsgId, msg proto.Message) {
	switch id {
	case GameMsg.MsgId_S2C_SyncPlayer,
		GameMsg.MsgId_S2C_UpdateInfo,
		GameMsg.MsgId_S2C_SyncMainlineTask,
		GameMsg.MsgId_S2C_SyncPlayerTalentList,
		GameMsg.MsgId_S2C_ShowWebViewRs,
		GameMsg.MsgId_S2C_TaskInfo,
		GameMsg.MsgId_S2C_PlayerOffline,
		GameMsg.MsgId_S2C_SyncHeroValidTalentPage:
		doLogMessage(id, msg, ColorNotify, "<")
	default:
		doLogMessage(id, msg, ColorRcv, "<")

	}

}

func LogErrMsg(id GameMsg.MsgId, msg proto.Message) {
	doLogMessage(id, msg, ColorRcvErr, "<")
}

func LogSndMsg(id GameMsg.MsgId, msg proto.Message) {
	doLogMessage(id, msg, ColorSnd, ">")
}

func doLogMessage(id GameMsg.MsgId, msg proto.Message, a string, tag string) {
	log.Printf("%s%s %-30s| %v\u001B[0m\n", a, tag, id, js.PbMinifyJson(msg))
}
