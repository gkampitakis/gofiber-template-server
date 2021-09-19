package utils

import "github.com/gofiber/fiber/v2"

func ErrorHandler(c *fiber.Ctx, e error) error {
	return c.Status(500).JSON(map[string]interface{}{
		"status":  500,
		"message": "Internal Server Error",
	})
}
