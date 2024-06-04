package server

import (
	"catalog-be/internal/database"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type FiberServer struct {
	App       *fiber.App
	PG        *gorm.DB
	validator *validator.Validate
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "catalog-be",
			AppName:      "catalog-be",
		}),
		PG:        database.New(),
		validator: validator.New(),
	}

	return server
}
