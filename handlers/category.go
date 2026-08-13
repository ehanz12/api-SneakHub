package handlers

import (
	"github.com/ehanz12/api-SneakHub/mappers"
	"github.com/ehanz12/api-SneakHub/requests"
	"github.com/ehanz12/api-SneakHub/services"
	"github.com/gofiber/fiber/v2"
)

func CategoryCreateHandler(c *fiber.Ctx) error {
	var r requests.CategoryRequest
	if err := c.BodyParser(&r); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "request gagal", "errors": err.Error()})
	}
	if len(r.NamaKategori) < 3 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "nama kategori harus lebih 3 karakter", "errors": "kesalahan"})
	}
	category, err := services.CreateCategoryService(r)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "terjadi kesalahan server", "errors": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"success": true, "message": "berhasil membuat category", "data": mappers.ToCategoryRes(category)})
}

func GetCategoryHandler(c *fiber.Ctx) error {
	category, err := services.GetCategoryService()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "terjadi kesalahan server", "errors": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"success": true, "message": "success memuat category", "data": mappers.ListToCategoryResponse(category)})
}

func UpdateCategoryHandler(c *fiber.Ctx) error {
	cID := c.Params("category_id")
	var r requests.CategoryRequest
	if err := c.BodyParser(&r); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "request gagal", "errors": err.Error()})
	}
	if len(r.NamaKategori) < 3 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "nama kategori harus lebih 3 karakter", "errors": "kesalahan"})
	}
	category, err := services.UpdateCategoryService(cID, r)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "terjadi kesalahan server", "errors": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"success": true, "message": "berhasil update category", "data": mappers.ToCategoryRes(category)})
}

func DeleteCategoryHandler(c *fiber.Ctx) error {
	cID := c.Params("category_id")
	if err := services.DeleteCategoryService(cID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "terjadi kesalahan server", "errors": err.Error()})
	}
	return c.Status(200).JSON(fiber.Map{"success": true, "message": "berhasil menghapus category", "data": nil})
}
