package configs

import "time"

type HealthcheckConfig struct {
	TimeoutPeriod  time.Duration
	TimeoutEnabled bool
}

// FIXME: values from env

func NewHealthcheckConfig() *HealthcheckConfig {
	return &HealthcheckConfig{
		TimeoutPeriod:  5,
		TimeoutEnabled: false,
	}
}
