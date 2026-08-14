package handlers

import (
	"github.com/ehanz12/api-SneakHub/mappers"
	"github.com/ehanz12/api-SneakHub/requests"
	"github.com/ehanz12/api-SneakHub/services"
	"github.com/gofiber/fiber/v2"
)

func GetWishlistHandler(c *fiber.Ctx) error {
	customerID := c.Locals("user_id").(string)

	wishlists, err := services.GetWishlistService(customerID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Wishlist berhasil diambil.", "data": mappers.ToWishlistListResponse(wishlists)})
}

func CreateWishlistHandler(c *fiber.Ctx) error {
	customerID := c.Locals("user_id").(string)

	var r requests.CreateWishlistRequest
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

	wishlist, err := services.CreateWishlistService(customerID, r.ProductID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"success": true, "message": "Produk ditambahkan ke wishlist.", "data": mappers.ToCreateWishlistResponse(*wishlist)})
}

func DeleteWishlistHandler(c *fiber.Ctx) error {
	customerID := c.Locals("user_id").(string)
	productID := c.Params("product_id")

	if err := services.DeleteWishlistService(customerID, productID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Produk dihapus dari wishlist.", "data": nil})
}

func SetPriceAlertHandler(c *fiber.Ctx) error {
	customerID := c.Locals("user_id").(string)
	productID := c.Params("product_id")

	var r requests.PriceAlertRequest
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

	wishlist, err := services.SetPriceAlertService(customerID, productID, r.TargetPrice)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Price Alert berhasil diaktifkan.", "data": mappers.ToPriceAlertResponse(*wishlist)})
}

func DisablePriceAlertHandler(c *fiber.Ctx) error {
	customerID := c.Locals("user_id").(string)
	productID := c.Params("product_id")

	if err := services.DisablePriceAlertService(customerID, productID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Price Alert berhasil dinonaktifkan.", "data": nil})
}

func SetRestockAlertHandler(c *fiber.Ctx) error {
	customerID := c.Locals("user_id").(string)
	productID := c.Params("product_id")

	wishlist, err := services.SetRestockAlertService(customerID, productID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Restock Alert berhasil diaktifkan.", "data": mappers.ToRestockAlertResponse(*wishlist)})
}

func DisableRestockAlertHandler(c *fiber.Ctx) error {
	customerID := c.Locals("user_id").(string)
	productID := c.Params("product_id")

	if err := services.DisableRestockAlertService(customerID, productID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Restock Alert berhasil dinonaktifkan.", "data": nil})
}
