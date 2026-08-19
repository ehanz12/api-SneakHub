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
		return nil, errors.New("koneksi database gagal")
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
		return nil, errors.New("review gagal tersimpan")
	}

	return &review, nil
}

func reviewQueryWithCustomer(db *gorm.DB) *gorm.DB {
	return db.Preload("Customer", func(db *gorm.DB) *gorm.DB {
		return db.Select("user_id", "nama")
	})
}

func GetProductReviewsService(productID string, page, limit int) ([]models.Review, float64, int64, error) {
	var total int64
	if err := database.DB.Model(&models.Review{}).
		Where("product_id = ?", productID).Count(&total).Error; err != nil {
		return nil, 0, 0, errors.New("gagal menghitung review")
	}

	var avg float64
	if err := database.DB.Model(&models.Review{}).
		Where("product_id = ?", productID).
		Select("COALESCE(AVG(rating), 0)").Scan(&avg).Error; err != nil {
		return nil, 0, 0, errors.New("gagal menghitung rating")
	}

	var reviews []models.Review
	if err := reviewQueryWithCustomer(database.DB).
		Where("product_id = ?", productID).
		Order("created_at desc").
		Offset((page - 1) * limit).Limit(limit).
		Find(&reviews).Error; err != nil {
		return nil, 0, 0, errors.New("gagal memuat review")
	}

	return reviews, avg, total, nil
}

func GetSellerReviewsService(sellerID string, page, limit int) ([]models.Review, float64, int64, error) {
	var sellerCheck int64
	if err := database.DB.Model(&models.Seller{}).
		Where("seller_id = ?", sellerID).Count(&sellerCheck).Error; err != nil {
		return nil, 0, 0, errors.New("gagal memvalidasi toko")
	}
	if sellerCheck == 0 {
		return nil, 0, 0, errors.New("toko seller tidak ditemukan")
	}

	var total int64
	if err := database.DB.Model(&models.Review{}).
		Joins("JOIN products ON products.product_id = reviews.product_id").
		Where("products.seller_id = ?", sellerID).Count(&total).Error; err != nil {
		return nil, 0, 0, errors.New("gagal menghitung review")
	}

	var avg float64
	if err := database.DB.Model(&models.Review{}).
		Joins("JOIN products ON products.product_id = reviews.product_id").
		Where("products.seller_id = ?", sellerID).
		Select("COALESCE(AVG(reviews.rating), 0)").Scan(&avg).Error; err != nil {
		return nil, 0, 0, errors.New("gagal menghitung rating")
	}

	var reviews []models.Review
	if err := reviewQueryWithCustomer(database.DB).
		Joins("JOIN products ON products.product_id = reviews.product_id").
		Where("products.seller_id = ?", sellerID).
		Order("reviews.created_at desc").
		Offset((page - 1) * limit).Limit(limit).
		Find(&reviews).Error; err != nil {
		return nil, 0, 0, errors.New("gagal memuat review")
	}

	return reviews, avg, total, nil
}
