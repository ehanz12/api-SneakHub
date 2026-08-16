package routes

import (
	"github.com/ehanz12/api-SneakHub/handlers"
	"github.com/ehanz12/api-SneakHub/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupAddressRoute(api fiber.Router) {
	address := api.Group("/addresses")
	address.Get("/", middleware.CustomerSellerOnly, handlers.GetAddressesHandler)
	address.Get("/:address_id", middleware.CustomerSellerOnly, handlers.GetAddressHandler)
	address.Post("/", middleware.CustomerSellerOnly, handlers.CreateAddressHandler)
	address.Put("/:address_id", middleware.CustomerSellerOnly, handlers.UpdateAddressHandler)
	address.Delete("/:address_id", middleware.CustomerSellerOnly, handlers.DeleteAddressHandler)
}
