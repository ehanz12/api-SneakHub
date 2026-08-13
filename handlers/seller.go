package handlers

import (
	"github.com/ehanz12/api-SneakHub/mappers"
	"github.com/ehanz12/api-SneakHub/requests"
	"github.com/ehanz12/api-SneakHub/services"
	"github.com/gofiber/fiber/v2"
)

func CreateSellerHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	var r requests.CreateSellerRequest
	if err := c.BodyParser(&r); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "request gagal", "errors": err.Error()})
	}
	if errs := r.Validation(); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"success": false,
			"message": "kesalahan validasi",
			"errors":  errs,
		})
	}

	seller, err := services.CreateSellerService(userID, r)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "terjadi kesalahan server", "errors": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"success": true, "message": "Akun Seller berhasil diaktifkan.", "data": mappers.ToSellerCreate(seller)})
}
