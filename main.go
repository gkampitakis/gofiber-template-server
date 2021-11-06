package main

import (
	_ "github.com/gkampitakis/gofiber-template-server/docs"
	"github.com/gkampitakis/gofiber-template-server/pkg/configs"
	"github.com/gkampitakis/gofiber-template-server/pkg/middleware"
	"github.com/gkampitakis/gofiber-template-server/pkg/routes"
	"github.com/gkampitakis/gofiber-template-server/pkg/utils"

	hc "github.com/gkampitakis/fiber-modules/healthcheck"

	"github.com/gofiber/fiber/v2"
)

// @title Gofiber Template Server
// @version 1.0.1
// @description Template for spinning up a gofiber server
// @contact.name gkampitakis
// @contact.email gkabitakis@gmail.com
// @license.name MIT
// @license.name https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
func main() {
	app_cfg := configs.NewAppConfig()
	hc_cfg := configs.NewHealthcheckConfig()
	app := SetupServer(hc_cfg, app_cfg.IsDevelopment)

	utils.StartServer(app, app_cfg)
}

func SetupServer(cfg *configs.HealthcheckConfig, isDevelopment bool) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler:          utils.ErrorHandler,
		DisableStartupMessage: !isDevelopment,
	})
	middleware.FiberMiddleware(app, isDevelopment)

	/**
	Register Routes
	*/
	routes.AppRoutes(app)

	/**
	Special Setup for development
	*/
	if isDevelopment {
		routes.SwaggerRoute(app)
	}

	app.Get("/health", hc.New(
		hc.SetServiceName(cfg.Service),
		hc.ShowErrors(),
		hc.EnableTimeout(),
	))

	return app
}
