package routes

import (
	"github.com/ehanz12/api-SneakHub/handlers"
	"github.com/ehanz12/api-SneakHub/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupSellerRoute(api fiber.Router) {
	seller := api.Group("/seller")
	seller.Get("/products", middleware.SellerOnly, handlers.GetSellerProductsHandler)
	seller.Get("/orders", middleware.SellerOnly, handlers.GetSellerOrdersHandler)
	seller.Get("/dashboard", middleware.SellerOnly, handlers.GetSellerDashboardHandler)

	sellerUser := api.Group("/users")
	sellerUser.Post("/me/seller-activation", middleware.AllRoles, handlers.CreateSellerHandler)
}
