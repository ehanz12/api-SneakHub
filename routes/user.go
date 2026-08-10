package routes

import (
	"github.com/ehanz12/api-SneakHub/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(api fiber.Router) {
	auth := api.Group("/auth")

	auth.Post("/register", handlers.RegisterHandler)
	auth.Post("/login", handlers.LoginHandler)
}
