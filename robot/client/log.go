package client

import (
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"log"
	"market/GameMsg"
	"market/js"
	"os"
)

var Log *zap.Logger

func init() {

	//zap.Development()
	l, _ := zap.NewDevelopment(zap.Development())

	//l, _ := zap.NewProduction()
	//l.Info("hello info")
	//zap.New()
	Log = l

	log.SetFlags(0)
	log.SetOutput(os.Stdout)
}

const (
	ColorRcvErr = "\u001B[1;31m"
	ColorRcv    = "\u001B[2;36m"
	ColorNotify = "\u001b[38;5;130m" //"\u001B[3;33m"
	ColorSnd    = "\u001B[1;32m"
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
		GameMsg.MsgId_S2C_Strength,
		GameMsg.MsgId_S2C_SyncHeroValidTalentPage:
		doLogMessage(id, msg, ColorNotify, "☀\t")
	default:
		doLogMessage(id, msg, ColorRcv, "▼\t")

	}

}

func LogErrMsg(id GameMsg.MsgId, msg proto.Message) {
	doLogMessage(id, msg, ColorRcvErr, "☢\t")
}

func LogSndMsg(id GameMsg.MsgId, msg proto.Message) {
	doLogMessage(id, msg, ColorSnd, "▲\t")
}

func doLogMessage(id GameMsg.MsgId, msg proto.Message, a string, tag string) {
	//log.SetFlags(log.LstdFlags)
	//tag = "■"
	tag = a + "\u001B[27m" + tag + "\u001B[0m"
	log.Printf("%s %s%-27s\u001B[0m %s\n", tag, a, id, js.PbMinifyJson(msg))
}

//◀▶◁▷☀■□☢☠
