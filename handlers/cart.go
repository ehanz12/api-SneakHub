package handlers

import (
	"github.com/ehanz12/api-SneakHub/mappers"
	"github.com/ehanz12/api-SneakHub/requests"
	"github.com/ehanz12/api-SneakHub/responses"
	"github.com/ehanz12/api-SneakHub/services"
	"github.com/gofiber/fiber/v2"
)

func GetCartHandler(c *fiber.Ctx) error {
	customerID := c.Locals("user_id").(string)

	cart, err := services.GetCartService(customerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Cart berhasil diambil.",
		"data":    mappers.ToCartResponse(*cart),
	})
}

func AddCartItemsHandler(c *fiber.Ctx) error {
	customerID := c.Locals("user_id").(string)

	var r requests.AddCartItemsRequest
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

	items, total, err := services.AddCartItemsService(customerID, r)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	data := responses.CartItemsResponse{
		Items: mappers.ToCartItemListResponse(items),
		Total: total,
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Produk ditambahkan ke cart.",
		"data":    data,
	})
}

func UpdateCartItemHandler(c *fiber.Ctx) error {
	customerID := c.Locals("user_id").(string)
	cartItemID := c.Params("cart_item_id")

	var r requests.UpdateCartItemRequest
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

	item, err := services.UpdateCartItemService(customerID, cartItemID, r.Jumlah)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Jumlah item berhasil diperbarui.",
		"data":    mappers.ToCartItemResponse(*item),
	})
}

func DeleteCartItemHandler(c *fiber.Ctx) error {
	customerID := c.Locals("user_id").(string)
	cartItemID := c.Params("cart_item_id")

	if err := services.DeleteCartItemService(customerID, cartItemID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Item cart berhasil dihapus.",
		"data":    nil,
	})
}
