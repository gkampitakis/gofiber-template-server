package config

import (
	"os"
	"testing"
	"time"

	"github.com/gkampitakis/fiber-modules/gracefulshutdown"
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

		config := New()
		assert.Equal(t, true, config.IsDevelopment)
		assert.Equal(t, "localhost:8080", config.Addr)
	})

	t.Run("should return explicit values and host as 0.0.0.0", func(t *testing.T) {
		defer clearEnv()
		clearEnv()
		os.Setenv("GO_ENV", "production")
		os.Setenv("APP_PORT", "1000")
		os.Setenv("GRACE_PERIOD", "35")

		fns := []func() error{}

		config := New()
		assert.Equal(t, false, config.IsDevelopment)
		assert.Equal(t, "0.0.0.0:1000", config.Addr)
		assert.Equal(t, gracefulshutdown.Config{
			Period:      time.Duration(35),
			Enabled:     true,
			ShutdownFns: fns,
		}, config.GracefulshutdownConfig(fns...))
	})
}
