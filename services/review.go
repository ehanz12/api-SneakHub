package services

import (
	"errors"

	"github.com/ehanz12/api-SneakHub/database"
	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/requests"
	"gorm.io/gorm"
)

func CreateReviewService(customerID, orderID string, r requests.CreateReviewRequest) (*models.Review, error) {
	tx := database.DB.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return nil, errors.New("gagal menyambung server")
	}

	var order models.Order
	if err := tx.Where("order_id = ? AND customer_id = ?", orderID, customerID).First(&order).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order tidak ditemukan")
		}
		return nil, errors.New("gagal memuat order")
	}

	if order.StatusOrder != "selesai" {
		tx.Rollback()
		return nil, errors.New("order hanya bisa direview setelah status selesai")
	}

	var orderItem models.OrderItem
	if err := tx.Where("order_id = ? AND product_id = ?", orderID, r.ProductID).First(&orderItem).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("produk tidak ada di dalam order ini")
	}

	var existing int64
	if err := tx.Model(&models.Review{}).
		Where("order_id = ? AND product_id = ?", orderID, r.ProductID).
		Count(&existing).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal memeriksa review")
	}
	if existing > 0 {
		tx.Rollback()
		return nil, errors.New("review untuk produk ini sudah ada")
	}

	review := models.Review{
		OrderID:    orderID,
		CustomerID: customerID,
		ProductID:  r.ProductID,
		Rating:     r.Rating,
		Komentar:   r.Komentar,
	}
	if err := tx.Create(&review).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal membuat review")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal menyimpan data")
	}

	return &review, nil
}
