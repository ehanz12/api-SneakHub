package routes

import (
	"github.com/ehanz12/api-SneakHub/handlers"
	"github.com/ehanz12/api-SneakHub/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupCartRoute(api fiber.Router) {
	cart := api.Group("/cart")
	cart.Get("/", middleware.CustomerOnly, handlers.GetCartHandler)
	cart.Post("/items", middleware.CustomerOnly, handlers.AddCartItemsHandler)
	cart.Put("/items/:cart_item_id", middleware.CustomerOnly, handlers.UpdateCartItemHandler)
	cart.Delete("/items/:cart_item_id", middleware.CustomerOnly, handlers.DeleteCartItemHandler)
}
