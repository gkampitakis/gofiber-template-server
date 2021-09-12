package middleware

import (
	"github.com/gofiber/fiber/v2/middleware/pprof"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// FiberMiddleware provide Fiber's built-in middlewares.
// See: https://docs.gofiber.io/api/middleware
func FiberMiddleware(a *fiber.App, isDevelopment bool) {

	if isDevelopment {
		a.Use(pprof.New())
	}

	a.Use(
		cors.New(),
	)
}
