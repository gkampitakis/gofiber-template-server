package utils

import "github.com/gofiber/fiber/v2"

type HealthcheckResult struct {
	_      struct{}
	Result bool
	Label  string
}
type HealthcheckControll func(result chan HealthcheckResult)

func RegisterHealthchecks(app *fiber.App, controls ...HealthcheckControll) {
	app.Get("/health", registerHealthRoute(controls...))
}

func registerHealthRoute(controls ...HealthcheckControll) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		health := make(chan HealthcheckResult)
		response := make(map[string]string)
		status := 200
		controlsLength := len(controls)

		if controlsLength == 0 {
			response["status"] = "healthy"
		}

		for _, control := range controls {
			control(health)
		}

		for i := 0; i < controlsLength; i++ {
			c := <-health
			if c.Result {
				response[c.Label] = "healthy"
				continue
			}

			response[c.Label] = "unhealthy"
			status = 500
		}

		return c.Status(status).JSON(response)
	}
}
