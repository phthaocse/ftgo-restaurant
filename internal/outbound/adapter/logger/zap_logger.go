package logger

import (
	"go.uber.org/zap"
)

var ZapLogger *zap.SugaredLogger

func init() {
	devLogger, err := zap.NewDevelopment()
	defer devLogger.Sync()
	if err != nil {
		panic("Can not construct logger")
	}
	ZapLogger = devLogger.Sugar()
}
