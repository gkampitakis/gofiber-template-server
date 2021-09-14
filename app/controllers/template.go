package controllers

import "github.com/gofiber/fiber/v2"

// @Description Default func prints hello world
// @Summary Prints Hello World
// @Tags template
// @Accept text/plain
// @Produce text/plain
// @Success 200 {string} greet
// @Router / [get]
func Default(c *fiber.Ctx) error {
	return c.Send([]byte("Hello World from template"))
}

// @Description Hello func just Greets with the name provided in params
// @Summary Greets user
// @Tags template
// @Accept text/plain
// @Produce text/plain
// @Success 200 {string} greet
// @Router /hello/:name [get]
func Hello(c *fiber.Ctx) error {
	name := c.Params("name")

	return c.Send([]byte("Hello " + name + " 👋"))
}
