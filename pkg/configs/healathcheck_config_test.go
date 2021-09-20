package configs

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHealthcheckConfig(t *testing.T) {
	t.Run("should return default values", func(t *testing.T) {
		clearEnv()
		config := NewHealthcheckConfig()

		assert.Equal(t, time.Duration(5), config.TimeoutPeriod)
		assert.Equal(t, true, config.TimeoutEnabled)
	})

	t.Run("should return default values in case of wrong values", func(t *testing.T) {
		clearEnv()
		os.Setenv("HC_TIMEOUT_PERIOD", "wrong")
		os.Setenv("HC_TIMEOUT_ENABLED", "wrong")

		config := NewHealthcheckConfig()
		assert.Equal(t, time.Duration(5), config.TimeoutPeriod)
		assert.Equal(t, true, config.TimeoutEnabled)
	})

	t.Run("should return setted values", func(t *testing.T) {
		clearEnv()

		os.Setenv("HC_TIMEOUT_PERIOD", "10")
		os.Setenv("HC_TIMEOUT_ENABLED", "false")

		config := NewHealthcheckConfig()
		assert.Equal(t, time.Duration(10), config.TimeoutPeriod)
		assert.Equal(t, false, config.TimeoutEnabled)
	})
}
