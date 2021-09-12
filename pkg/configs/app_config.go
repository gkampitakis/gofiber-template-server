package configs

import (
	"os"
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
	isDevelopment := os.Getenv("GO_ENV") != "production"
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	shutdownPeriod, err := time.ParseDuration(os.Getenv("SHUTDOWN_PERIOD"))
	if err != nil {
		shutdownPeriod = 10
	}

	return &AppConfig{
		Host:           "0.0.0.0",
		Port:           port,
		ShutdownPeriod: shutdownPeriod,
		IsDevelopment:  isDevelopment,
	}
}
