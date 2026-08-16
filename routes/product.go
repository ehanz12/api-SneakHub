package routes

import (
	"github.com/ehanz12/api-SneakHub/handlers"
	"github.com/ehanz12/api-SneakHub/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupProductRoute(api fiber.Router) {
	product := api.Group("/products")
	product.Get("/", middleware.AllRoles, handlers.GetProductsHandler)
	product.Get("/:product_id", middleware.AllRoles, handlers.GetProductByIDHandler)
	product.Post("/search-by-image", middleware.AllRoles, handlers.SearchProductByImageHandler)
	product.Post("/", middleware.AdminSellerOnly, handlers.CreateProductHandler)
	product.Put("/:product_id", middleware.AdminSellerOnly, handlers.UpdateProductHandler)
	product.Delete("/:product_id", middleware.AdminSellerOnly, handlers.DeleteProductHandler)
	product.Post("/:product_id/images", middleware.AdminSellerOnly, handlers.UploadProductImageHandler)
	product.Get("/:product_id/images", middleware.AllRoles, handlers.ListProductImagesHandler)
	product.Delete("/:product_id/images/:image_id", middleware.AdminSellerOnly, handlers.DeleteProductImageHandler)
	product.Get("/:product_id/price-insight", middleware.AllRoles, handlers.PriceInsightHandler)
	product.Get("/:product_id/reviews", middleware.AllRoles, handlers.GetProductReviewsHandler)
	product.Post("/:product_id/condition-score", middleware.AdminSellerOnly, handlers.CreateConditionScoreHandler)
	product.Get("/:product_id/condition-score", middleware.AllRoles, handlers.GetConditionScoreHandler)
}
