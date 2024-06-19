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

	server.Validator.RegisterValidation("url_or_empty", func(fl validator.FieldLevel) bool {
		urlString := fl.Field().String()
		if urlString == "" {
			return true
		}
		return server.Validator.Var(urlString, "url") == nil
	})

	server.Validator.RegisterValidation("day_or_empty", func(fl validator.FieldLevel) bool {
		day := fl.Field().String()
		if day == "" {
			return true
		}
		// day is enum oneof first second both
		return day == "first" || day == "second" || day == "both"
	})

	return server
}
