package middleware

import (
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// FiberMiddleware provide Fiber's built-in middlewares.
// See: https://docs.gofiber.io/api/middleware
func FiberMiddleware(a *fiber.App, isDevelopment bool) {
	if isDevelopment {
		a.Use(pprof.New())
		a.Get("/monitor", monitor.New())
	}

	a.Use(
		cors.New(),
		recover.New(recover.Config{
			EnableStackTrace: true,
		}),
		requestid.New(),
		logger.New(),
	)
}
