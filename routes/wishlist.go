package routes

import (
	"github.com/ehanz12/api-SneakHub/handlers"
	"github.com/ehanz12/api-SneakHub/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupWishlistRoute(api fiber.Router) {
	wishlist := api.Group("/wishlist", middleware.CustomerOnly)
	wishlist.Get("/", handlers.GetWishlistHandler)
	wishlist.Post("/", handlers.CreateWishlistHandler)
	wishlist.Delete("/:product_id", handlers.DeleteWishlistHandler)
	wishlist.Post("/:product_id/price-alert", handlers.SetPriceAlertHandler)
	wishlist.Delete("/:product_id/price-alert", handlers.DisablePriceAlertHandler)
	wishlist.Post("/:product_id/restock-alert", handlers.SetRestockAlertHandler)
	wishlist.Delete("/:product_id/restock-alert", handlers.DisableRestockAlertHandler)
}

func SetupNotificationRoute(api fiber.Router) {
	notification := api.Group("/notifications", middleware.AllRoles)
	notification.Get("/", handlers.GetNotificationsHandler)
	notification.Patch("/:notification_id/read", handlers.MarkNotificationReadHandler)
}
