package internal

import "go.uber.org/zap"

var Logger *zap.Logger

func init() {
	Logger, _ = zap.NewDevelopment()
}
