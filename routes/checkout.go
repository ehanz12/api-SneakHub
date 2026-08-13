package routes

import (
	"github.com/ehanz12/api-SneakHub/config"
	"github.com/ehanz12/api-SneakHub/handlers"
	"github.com/ehanz12/api-SneakHub/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupCheckoutRoute(api fiber.Router) {
	api.Post("/checkout", middleware.CustomerOnly, handlers.CheckoutHandler)
	api.Post("/payments/notification", handlers.PaymentNotificationHandler)

	if config.AppConfig.PaymentMode == "mock" {
		api.Get("/mock-pay/:order_id", handlers.MockPayPageHandler)
		api.Post("/payments/mock-settle", handlers.MockSettleHandler)
	}

	if config.AppConfig.PaymentMode == "midtrans" {
		api.Post("/payments/midtrans/notification", handlers.MidtransNotificationHandler)
	}
}
