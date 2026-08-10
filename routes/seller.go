package routes

import (
	"github.com/ehanz12/api-SneakHub/handlers"
	"github.com/ehanz12/api-SneakHub/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupSellerRoute(api fiber.Router) {
	seller := api.Group("/users")
	seller.Post("/me/seller-activation", middleware.AllRoles, handlers.CreateSellerHandler)
}
