package server

import (
	"github.com/gofiber/fiber/v2"

	"catalog-be/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "catalog-be",
			AppName:      "catalog-be",
		}),

		db: database.New(),
	}

	return server
}
