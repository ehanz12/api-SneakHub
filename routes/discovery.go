package routes

import (
	"github.com/ehanz12/api-SneakHub/handlers"
	"github.com/ehanz12/api-SneakHub/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupDiscoveryRoute(api fiber.Router) {
	discovery := api.Group("/discovery")
	discovery.Post("/smart-filter", middleware.AllRoles, handlers.SmartFilterHandler)

	home := api.Group("/home")
	home.Get("/personalized", middleware.AllRoles, handlers.HomePersonalizedHandler)

	recommendation := api.Group("/recommendation")
	recommendation.Get("/cocok-untuk-kamu", middleware.AllRoles, handlers.PersonalizedRecommendationHandler)

	api.Get("/trending", middleware.AllRoles, handlers.TrendingHandler)
	api.Get("/best-seller/weekly", middleware.AllRoles, handlers.BestSellerWeeklyHandler)

	price := api.Group("/price")
	price.Post("/prediction", middleware.SellerOnly, handlers.PricePredictionHandler)
}
