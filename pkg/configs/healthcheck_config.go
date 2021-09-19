package configs

import (
	"log"
	"strconv"
	"time"
)

type HealthcheckConfig struct {
	TimeoutPeriod  time.Duration
	TimeoutEnabled bool
}

func NewHealthcheckConfig() *HealthcheckConfig {
	timeoutPeriod, err := strconv.Atoi(GetEnv("HC_TIMEOUT_PERIOD", "5"))
	if err != nil {
		timeoutPeriod = 5
		log.Printf("[HC_TIMEOUT_PERIOD] incorrect value, defaulting to %d", timeoutPeriod)
	}
	timeoutEnabled := GetEnv("HC_TIMEOUT_ENABLED", "true") == "true"

	return &HealthcheckConfig{
		TimeoutPeriod:  time.Duration(timeoutPeriod),
		TimeoutEnabled: timeoutEnabled,
	}
}
