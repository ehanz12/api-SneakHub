package handlers

import (
	"strconv"

	"github.com/ehanz12/api-SneakHub/mappers"
	"github.com/ehanz12/api-SneakHub/requests"
	"github.com/ehanz12/api-SneakHub/responses"
	"github.com/ehanz12/api-SneakHub/services"
	"github.com/gofiber/fiber/v2"
)

func GetOrdersHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	orders, total, err := services.GetOrdersService(userID, role, page, limit, c.Query("status"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	totalPages := 0
	if total > 0 {
		totalPages = (int(total) + limit - 1) / limit
	}

	data := responses.OrderListDataResponse{
		Items: mappers.ToOrderListResponse(orders),
		Pagination: responses.PaginationResponse{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Order berhasil diambil.", "data": data})
}

func GetOrderHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)
	orderID := c.Params("order_id")

	order, err := services.GetOrderService(userID, role, orderID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Detail order berhasil diambil.", "data": mappers.ToOrderDetailResponse(*order)})
}

func CreateOrderHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	var r requests.CreateOrderRequest
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

	order, err := services.CreateOrderService(userID, role, r)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"success": true, "message": "Order berhasil dibuat.", "data": mappers.ToOrderDetailResponse(*order)})
}

func UpdateOrderHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)
	orderID := c.Params("order_id")

	var r requests.UpdateOrderStatusRequest
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

	order, err := services.UpdateOrderStatusService(userID, role, orderID, *r.StatusOrder)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Status order berhasil diperbarui.", "data": mappers.ToOrderDetailResponse(*order)})
}

func DeleteOrderHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)
	orderID := c.Params("order_id")

	if err := services.DeleteOrderService(userID, role, orderID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Order berhasil dihapus.", "data": nil})
}
