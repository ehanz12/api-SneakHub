package handlers

import (
	"github.com/ehanz12/api-SneakHub/mappers"
	"github.com/ehanz12/api-SneakHub/requests"
	"github.com/ehanz12/api-SneakHub/services"
	"github.com/gofiber/fiber/v2"
)

func BrandCreateHandler(c *fiber.Ctx) error {
	var r requests.BrandRequest
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
	brand, err := services.CreateBrandService(r)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"success": true, "message": "berhasil membuat brand", "data": mappers.ToBrandResponse(brand)})
}

func GetBrandHandler(c *fiber.Ctx) error {
	brands, err := services.GetBrandService()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "terjadi kesalahan server", "errors": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"success": true, "message": "success memuat brand", "data": mappers.ListToBrandResponse(brands)})
}

func UpdateBrandHandler(c *fiber.Ctx) error {
	brandID := c.Params("brand_id")
	var r requests.BrandRequest
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
	brand, err := services.UpdateBrandService(brandID, r)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"success": true, "message": "berhasil update brand", "data": mappers.ToBrandResponse(brand)})
}

func DeleteBrandHandler(c *fiber.Ctx) error {
	brandID := c.Params("brand_id")
	if err := services.DeleteBrandService(brandID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"success": true, "message": "berhasil menghapus brand", "data": nil})
}
