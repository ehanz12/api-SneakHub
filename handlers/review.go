package handlers

import (
	"strconv"

	"github.com/ehanz12/api-SneakHub/mappers"
	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/requests"
	"github.com/ehanz12/api-SneakHub/responses"
	"github.com/ehanz12/api-SneakHub/services"
	"github.com/gofiber/fiber/v2"
)

func CreateReviewHandler(c *fiber.Ctx) error {
	customerID := c.Locals("user_id").(string)
	orderID := c.Params("order_id")

	var r requests.CreateReviewRequest
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

	review, err := services.CreateReviewService(customerID, orderID, r)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"success": true, "message": "Review berhasil dibuat.", "data": mappers.ToReviewResponse(*review)})
}

func buildReviewListData(reviews []models.Review, avg float64, total int64, page, limit int) responses.ReviewListDataResponse {
	totalPages := 0
	if total > 0 {
		totalPages = (int(total) + limit - 1) / limit
	}
	return responses.ReviewListDataResponse{
		Items:          mappers.ToReviewListResponse(reviews),
		RatingRataRata: avg,
		TotalReview:    total,
		Pagination: responses.PaginationResponse{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}
}

func GetProductReviewsHandler(c *fiber.Ctx) error {
	productID := c.Params("product_id")

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	reviews, avg, total, err := services.GetProductReviewsService(productID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Review produk berhasil diambil.",
		"data":    buildReviewListData(reviews, avg, total, page, limit),
	})
}

func GetSellerReviewsHandler(c *fiber.Ctx) error {
	sellerID := c.Params("seller_id")

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	reviews, avg, total, err := services.GetSellerReviewsService(sellerID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Review toko seller berhasil diambil.",
		"data":    buildReviewListData(reviews, avg, total, page, limit),
	})
}
