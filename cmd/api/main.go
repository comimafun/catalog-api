package main

import (
	"catalog-be/internal"
	"catalog-be/internal/modules/user"
	"catalog-be/internal/server"
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

	server.App.Use(cors.New())
	server.App.Use(requestid.New())
	server.App.Use(logger.New(logger.Config{
		TimeFormat: "02-Jan-2006, 15:04:05",
		TimeZone:   "Asia/Jakarta",
		Format:     "{locals:requestid} | ${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
	}))

	user.NewUserRepo(server.PG)

	internal.NewHTTP().RegisterRoutes(server.App)

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	err := server.App.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
