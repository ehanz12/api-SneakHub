package handlers

import (
	"github.com/ehanz12/api-SneakHub/mappers"
	"github.com/ehanz12/api-SneakHub/requests"
	"github.com/ehanz12/api-SneakHub/services"
	"github.com/gofiber/fiber/v2"
)

func CreateConditionScoreHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)
	productID := c.Params("product_id")

	var r requests.ConditionScoreRequest
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

	cs, err := services.CreateConditionScoreService(userID, role, productID, r)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Condition Score berhasil dihitung.", "data": mappers.ToConditionScoreCreateResponse(*cs)})
}

func GetConditionScoreHandler(c *fiber.Ctx) error {
	productID := c.Params("product_id")

	cs, err := services.GetConditionScoreService(productID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Condition Score berhasil diambil.", "data": mappers.ToConditionScoreGetResponse(*cs)})
}

func SellerTrustScoreHandler(c *fiber.Ctx) error {
	sellerID := c.Params("seller_id")

	result, err := services.GetSellerTrustScoreService(sellerID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Seller Trust Score berhasil diambil.", "data": result})
}
