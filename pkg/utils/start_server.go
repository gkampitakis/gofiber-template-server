package utils

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gkampitakis/gofiber-template-server/pkg/configs"

	"github.com/gofiber/fiber/v2"
)

func StartServerWithGracefulShutdown(app *fiber.App, config *configs.AppConfig) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		_ = <-c
		Logger.Info("Gracefully shutting down ...")
		err := app.Shutdown()
		if err != nil {
			Logger.Fatal(err.Error())
		}
	}()

	if err := app.Listen(config.Host + ":" + config.Port); err != nil {
		Logger.Panic(err.Error())
	}

	select {
	case <-c:
		Logger.Info("Received second signal")
		break
	case <-time.After(time.Second * config.ShutdownPeriod):
		break
	}
}

func StartServer(app *fiber.App, config *configs.AppConfig) {
	if !config.IsDevelopment {
		StartServerWithGracefulShutdown(app, config)
		return
	}

	if err := app.Listen(config.Host + ":" + config.Port); err != nil {
		Logger.Panic(err.Error())
	}
}
