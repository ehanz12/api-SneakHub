package routes

import (
	"github.com/ehanz12/api-SneakHub/handlers"
	"github.com/ehanz12/api-SneakHub/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupBrandRoute(api fiber.Router) {
	b := api.Group("/brand")
	b.Get("/", middleware.AllRoles, handlers.GetBrandHandler)
	b.Post("/", middleware.AdminOnly, handlers.BrandCreateHandler)
	b.Put("/:brand_id", middleware.AdminOnly, handlers.UpdateBrandHandler)
	b.Delete("/:brand_id", middleware.AdminOnly, handlers.DeleteBrandHandler)
}
