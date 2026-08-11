// @title My API
// @version 1.0
// @description Ini adalah dokumentasi API gue
// @host www.reihan.biz.id
// @BasePath /api/v1

package main

import (
	"log"
	"os"

	"github.com/ehanz12/api-SneakHub/config"
	"github.com/ehanz12/api-SneakHub/database"
	"github.com/ehanz12/api-SneakHub/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// load env
	config.LoadEnv()

	// connect to database
	database.ConnectDB()
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3060",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true, // jika pake jwt
	}))

	if err := os.MkdirAll("uploads", 0755); err != nil {
		log.Fatal("⚠️ GAGAL MEMBUAT FOLDER UPLOADS !", err)
	}
	app.Static("/uploads", "./uploads")

	routes.SetupRoutes(app)

	log.Println(app.Listen(":" + config.AppConfig.Port))
}
