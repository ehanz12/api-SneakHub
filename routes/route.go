package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	apiLimiter := limiter.New(limiter.Config{
		Max:        200,             // Maksimal request
		Expiration: 1 * time.Minute, // Durasi limit
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // Berdasarkan alamat IP klien
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"status":  "error",
				"message": "Too many requests. Please try again later.",
			})
		},
	})
	api.Use(apiLimiter)
	SetupUserRoutes(api)
	SetupSellerRoute(api)
	SetupCategoryRoute(api)
	SetupProductRoute(api)
}
