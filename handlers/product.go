package handlers

import (
	"github.com/ehanz12/api-SneakHub/mappers"
	"github.com/ehanz12/api-SneakHub/requests"
	"github.com/ehanz12/api-SneakHub/services"
	"github.com/gofiber/fiber/v2"
)

func CreateProductHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	var r requests.CreateProduct
	if err := c.BodyParser(&r); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": "false", "message": "request gagal", "errors": err.Error()})
	}
	if errs := r.Validate(); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "error",
			"message": "kesalahan validasi",
			"errors":  errs,
		})
	}
	product, err := services.CreateProductService(userID, r)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": "false", "message": "request gagal", "errors": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"status": "true", "message": "product berhasil dibuat", "data": mappers.ToProductResponse(*product)})
}

func UpdateProductHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	productID := c.Params("product_id")
	var r requests.CreateProduct
	if err := c.BodyParser(&r); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": "false", "message": "request gagal", "errors": err.Error()})
	}
	if errs := r.Validate(); len(errs) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "error",
			"message": "kesalahan validasi",
			"errors":  errs,
		})
	}
	product, err := services.UpdateProductService(userID, productID, r)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": "false", "message": "request gagal", "errors": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"status": "true", "message": "product berhasil di update", "data": mappers.ToProductUpdateResponse(*product)})
}
