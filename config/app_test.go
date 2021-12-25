package config

import (
	"os"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestMain(t *testing.M) {
	v := t.Run()

	snaps.Clean()

	os.Exit(v)
}

func TestAppConfig(t *testing.T) {
	t.Run("should return default values", func(t *testing.T) {
		config := New()

		snaps.MatchSnapshot(t, config, config.GracefulshutdownConfig())
	})

	t.Run("should return explicit values and host as 0.0.0.0", func(t *testing.T) {
		t.Setenv("GO_ENV", "production")
		t.Setenv("APP_PORT", "1000")
		t.Setenv("GRACE_PERIOD", "35")

		fns := []func() error{}
		config := New()

		snaps.MatchSnapshot(t, config, config.GracefulshutdownConfig(fns...))
	})
}
