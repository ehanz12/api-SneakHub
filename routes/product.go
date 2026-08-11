package routes

import (
	"github.com/ehanz12/api-SneakHub/handlers"
	"github.com/ehanz12/api-SneakHub/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupProductRoute(api fiber.Router) {
	product := api.Group("/products")
	product.Post("/search-by-image", middleware.AllRoles, handlers.SearchProductByImageHandler)
	product.Post("/", middleware.AdminSellerOnly, handlers.CreateProductHandler)
	product.Put("/:product_id", middleware.AdminSellerOnly, handlers.UpdateProductHandler)
	product.Post("/:product_id/images", middleware.AdminSellerOnly, handlers.UploadProductImageHandler)
	product.Get("/:product_id/images", middleware.AllRoles, handlers.ListProductImagesHandler)
}
