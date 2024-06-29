package main

import (
	"catalog-be/internal"
	"catalog-be/internal/server"
	"catalog-be/internal/utils"
	"catalog-be/internal/validation"
	"catalog-be/seed"
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	server := server.New()
	validation.NewCustomValidation(server.Validator).Init()

	seed := seed.NewSeed(server.Pg)
	seed.Run()

	server.App.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("ALLOWED_ORIGINS"),
		AllowCredentials: true,
	}))
	server.App.Use(requestid.New(requestid.Config{
		Generator: func() string {
			return utils.NewUtils().GenerateRandomCode(7)
		},
		ContextKey: "requestid",
	}))
	server.App.Use(logger.New(logger.Config{
		TimeFormat: "02-Jan-2006, 15:04:05",
		TimeZone:   "Asia/Jakarta",
		Format:     "${locals:requestid} | ${time}WIB | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
	}))

	internal.InitializeServer(
		server.Pg,
		server.Validator,
		server.S3,
	).RegisterRoutes(server.App)

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	err := server.App.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
