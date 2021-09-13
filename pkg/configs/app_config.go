package configs

import (
	"log"
	"os"
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
	isDevelopment := os.Getenv("GO_ENV") != "production"

	if isDevelopment {
		err := godotenv.Load()
		if err != nil {
			log.Println("[App Config]", err.Error())
		}
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	host := "0.0.0.0"
	if isDevelopment {
		host = "localhost"
	}

	shutdownPeriod, err := time.ParseDuration(os.Getenv("SHUTDOWN_PERIOD"))
	if err != nil {
		shutdownPeriod = 10
	}

	return &AppConfig{
		Host:           host,
		Port:           port,
		ShutdownPeriod: shutdownPeriod,
		IsDevelopment:  isDevelopment,
	}
}
