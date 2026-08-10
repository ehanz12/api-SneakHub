package routes

import (
	"github.com/ehanz12/api-SneakHub/handlers"
	"github.com/ehanz12/api-SneakHub/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(api fiber.Router) {
	auth := api.Group("/auth")

	auth.Post("/register", handlers.RegisterHandler)
	auth.Post("/login", handlers.LoginHandler)

	user := api.Group("/users")
	user.Get("/me", middleware.AllRoles, handlers.MeUserHandler)
	user.Put("/me", middleware.AllRoles, handlers.UpdateUserHandler)
}
