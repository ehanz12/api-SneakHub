package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	apiLimiter := limiter.New(limiter.Config{
		Max:        200,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
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
	SetupAdminRoute(api)
	SetupDiscoveryRoute(api)
	SetupCategoryRoute(api)
	SetupBrandRoute(api)
	SetupWishlistRoute(api)
	SetupNotificationRoute(api)
	SetupProductRoute(api)
	SetupCartRoute(api)
	SetupAddressRoute(api)
	SetupCheckoutRoute(api)
	SetupOrderRoute(api)
}
