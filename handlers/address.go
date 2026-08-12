package handlers

import (
	"github.com/ehanz12/api-SneakHub/mappers"
	"github.com/ehanz12/api-SneakHub/requests"
	"github.com/ehanz12/api-SneakHub/services"
	"github.com/gofiber/fiber/v2"
)

func GetAddressesHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	addresses, err := services.GetAddressesService(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Daftar alamat berhasil diambil.",
		"data":    mappers.ToAddressListResponse(addresses),
	})
}

func GetAddressHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	addressID := c.Params("address_id")

	address, err := services.GetAddressService(userID, addressID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Alamat berhasil diambil.",
		"data":    mappers.ToAddressResponse(*address),
	})
}

func CreateAddressHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var r requests.AddressCreateRequest
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

	address, err := services.CreateAddressService(userID, r)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"success": true,
		"message": "Alamat berhasil ditambahkan.",
		"data":    mappers.ToAddressResponse(*address),
	})
}

func UpdateAddressHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	addressID := c.Params("address_id")

	var r requests.AddressUpdateRequest
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

	address, err := services.UpdateAddressService(userID, addressID, r)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Alamat berhasil diperbarui.",
		"data":    mappers.ToAddressResponse(*address),
	})
}

func DeleteAddressHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	addressID := c.Params("address_id")

	if err := services.DeleteAddressService(userID, addressID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Alamat berhasil dihapus.",
		"data":    nil,
	})
}