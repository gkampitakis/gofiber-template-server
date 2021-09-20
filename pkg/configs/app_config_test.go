package configs

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func clearEnv() {
	os.Unsetenv("GO_ENV")
	os.Unsetenv("APP_PORT")
	os.Unsetenv("APP_SHUTDOWN_PERIOD")
	os.Unsetenv("HC_TIMEOUT_PERIOD")
	os.Unsetenv("HC_TIMEOUT_ENABLED")
}

func TestAppConfig(t *testing.T) {
	t.Run("should return default values", func(t *testing.T) {
		defer clearEnv()
		clearEnv()

		config := NewAppConfig()
		assert.Equal(t, true, config.IsDevelopment)
		assert.Equal(t, "localhost", config.Host)
		assert.Equal(t, "8080", config.Port)
		assert.Equal(t, time.Duration(10), config.ShutdownPeriod)
	})

	t.Run("should return setted values and host as 0.0.0.0", func(t *testing.T) {
		defer clearEnv()
		clearEnv()
		os.Setenv("GO_ENV", "production")
		os.Setenv("APP_PORT", "1000")
		os.Setenv("APP_SHUTDOWN_PERIOD", "1000")

		config := NewAppConfig()
		assert.Equal(t, false, config.IsDevelopment)
		assert.Equal(t, "0.0.0.0", config.Host)
		assert.Equal(t, "1000", config.Port)
		assert.Equal(t, time.Duration(1000), config.ShutdownPeriod)
	})
}
