package configs

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func getEnvFilePath() string {
	root, _ := os.Getwd()
	return filepath.Join(root, ".env")
}

func clearEnv() {
	os.Unsetenv("GO_ENV")
	os.Unsetenv("APP_PORT")
	os.Unsetenv("APP_SHUTDOWN_PERIOD")
	os.Remove(getEnvFilePath())
}

func setEnvFile(t *testing.T, port, shutdownPeriod string) {
	data := []byte(fmt.Sprintf("APP_PORT=%s\nAPP_SHUTDOWN_PERIOD=%s", port, shutdownPeriod))
	err := os.WriteFile(getEnvFilePath(), data, 0777)
	if err != nil {
		t.Error(err)
	}
}

func TestAppConfig(t *testing.T) {
	t.Run("development", func(t *testing.T) {
		t.Run("should load env from file", func(t *testing.T) {
			defer clearEnv()
			clearEnv()
			setEnvFile(t, "1000", "15")

			config := NewAppConfig()
			assert.Equal(t, true, config.IsDevelopment)
			assert.Equal(t, "localhost", config.Host)
			assert.Equal(t, "1000", config.Port)
			assert.Equal(t, time.Duration(15), config.ShutdownPeriod)
		})
	})

	t.Run("production", func(t *testing.T) {
		t.Run("should not load env from file and set host to 0.0.0.0", func(t *testing.T) {
			defer clearEnv()
			clearEnv()
			setEnvFile(t, "1000", "10")
			os.Setenv("GO_ENV", "production")

			config := NewAppConfig()
			assert.Equal(t, false, config.IsDevelopment)
			assert.Equal(t, "0.0.0.0", config.Host)
			assert.Equal(t, "8080", config.Port)
			assert.Equal(t, time.Duration(10), config.ShutdownPeriod)
		})
	})
}
