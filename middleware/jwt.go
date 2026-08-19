package middleware

import (
	"strings"

	"github.com/ehanz12/api-SneakHub/database"
	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/utils"
	"github.com/golang-jwt/jwt/v5"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return c.Status(401).JSON(fiber.Map{"error": "missing authorization header"})
		}

		if !strings.HasPrefix(header, "Bearer ") {
			return c.Status(401).JSON(fiber.Map{"error": "invalid authorization format"})
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")

		token, err := utils.VerifyToken(tokenStr)
		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{"error": err.Error()})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "invalid jwt claims"})
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			return c.Status(401).JSON(fiber.Map{"error": "invalid user_id claim"})
		}

		var user models.User
		if err := database.DB.Select("peran", "status_akun").
			Where("user_id = ?", userID).First(&user).Error; err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "user tidak ditemukan"})
		}
		if user.StatusAkun != "aktif" {
			return c.Status(403).JSON(fiber.Map{"error": "akun tidak aktif"})
		}
		role := user.Peran

		roleAllowed := false
		for _, r := range allowedRoles {
			if r == role {
				roleAllowed = true
				break
			}
		}
		if !roleAllowed {
			return c.Status(403).JSON(fiber.Map{"error": "akses ditolak untuk role: " + role})
		}

		c.Locals("user_id", userID)
		c.Locals("role", role)

		return c.Next()
	}
}

var (
	AllRoles           = AuthMiddleware("customer", "seller", "admin")
	CustomerOnly       = AuthMiddleware("customer")
	CustomerSellerOnly = AuthMiddleware("customer", "seller")
	AdminOnly          = AuthMiddleware("admin")
	SellerOnly         = AuthMiddleware("seller")
	AdminSellerOnly    = AuthMiddleware("seller", "admin")
)
