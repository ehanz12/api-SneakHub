package handlers

import (
	"strconv"

	"github.com/ehanz12/api-SneakHub/mappers"
	"github.com/ehanz12/api-SneakHub/responses"
	"github.com/ehanz12/api-SneakHub/services"
	"github.com/gofiber/fiber/v2"
)

func GetNotificationsHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	notifications, unread, err := services.GetNotificationsService(userID, page, limit, c.Query("type"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	data := responses.NotificationListDataResponse{
		Items:       mappers.ToNotificationListResponse(notifications),
		UnreadCount: unread,
	}
	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Notifikasi berhasil diambil.", "data": data})
}

func MarkNotificationReadHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	notificationID := c.Params("notification_id")

	notification, err := services.MarkNotificationReadService(userID, notificationID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"success": true, "message": "Notifikasi ditandai sudah dibaca.", "data": mappers.ToNotificationReadResponse(*notification)})
}
