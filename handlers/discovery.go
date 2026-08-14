package handlers

import (
	"strconv"

	"github.com/ehanz12/api-SneakHub/requests"
	"github.com/ehanz12/api-SneakHub/responses"
	"github.com/ehanz12/api-SneakHub/services"
	"github.com/gofiber/fiber/v2"
)

func SmartFilterHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var r requests.SmartFilterRequest
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

	items, err := services.SmartFilterService(userID, r)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	data := responses.SmartFilterDataResponse{Items: items}
	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Smart Filter berhasil diproses.", "data": data})
}

func HomePersonalizedHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	sections, err := services.HomePersonalizedService(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	data := responses.HomePersonalizedDataResponse{Sections: sections}
	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Homepage personalized berhasil diambil.", "data": data})
}

func PersonalizedRecommendationHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if limit < 1 {
		limit = 10
	}

	items, err := services.PersonalizedRecommendationService(userID, limit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	data := responses.RecommendationDataResponse{Items: items}
	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Rekomendasi berhasil diambil.", "data": data})
}

func TrendingHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if limit < 1 {
		limit = 10
	}

	items, period, err := services.TrendingService(userID, c.Query("period"), limit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	data := responses.TrendingDataResponse{Period: period, Items: items}
	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Trending shoes berhasil diambil.", "data": data})
}

func BestSellerWeeklyHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	if limit < 1 {
		limit = 10
	}

	items, periodStart, periodEnd, err := services.BestSellerWeeklyService(userID, limit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	data := responses.BestSellerDataResponse{
		PeriodStart: periodStart.Format("2006-01-02"),
		PeriodEnd:   periodEnd.Format("2006-01-02"),
		Items:       items,
	}
	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Weekly best seller berhasil diambil.", "data": data})
}
