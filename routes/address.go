package routes

import (
	"github.com/ehanz12/api-SneakHub/handlers"
	"github.com/ehanz12/api-SneakHub/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupAddressRoute(api fiber.Router) {
	address := api.Group("/addresses")
	address.Get("/", middleware.CustomerOnly, handlers.GetAddressesHandler)
	address.Get("/:address_id", middleware.CustomerOnly, handlers.GetAddressHandler)
	address.Post("/", middleware.CustomerOnly, handlers.CreateAddressHandler)
	address.Put("/:address_id", middleware.CustomerOnly, handlers.UpdateAddressHandler)
	address.Delete("/:address_id", middleware.CustomerOnly, handlers.DeleteAddressHandler)
}