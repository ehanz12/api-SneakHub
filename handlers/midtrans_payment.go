package handlers

import (
	"fmt"
	"strings"

	"github.com/ehanz12/api-SneakHub/services"
	"github.com/gofiber/fiber/v2"
)

// mapMidtransStatus memetakan transaction_status/fraud_status Midtrans
// ke status internal (PAID/EXPIRED/FAILED/REFUND). Status yang belum final
// (pending, authorize, dst) mengembalikan ok=false sehingga diabaikan.
func mapMidtransStatus(txStatus, fraudStatus string) (string, bool) {
	switch strings.ToLower(strings.TrimSpace(txStatus)) {
	case "settlement":
		return "PAID", true
	case "capture":
		if strings.EqualFold(strings.TrimSpace(fraudStatus), "accept") {
			return "PAID", true
		}
	case "deny", "cancel":
		return "FAILED", true
	case "expire":
		return "EXPIRED", true
	case "refund", "partial_refund":
		return "REFUND", true
	}
	return "", false
}

func MidtransNotificationHandler(c *fiber.Ctx) error {
	rawBody := c.Body()

	var body services.MidtransNotificationPayload
	if err := c.BodyParser(&body); err != nil {
		return c.JSON(fiber.Map{"success": false, "message": "request gagal"})
	}

	if !services.GetPaymentProvider().VerifySignature(rawBody, body.SignatureKey) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "signature tidak valid"})
	}

	if body.OrderID == "" {
		return c.JSON(fiber.Map{"success": false, "message": "data callback tidak lengkap"})
	}

	status, ok := mapMidtransStatus(body.TransactionStatus, body.FraudStatus)
	if !ok {
		return c.JSON(fiber.Map{"success": true, "message": "status belum final, diabaikan"})
	}

	grossAmount := 0.0
	if body.GrossAmount != "" {
		if _, err := fmt.Sscanf(body.GrossAmount, "%f", &grossAmount); err != nil {
			return c.JSON(fiber.Map{"success": false, "message": "gross_amount tidak valid"})
		}
	}

	if err := services.HandlePaymentNotificationService(body.OrderID, status, body.TransactionID, grossAmount); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}
