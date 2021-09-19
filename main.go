package main

import (
	"time"

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
	app := fiber.New(fiber.Config{
		ErrorHandler:          utils.ErrorHandler,
		DisableStartupMessage: !isDevelopment,
	})
	middleware.FiberMiddleware(app, isDevelopment)

	/**
	Register Routes
	*/
	routes.AppRoutes(app)
	if isDevelopment {
		routes.SwaggerRoute(app)
	}

	/**
	--- Example healthcheck ---
	*/

	checks := utils.HealthcheckMap{
		"test": func() bool {
			time.Sleep(4 * time.Second)
			return true
		},
		"test2": func() bool {
			time.Sleep(3 * time.Second)
			return true
		},
		"test3": func() bool {
			time.Sleep(10 * time.Second)
			return true
		},
	}

	utils.RegisterHealthchecks(app, checks)
	return app
}
