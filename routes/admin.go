package routes

import (
	"github.com/ehanz12/api-SneakHub/handlers"
	"github.com/ehanz12/api-SneakHub/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupAdminRoute(api fiber.Router) {
	admin := api.Group("/admin", middleware.AdminOnly)
	admin.Get("/users", handlers.AdminGetUsersHandler)
	admin.Patch("/users/:user_id/status", handlers.AdminUpdateUserStatusHandler)
	admin.Patch("/users/:user_id/role", handlers.AdminUpdateUserRoleHandler)
	admin.Get("/sellers", handlers.AdminGetSellersHandler)
	admin.Patch("/sellers/:seller_id/verification", handlers.AdminVerifySellerHandler)
	admin.Get("/products", handlers.AdminGetProductsHandler)
	admin.Patch("/products/:product_id/status", handlers.AdminUpdateProductStatusHandler)
	admin.Get("/orders", handlers.AdminGetOrdersHandler)
	admin.Get("/reports", handlers.AdminGetReportsHandler)
}
