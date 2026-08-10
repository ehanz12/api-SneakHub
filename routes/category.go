package routes

import (
	"github.com/ehanz12/api-SneakHub/handlers"
	"github.com/ehanz12/api-SneakHub/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupCategoryRoute(api fiber.Router) {
	c := api.Group("/category")
	c.Get("/", middleware.AllRoles, handlers.GetCategoryHandler)
	c.Post("/", middleware.AdminOnly, handlers.CategoryCreateHandler)
	c.Put("/:category_id", middleware.AdminOnly, handlers.UpdateCategoryHandler)
	c.Delete("/:category_id", middleware.AdminOnly, handlers.DeleteCategoryHandler)
}
