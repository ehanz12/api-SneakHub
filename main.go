package main

import (
	"log"
	"os"
	"strings"

	"github.com/ehanz12/api-SneakHub/config"
	"github.com/ehanz12/api-SneakHub/database"
	"github.com/ehanz12/api-SneakHub/routes"
	"github.com/ehanz12/api-SneakHub/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	config.LoadEnv()

	database.ConnectDB()
	services.StartShippingScheduler()
	app := fiber.New()

	corsOrigins := config.AppConfig.CORSOrigins
	if strings.TrimSpace(corsOrigins) == "" {
		corsOrigins = "http://localhost:3060"
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     corsOrigins,
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	if err := os.MkdirAll("uploads", 0755); err != nil {
		log.Fatal("⚠️ GAGAL MEMBUAT FOLDER UPLOADS !", err)
	}
	app.Static("/uploads", "./uploads")

	routes.SetupRoutes(app)

	log.Println(app.Listen(":" + config.AppConfig.Port))
}
