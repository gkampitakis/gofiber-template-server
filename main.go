package main

import (
	_ "github.com/gkampitakis/gofiber-template-server/docs"
	"github.com/gkampitakis/gofiber-template-server/pkg/configs"
	"github.com/gkampitakis/gofiber-template-server/pkg/middleware"
	"github.com/gkampitakis/gofiber-template-server/pkg/routes"
	"github.com/gkampitakis/gofiber-template-server/pkg/utils"

	"github.com/gkampitakis/fiber-modules/gracefulshutdown"
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
	cfg := configs.New()
	app := SetupServer(cfg)

	gracefulshutdown.Listen(app, cfg.Addr, cfg.GracefulshutdownConfig())
}

func SetupServer(cfg *configs.AppConfig) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler:          utils.ErrorHandler,
		DisableStartupMessage: !cfg.IsDevelopment,
	})
	middleware.FiberMiddleware(app, cfg.IsDevelopment)

	/**
	Register Routes
	*/
	routes.AppRoutes(app)

	/**
	Special Setup for development
	*/
	if cfg.IsDevelopment {
		routes.SwaggerRoute(app)
	}

	app.Get("/health", hc.New(
		hc.SetServiceName(cfg.Service),
		hc.ShowErrors(),
		hc.EnableTimeout(),
	))

	return app
}
