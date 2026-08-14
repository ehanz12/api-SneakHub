package handlers

import (
	"strconv"

	"github.com/ehanz12/api-SneakHub/mappers"
	"github.com/ehanz12/api-SneakHub/requests"
	"github.com/ehanz12/api-SneakHub/responses"
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

func GetSellerProductsHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	products, total, err := services.GetSellerProductsService(userID, page, limit, c.Query("status"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	totalPages := 0
	if total > 0 {
		totalPages = (int(total) + limit - 1) / limit
	}

	data := responses.SellerProductListDataResponse{
		Items: mappers.ToSellerProductListResponse(products),
		Pagination: responses.PaginationResponse{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Produk seller berhasil diambil.", "data": data})
}

func GetSellerOrdersHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	orders, total, err := services.GetSellerOrdersService(userID, page, limit, c.Query("status"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	totalPages := 0
	if total > 0 {
		totalPages = (int(total) + limit - 1) / limit
	}

	data := responses.SellerOrderListDataResponse{
		Items: mappers.ToSellerOrderListResponse(orders),
		Pagination: responses.PaginationResponse{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Order seller berhasil diambil.", "data": data})
}

func GetSellerDashboardHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	dashboard, err := services.GetSellerDashboardService(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Dashboard seller berhasil diambil.", "data": mappers.ToSellerDashboardResponse(dashboard)})
}
