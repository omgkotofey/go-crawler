package logger

import (
	"go.uber.org/zap"
)

var (
	log *zap.Logger
)

func InitLogger(isProduction bool) *zap.Logger {
	if log != nil {
		return log
	}

	var err error
	if isProduction {
		log, err = zap.NewProduction()
	} else {
		log, err = zap.NewDevelopment()
	}

	if err != nil {
		panic(err)
	}

	return log
}
