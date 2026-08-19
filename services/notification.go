package services

import (
	"errors"
	"strings"

	"github.com/ehanz12/api-SneakHub/database"
	"github.com/ehanz12/api-SneakHub/models"
	"gorm.io/gorm"
)

func createNotification(db *gorm.DB, userID, jenis, isi string) error {
	notification := models.Notification{
		UserID:          userID,
		JenisNotifikasi: jenis,
		IsiNotifikasi:   isi,
	}
	if err := db.Create(&notification).Error; err != nil {
		return errors.New("gagal membuat notifikasi")
	}
	return nil
}

func normalizeNotificationType(jenis string) string {
	switch strings.ToUpper(strings.TrimSpace(jenis)) {
	case "PRICE_ALERT":
		return "price_alert"
	case "RESTOCK_ALERT":
		return "restock_alert"
	case "ORDER_UPDATE":
		return "order_update"
	case "PROMO":
		return "promo"
	case "DLL", "OTHER":
		return "dll"
	}
	return ""
}

func GetNotificationsService(userID string, page, limit int, jenis string) ([]models.Notification, int64, error) {
	query := database.DB.Model(&models.Notification{}).Where("user_id = ?", userID)

	if j := normalizeNotificationType(jenis); j != "" {
		query = query.Where("jenis_notifikasi = ?", j)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.New("gagal memuat notifikasi")
	}

	var notifications []models.Notification
	if err := query.Order("created_at desc").
		Offset((page - 1) * limit).Limit(limit).Find(&notifications).Error; err != nil {
		return nil, 0, errors.New("gagal memuat notifikasi")
	}

	var unread int64
	if err := database.DB.Model(&models.Notification{}).
		Where("user_id = ? AND status_baca = ?", userID, false).
		Count(&unread).Error; err != nil {
		return nil, 0, errors.New("gagal menghitung notifikasi belum dibaca")
	}

	return notifications, unread, nil
}

func MarkNotificationReadService(userID, notificationID string) (*models.Notification, error) {
	var notification models.Notification
	if err := database.DB.Where("notification_id = ? AND user_id = ?", notificationID, userID).
		First(&notification).Error; err != nil {
		return nil, errors.New("notifikasi tidak ditemukan")
	}

	if err := database.DB.Model(&models.Notification{}).
		Where("notification_id = ?", notificationID).
		Update("status_baca", true).Error; err != nil {
		return nil, errors.New("gagal menandai notifikasi")
	}
	notification.StatusBaca = true
	return &notification, nil
}
