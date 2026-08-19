package handlers

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/ehanz12/api-SneakHub/services"
	"github.com/gofiber/fiber/v2"
)

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
	log.Printf("[midtrans-notification] callback masuk: %s", string(rawBody))

	var body services.MidtransNotificationPayload
	if err := c.BodyParser(&body); err != nil {
		log.Printf("[midtrans-notification] body tidak valid: %v", err)
		return c.JSON(fiber.Map{"success": false, "message": "request gagal"})
	}

	if !services.GetPaymentProvider().VerifySignature(rawBody, body.SignatureKey) {
		log.Printf("[midtrans-notification] signature TIDAK valid untuk order %q", body.OrderID)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "signature tidak valid"})
	}
	log.Printf("[midtrans-notification] signature valid untuk order %q (status: %s)", body.OrderID, body.TransactionStatus)

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
		if errors.Is(err, services.ErrPaymentNotFound) {
			log.Printf("[midtrans-notification] order %q tidak ditemukan di DB, diabaikan", body.OrderID)
			return c.JSON(fiber.Map{"success": true, "message": "order tidak ditemukan, diabaikan"})
		}
		log.Printf("[midtrans-notification] gagal proses order %q: %v", body.OrderID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	log.Printf("[midtrans-notification] order %q berhasil diproses -> %s", body.OrderID, status)
	return c.JSON(fiber.Map{"success": true})
}
