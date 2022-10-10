package internal

import (
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

var db *gorm.DB

func InitDb(dsn string) *gorm.DB {

	l := logger.New(log.New(os.Stderr, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             1 * time.Second,
		Colorful:                  false,
		IgnoreRecordNotFoundError: false,
		LogLevel:                  logger.Info,
	})

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: l,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		Logger.Fatal("open db", zap.Error(err))
	}
	if sqlDb, err := database.DB(); err != nil {
		sqlDb.SetMaxOpenConns(3)
	}

	return database
}

func init() {
	db = InitDb("game:game123@tcp(127.0.0.1:3306)/game?charset=utf8mb4&parseTime=True&loc=Local")
}
