package main

import (
	_ "github.com/gkampitakis/gofiber-template-server/docs"
	"github.com/gkampitakis/gofiber-template-server/pkg/configs"
	"github.com/gkampitakis/gofiber-template-server/pkg/middleware"
	"github.com/gkampitakis/gofiber-template-server/pkg/routes"
	"github.com/gkampitakis/gofiber-template-server/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// @title Gofiber Template Server
// @version 1.0.0
// @description Template for spinning up a gofiber server
// @contact.name gkampitakis
// @contact.email gkabitakis@gmail.com
// @license.name MIT
// @license.name https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
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
	Register Routes
	*/
	routes.AppRoutes(app)
	routes.SwaggerRoute(app)

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
