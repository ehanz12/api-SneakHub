package handlers

import (
	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/requests"
	"github.com/ehanz12/api-SneakHub/responses"
	"github.com/ehanz12/api-SneakHub/services"
	"github.com/gofiber/fiber/v2"
)

func toCheckoutResponse(order models.Order) responses.CheckoutResponse {
	resp := responses.CheckoutResponse{
		OrderID:          order.OrderID,
		StatusOrder:      order.StatusOrder,
		Subtotal:         order.Subtotal,
		BiayaPengiriman:  order.BiayaPengiriman,
		TotalPembayaran:  order.TotalPesanan,
		MetodePembayaran: order.MetodePembayaran,
	}
	if order.Payment != nil && order.Payment.PaymentURL != nil {
		resp.PaymentURL = *order.Payment.PaymentURL
	}
	return resp
}

func CheckoutHandler(c *fiber.Ctx) error {
	customerID := c.Locals("user_id").(string)

	var r requests.CheckoutRequest
	if err := c.BodyParser(&r); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "request gagal", "errors": err.Error()})
	}
	if errs := r.Validate(); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "kesalahan validasi",
			"errors":  errs,
		})
	}

	orders, err := services.CheckoutService(customerID, r)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	data := make([]responses.CheckoutResponse, 0, len(orders))
	for _, order := range orders {
		data = append(data, toCheckoutResponse(order))
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "Checkout berhasil.",
		"data":    data,
	})
}

func PaymentNotificationHandler(c *fiber.Ctx) error {
	rawBody := c.Body()

	signature := c.Get("X-Callback-Signature")
	if !services.GetPaymentProvider().VerifySignature(rawBody, signature) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "signature tidak valid"})
	}

	if c.Get("X-Callback-Event") != "payment_status" {
		return c.JSON(fiber.Map{"success": false, "message": "event callback tidak dikenal"})
	}

	var body struct {
		Reference   string  `json:"reference"`
		MerchantRef string  `json:"merchant_ref"`
		TotalAmount float64 `json:"total_amount"`
		Status      string  `json:"status"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.JSON(fiber.Map{"success": false, "message": "request gagal"})
	}

	if body.MerchantRef == "" || body.Status == "" {
		return c.JSON(fiber.Map{"success": false, "message": "data callback tidak lengkap"})
	}

	if err := services.HandlePaymentNotificationService(body.MerchantRef, body.Status, body.Reference, body.TotalAmount); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"success": true})
}