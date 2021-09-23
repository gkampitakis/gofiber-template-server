package utils

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, e error) error {
	return c.Status(http.StatusInternalServerError).JSON(map[string]interface{}{
		"status":  http.StatusInternalServerError,
		"message": "Internal Server Error",
	})
}
