package client

import "go.uber.org/zap"

var Log *zap.Logger

func init() {

	//zap.Development()
	l, _ := zap.NewDevelopment(zap.Development())

	//l, _ := zap.NewProduction()
	//l.Info("hello info")
	//zap.New()
	Log = l
}
