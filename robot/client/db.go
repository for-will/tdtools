package client

import (
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func InitDb(dsn string) *gorm.DB {

	l := logger.New(log.New(os.Stderr, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             1 * time.Second,
		Colorful:                  false,
		IgnoreRecordNotFoundError: false,
		LogLevel:                  logger.Info,
	})

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: l,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		Log.Fatal("open db", zap.Error(err))
	}
	if sqlDb, err := db.DB(); err != nil {
		sqlDb.SetMaxOpenConns(3)
	}

	return db
}

type Player struct {
	Id       int
	NickName string
	Lv       int32
	Exp      int32
	Gold     int32
	Diamond  int32
}

func dbJobs(db *gorm.DB, upt *ModifyDb) {

	player := &Player{}
	if err := db.First(player, "nick_name=?", upt.NickName).Error; err != nil {
		Log.Fatal("query player failed", zap.Error(err))
	}

	Log.Debug("player_info", zap.Any("player", player))

	for _, card := range upt.Cards {
		if err := db.Exec("UPDATE card SET talent_points = ? WHERE player_id=? AND id = ?",
			card.TalentPoint, player.Id, card.CardId,
		).Error; err != nil {
			Log.Error("db error", zap.Error(err))
		}
	}
}

func readUpdate() *ModifyDb {
	update := &ModifyDb{}

	b, err := ioutil.ReadFile("config.json")
	if err != nil {
		Log.Fatal("read config failed", zap.Error(err))
	}
	if err := jsoniter.Unmarshal(b, update); err != nil {
		Log.Fatal("unmarshal config", zap.Error(err))
	}

	Log.Debug("config", zap.Any("update", update))
	return update
}

func UpdateDb() {

	defer func() {
		if err := recover(); err != nil {
			Log.Error("recover", zap.Any("err", err))
		}
		var press [1]byte
		os.Stdin.Read(press[:])

		//signalCh := make(chan os.Signal)
		//signal.Notify(signalCh, os.Interrupt)
		//<-signalCh
	}()
	upt := readUpdate()
	db := InitDb(upt.Dsn)
	dbJobs(db, upt)
}

type CardUpdate struct {
	CardId      int32 `json:"card_id"`
	TalentPoint int32 `json:"talent_point"`
}

type ModifyDb struct {
	Dsn      string        `json:"dsn"`
	NickName string        `json:"nick_name"`
	Cards    []*CardUpdate `json:"cards"`
}
