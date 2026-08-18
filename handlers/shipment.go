package handlers

import (
	"github.com/ehanz12/api-SneakHub/mappers"
	"github.com/ehanz12/api-SneakHub/requests"
	"github.com/ehanz12/api-SneakHub/services"
	"github.com/gofiber/fiber/v2"
)

// ShipOrderHandler (seller) mengirim pesanan: booking kurir via Biteship
// untuk mendapatkan resi otomatis, atau resi manual jika dikirim di body.
func ShipOrderHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	orderID := c.Params("order_id")

	var r requests.ShipOrderRequest
	if err := c.BodyParser(&r); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "request gagal", "errors": err.Error()})
	}

	order, err := services.ShipOrderService(userID, orderID, r)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Pesanan berhasil dikirim.", "data": mappers.ToOrderDetailResponse(*order)})
}

// ConfirmReceivedHandler (customer) menandai pesanan selesai setelah
// barang diterima.
func ConfirmReceivedHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	orderID := c.Params("order_id")

	order, err := services.ConfirmReceivedService(userID, orderID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Pesanan selesai.", "data": mappers.ToOrderDetailResponse(*order)})
}
