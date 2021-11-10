package configs

import (
	"log"
	"strconv"
	"time"

	"github.com/gkampitakis/fiber-modules/gracefulshutdown"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	_                      struct{}
	Addr                   string
	Service                string
	IsDevelopment          bool
	GracefulshutdownConfig func(fns ...func() error) gracefulshutdown.Config
}

func New() *AppConfig {
	isDevelopment := GetEnv("GO_ENV") != "production"
	host := "0.0.0.0"

	if isDevelopment {
		host = "localhost"

		err := godotenv.Load()
		if err != nil {
			log.Printf("[App Config] " + err.Error())
		}
	}

	port := GetEnv("APP_PORT", "8080")

	gracePeriod, err := strconv.Atoi(GetEnv("GRACE_PERIOD", "15"))
	if err != nil {
		gracePeriod = 15
	}

	return &AppConfig{
		Addr:          host + ":" + port,
		IsDevelopment: isDevelopment,
		Service:       GetEnv("HC_APP_NAME", "gofiber-template"),
		GracefulshutdownConfig: func(fns ...func() error) gracefulshutdown.Config {
			return gracefulshutdown.Config{
				Period:      time.Duration(gracePeriod),
				Enabled:     !isDevelopment,
				ShutdownFns: fns,
			}
		},
	}
}
