package router

import (
	"catalog-be/internal/modules/auth"

	"github.com/gofiber/fiber/v2"
)

type HTTP struct {
	auth *auth.AuthHandler
}

func (h *HTTP) RegisterRoutes(app *fiber.App) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":     "UP",
			"request_id": c.Locals("requestid"),
		})
	})

	v1 := app.Group("/api/v1")

	auth := v1.Group("/auth")
	auth.Get("/google", h.auth.GetAuthURL)
	auth.Get("/google/callback", h.auth.GetGoogleCallback)
	auth.Post("/google/callback", h.auth.PostGoogleCallback)
}

func NewHTTP(
	auth *auth.AuthHandler,
) *HTTP {
	return &HTTP{
		auth,
	}
}
