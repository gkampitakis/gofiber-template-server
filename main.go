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

	utils.StartServerWithGracefulShutdown(app, app_config)
}

func SetupServer(isDevelopment bool) *fiber.App {
	app := fiber.New()
	middleware.FiberMiddleware(app, isDevelopment)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(map[string]string{
			"status": "healthy",
		})
	})

	return app
}
