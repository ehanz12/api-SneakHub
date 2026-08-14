package handlers

import (
	"github.com/ehanz12/api-SneakHub/requests"
	"github.com/ehanz12/api-SneakHub/services"
	"github.com/gofiber/fiber/v2"
)

func PricePredictionHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var r requests.PricePredictionRequest
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

	result, err := services.PredictPriceService(userID, r)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Prediksi harga berhasil dibuat.", "data": result})
}

func PriceInsightHandler(c *fiber.Ctx) error {
	productID := c.Params("product_id")

	result, err := services.PriceInsightService(productID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Price Insight berhasil diambil.", "data": result})
}
