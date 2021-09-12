package utils

import (
	"log"
	"os"

	"go.uber.org/zap"
)

var Logger *zap.Logger
var err error

func init() {
	isProduction := os.Getenv("GO_ENV") == "production"
	if isProduction {
		Logger, err = zap.NewProduction()

		if err != nil {
			log.Panic(err)
		}
		return
	}

	Logger, err = zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}
}
