package configs

import (
	"log"
	"strconv"
	"time"
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

	port := GetEnv("APP_PORT", "8080")
	host := "0.0.0.0"
	if isDevelopment {
		host = "localhost"
	}
	shutdownPeriod, err := strconv.Atoi(GetEnv("APP_SHUTDOWN_PERIOD", "10"))
	if err != nil {
		shutdownPeriod = 10
		log.Printf("[APP_SHUTDOWN_PERIOD] incorrect value, defaulting to %d", shutdownPeriod)
	}

	return &AppConfig{
		Host:           host,
		Port:           port,
		ShutdownPeriod: time.Duration(shutdownPeriod),
		IsDevelopment:  isDevelopment,
	}
}
