package routes

import (
	"github.com/gkampitakis/gofiber-template-server/app/controllers"
	"github.com/gofiber/fiber/v2"
)

func AppRoutes(app *fiber.App) {
	route := app.Group("")

	route.Get("/", controllers.Default)
	route.Get("/hello/:name", controllers.Hello)
}
