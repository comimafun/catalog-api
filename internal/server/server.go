package server

import (
	"catalog-be/internal/database"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type FiberServer struct {
	App       *fiber.App
	Pg        *gorm.DB
	Validator *validator.Validate
	S3        *s3.Client
}

func New() *FiberServer {
	dsn := os.Getenv("DB_URL")
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "catalog-be",
			AppName:      "catalog-be",
		}),
		Pg:        database.New(dsn),
		Validator: validator.New(),
		S3:        database.NewS3(),
	}

	return server
}
