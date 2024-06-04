package internal

import "github.com/gofiber/fiber/v2"

type HTTP struct{}

func (h *HTTP) RegisterRoutes(app *fiber.App) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "UP",
		})
	})

	app.Group("/api/v1")
}

func NewHTTP() *HTTP {
	return &HTTP{}
}
