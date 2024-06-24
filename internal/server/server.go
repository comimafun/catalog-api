package server

import (
	"catalog-be/internal/database"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type FiberServer struct {
	App       *fiber.App
	Pg        *gorm.DB
	Validator *validator.Validate
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "catalog-be",
			AppName:      "catalog-be",
		}),
		Pg:        database.New(),
		Validator: validator.New(),
	}

	return server
}
