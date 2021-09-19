package configs

import (
	"log"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	_              struct{}
	Host           string
	ShutdownPeriod time.Duration
	Port           string
	IsDevelopment  bool
}

func NewAppConfig() *AppConfig {
	isDevelopment := GetEnv("GO_ENV") != "production"

	if isDevelopment {
		err := godotenv.Load()
		if err != nil {
			log.Println("[App Config]", err.Error())
		}
	}

	port := GetEnv("APP_PORT", "8080")
	host := "0.0.0.0"
	if isDevelopment {
		host = "localhost"
	}
	shutdownPeriod, err := strconv.Atoi(GetEnv("SHUTDOWN_PERIOD", "10"))
	if err != nil {
		shutdownPeriod = 10
		log.Printf("[SHUTDOWN_PERIOD] incorrect value, defaulting to %d", shutdownPeriod)
	}

	return &AppConfig{
		Host:           host,
		Port:           port,
		ShutdownPeriod: time.Duration(shutdownPeriod),
		IsDevelopment:  isDevelopment,
	}
}
