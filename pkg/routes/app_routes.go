package routes

import (
	"github.com/gkampitakis/fiber-modules/routevalidator"

	"github.com/gkampitakis/gofiber-template-server/app/controllers"
	"github.com/gofiber/fiber/v2"
)

func AppRoutes(app *fiber.App) {
	route := app.Group("")

	route.Get("/", controllers.Default)
	route.Get("/hello/:name", controllers.Hello)

	route.Post("/test/validation", validator(),
		func(c *fiber.Ctx) error {
			return c.Status(200).Send([]byte("Hello World"))
		})
}

var validator = func() fiber.Handler {
	return routevalidator.Validator(routevalidator.Config{
		SchemaPath: "pkg/routes/bigJSONValidator.json",
		// SchemaLiteral: t,
		ExposeErrors: true,
	})
}

const t = `{
  "body": {
    "type": "object",
    "properties": {
      "data": {
        "type": "array",
        "items": {
          "type": "integer"
        }
      }
    }
  },
  "querystring": {},
  "headers": {},
  "params": {}
}
`
