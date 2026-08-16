package handlers

import (
	"strconv"

	"github.com/ehanz12/api-SneakHub/mappers"
	"github.com/ehanz12/api-SneakHub/requests"
	"github.com/ehanz12/api-SneakHub/responses"
	"github.com/ehanz12/api-SneakHub/services"
	"github.com/gofiber/fiber/v2"
)

func AdminGetUsersHandler(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	users, total, err := services.GetAdminUsersService(page, limit, c.Query("status"), c.Query("role"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	totalPages := 0
	if total > 0 {
		totalPages = (int(total) + limit - 1) / limit
	}

	data := responses.AdminUserListDataResponse{
		Items: mappers.ToAdminUserListResponse(users),
		Pagination: responses.PaginationResponse{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Data pengguna berhasil diambil.", "data": data})
}

func AdminUpdateUserStatusHandler(c *fiber.Ctx) error {
	userID := c.Params("user_id")

	var r requests.UpdateUserStatusRequest
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

	user, err := services.UpdateUserStatusService(userID, r.StatusAkun)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Status akun berhasil diperbarui.", "data": mappers.ToAdminUserStatusResponse(user)})
}

func AdminGetSellersHandler(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	sellers, total, err := services.GetAdminSellersService(page, limit, c.Query("status"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	totalPages := 0
	if total > 0 {
		totalPages = (int(total) + limit - 1) / limit
	}

	data := responses.AdminSellerListDataResponse{
		Items: mappers.ToAdminSellerListResponse(sellers),
		Pagination: responses.PaginationResponse{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Data toko seller berhasil diambil.", "data": data})
}

func AdminVerifySellerHandler(c *fiber.Ctx) error {
	sellerID := c.Params("seller_id")

	var r requests.VerifySellerRequest
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

	seller, err := services.VerifySellerService(sellerID, r.Status)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Status verifikasi toko berhasil diperbarui.", "data": mappers.ToAdminSellerVerificationResponse(seller)})
}

func AdminUpdateUserRoleHandler(c *fiber.Ctx) error {
	userID := c.Params("user_id")

	var r requests.UpdateUserRoleRequest
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

	user, err := services.UpdateUserRoleService(userID, r.Peran)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Peran user berhasil diperbarui.", "data": mappers.ToAdminUserRoleResponse(user)})
}

func AdminGetProductsHandler(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	products, total, err := services.GetAdminProductsService(page, limit, c.Query("status"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	totalPages := 0
	if total > 0 {
		totalPages = (int(total) + limit - 1) / limit
	}

	data := responses.AdminProductListDataResponse{
		Items: mappers.ToAdminProductListResponse(products),
		Pagination: responses.PaginationResponse{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Data produk berhasil diambil.", "data": data})
}

func AdminUpdateProductStatusHandler(c *fiber.Ctx) error {
	productID := c.Params("product_id")

	var r requests.UpdateProductStatusRequest
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

	product, err := services.UpdateProductStatusService(productID, r.StatusPublikasi)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Status produk berhasil diperbarui.", "data": mappers.ToAdminProductStatusResponse(product)})
}

func AdminGetOrdersHandler(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	orders, total, err := services.GetAdminOrdersService(page, limit, c.Query("status"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	totalPages := 0
	if total > 0 {
		totalPages = (int(total) + limit - 1) / limit
	}

	data := responses.AdminOrderListDataResponse{
		Items: mappers.ToAdminOrderListResponse(orders),
		Pagination: responses.PaginationResponse{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Data order berhasil diambil.", "data": data})
}

func AdminGetReportsHandler(c *fiber.Ctx) error {
	report, err := services.GetAdminReportsService(c.Query("period"), c.Query("start_date"), c.Query("end_date"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Laporan berhasil diambil.", "data": mappers.ToAdminReportResponse(report)})
}
