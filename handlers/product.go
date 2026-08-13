package handlers

import (
	"strconv"

	"github.com/ehanz12/api-SneakHub/mappers"
	"github.com/ehanz12/api-SneakHub/requests"
	"github.com/ehanz12/api-SneakHub/responses"
	"github.com/ehanz12/api-SneakHub/services"
	"github.com/gofiber/fiber/v2"
)

func GetProductsHandler(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	minPrice, _ := strconv.ParseFloat(c.Query("min_price"), 64)
	maxPrice, _ := strconv.ParseFloat(c.Query("max_price"), 64)

	products, total, err := services.GetProductsService(
		page,
		limit,
		c.Query("search"),
		c.Query("brand_id"),
		c.Query("category_id"),
		c.Query("kondisi"),
		minPrice,
		maxPrice,
		c.Query("size"),
		c.Query("sort"),
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "terjadi kesalahan server", "errors": err.Error()})
	}

	totalPages := 0
	if total > 0 {
		totalPages = (int(total) + limit - 1) / limit
	}

	data := responses.ProductListDataResponse{
		Items: mappers.ToProductListResponse(products),
		Pagination: responses.PaginationResponse{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}
	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Produk berhasil diambil.", "data": data})
}

func GetProductByIDHandler(c *fiber.Ctx) error {
	productID := c.Params("product_id")
	product, err := services.GetProductByIDService(productID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Produk tidak ditemukan."})
	}
	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Detail produk berhasil diambil.", "data": mappers.ToProductDetailResponse(*product)})
}

func CreateProductHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	var r requests.CreateProduct
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
	product, err := services.CreateProductService(userID, r)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "request gagal", "errors": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"success": true, "message": "product berhasil dibuat", "data": mappers.ToProductResponse(*product)})
}

func UpdateProductHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	productID := c.Params("product_id")
	var r requests.CreateProduct
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
	product, err := services.UpdateProductService(userID, productID, r)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "request gagal", "errors": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"success": true, "message": "product berhasil di update", "data": mappers.ToProductUpdateResponse(*product)})
}

func DeleteProductHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	productID := c.Params("product_id")

	if err := services.DeleteProductService(userID, productID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Produk berhasil dihapus.",
		"data":    nil,
	})
}
