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
	seller.Post("/orders/:order_id/ship", middleware.SellerOnly, handlers.ShipOrderHandler)
	seller.Get("/dashboard", middleware.SellerOnly, handlers.GetSellerDashboardHandler)

	sellerUser := api.Group("/users")
	sellerUser.Post("/me/seller-activation", middleware.AllRoles, handlers.CreateSellerHandler)

	sellers := api.Group("/sellers")
	sellers.Get("/:seller_id/trust-score", middleware.AllRoles, handlers.SellerTrustScoreHandler)
	sellers.Get("/:seller_id/reviews", middleware.AllRoles, handlers.GetSellerReviewsHandler)
}
