package handlers

import (
	"github.com/ehanz12/api-SneakHub/mappers"
	"github.com/ehanz12/api-SneakHub/requests"
	"github.com/ehanz12/api-SneakHub/services"
	"github.com/gofiber/fiber/v2"
)

func CreateReviewHandler(c *fiber.Ctx) error {
	customerID := c.Locals("user_id").(string)
	orderID := c.Params("order_id")

	var r requests.CreateReviewRequest
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

	review, err := services.CreateReviewService(customerID, orderID, r)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"success": true, "message": "Review berhasil dibuat.", "data": mappers.ToReviewResponse(*review)})
}
