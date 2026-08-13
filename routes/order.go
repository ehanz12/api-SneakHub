package routes

import (
	"github.com/ehanz12/api-SneakHub/handlers"
	"github.com/ehanz12/api-SneakHub/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupOrderRoute(api fiber.Router) {
	order := api.Group("/orders")
	order.Get("/", middleware.AllRoles, handlers.GetOrdersHandler)
	order.Get("/:order_id", middleware.AllRoles, handlers.GetOrderHandler)
	order.Post("/", middleware.AllRoles, handlers.CreateOrderHandler)
	order.Put("/:order_id", middleware.AllRoles, handlers.UpdateOrderHandler)
	order.Delete("/:order_id", middleware.AllRoles, handlers.DeleteOrderHandler)
	order.Post("/:order_id/review", middleware.CustomerOnly, handlers.CreateReviewHandler)
}
