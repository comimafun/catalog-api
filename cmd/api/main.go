package main

import (
	"catalog-be/internal"
	"catalog-be/internal/server"
	"catalog-be/internal/utils"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"

	_ "github.com/joho/godotenv/autoload"
)

func main() {

	server := server.New()

	println("again", time.Now().Format(time.RFC3339))
	server.App.Use(cors.New())
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

	internal.InitializeServer(server.Pg, server.Validator).RegisterRoutes(server.App)

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	err := server.App.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}