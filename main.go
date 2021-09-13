package main

import (
	"github.com/gkampitakis/gofiber-template-server/pkg/configs"
	"github.com/gkampitakis/gofiber-template-server/pkg/middleware"
	"github.com/gkampitakis/gofiber-template-server/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app_config := configs.NewAppConfig()
	app := SetupServer(app_config.IsDevelopment)

	utils.StartServer(app, app_config)
}

func SetupServer(isDevelopment bool) *fiber.App {
	app := fiber.New()
	middleware.FiberMiddleware(app, isDevelopment)
	utils.RegisterHealthchecks(app)

	/**
	--- Example healthcheck ---
	checks := []func(c chan utils.HealthcheckResult){
		func(c chan utils.HealthcheckResult) {
			go func() {
				time.Sleep(time.Second * 6)
				c <- utils.HealthcheckResult{Label: "postgres", Result: true}
			}()
		},
	}
	*/

	return app
}
