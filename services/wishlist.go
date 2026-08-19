package services

import (
	"errors"

	"github.com/ehanz12/api-SneakHub/database"
	"github.com/ehanz12/api-SneakHub/models"
	"gorm.io/gorm"
)

// findWishlist mengambil wishlist milik customer untuk sebuah produk.
func findWishlist(customerID, productID string) (*models.Wishlist, error) {
	var wishlist models.Wishlist
	if err := database.DB.Where("customer_id = ? AND product_id = ?", customerID, productID).
		First(&wishlist).Error; err != nil {
		return nil, errors.New("produk tidak ada di wishlist")
	}
	return &wishlist, nil
}

// GetWishlistService mengambil semua wishlist customer beserta info produk,
// sekaligus menyegarkan status stok terakhir dari produk.
func GetWishlistService(customerID string) ([]models.Wishlist, error) {
	var wishlists []models.Wishlist
	if err := database.DB.
		Preload("Product", func(db *gorm.DB) *gorm.DB {
			return db.Select("product_id", "nama_produk", "harga", "stok", "status_publikasi")
		}).
		Preload("Product.Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("urutan_tampil asc")
		}).
		Where("customer_id = ?", customerID).
		Order("created_at desc").
		Find(&wishlists).Error; err != nil {
		return nil, errors.New("gagal memuat wishlist")
	}

	for i := range wishlists {
		expected := "out_of_stock"
		if wishlists[i].Product.Stok > 0 {
			expected = "available"
		}
		if wishlists[i].StatusStok != expected {
			if err := database.DB.Model(&models.Wishlist{}).
				Where("wishlist_id = ?", wishlists[i].WishlistID).
				Update("status_stok", expected).Error; err == nil {
				wishlists[i].StatusStok = expected
			}
		}
	}

	return wishlists, nil
}

// CreateWishlistService menambahkan produk ke wishlist customer.
func CreateWishlistService(customerID, productID string) (*models.Wishlist, error) {
	var product models.Product
	if err := database.DB.Select("product_id", "stok", "status_publikasi").
		Where("product_id = ?", productID).First(&product).Error; err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}
	if product.StatusPublikasi != "aktif" {
		return nil, errors.New("produk tidak aktif")
	}

	statusStok := "out_of_stock"
	if product.Stok > 0 {
		statusStok = "available"
	}

	wishlist := models.Wishlist{
		CustomerID: customerID,
		ProductID:  productID,
		StatusStok: statusStok,
	}
	if err := database.DB.Create(&wishlist).Error; err != nil {
		return nil, errors.New("produk sudah ada di wishlist")
	}
	return &wishlist, nil
}

// DeleteWishlistService menghapus produk dari wishlist customer.
func DeleteWishlistService(customerID, productID string) error {
	res := database.DB.Where("customer_id = ? AND product_id = ?", customerID, productID).
		Delete(&models.Wishlist{})
	if res.Error != nil {
		return errors.New("gagal menghapus wishlist")
	}
	if res.RowsAffected == 0 {
		return errors.New("produk tidak ada di wishlist")
	}
	return nil
}

// SetPriceAlertService mengaktifkan price alert dengan target harga.
func SetPriceAlertService(customerID, productID string, targetPrice *float64) (*models.Wishlist, error) {
	wishlist, err := findWishlist(customerID, productID)
	if err != nil {
		return nil, err
	}
	if err := database.DB.Model(&models.Wishlist{}).
		Where("wishlist_id = ?", wishlist.WishlistID).
		Updates(map[string]interface{}{
			"price_alert_enabled": true,
			"target_price":        targetPrice,
		}).Error; err != nil {
		return nil, errors.New("gagal mengaktifkan price alert")
	}
	wishlist.PriceAlertEnabled = true
	wishlist.TargetPrice = targetPrice
	return wishlist, nil
}

// DisablePriceAlertService menonaktifkan price alert wishlist.
func DisablePriceAlertService(customerID, productID string) error {
	wishlist, err := findWishlist(customerID, productID)
	if err != nil {
		return err
	}
	if err := database.DB.Model(&models.Wishlist{}).
		Where("wishlist_id = ?", wishlist.WishlistID).
		Updates(map[string]interface{}{
			"price_alert_enabled": false,
			"target_price":        nil,
		}).Error; err != nil {
		return errors.New("gagal menonaktifkan price alert")
	}
	return nil
}

// SetRestockAlertService mengaktifkan restock alert wishlist.
func SetRestockAlertService(customerID, productID string) (*models.Wishlist, error) {
	wishlist, err := findWishlist(customerID, productID)
	if err != nil {
		return nil, err
	}
	if err := database.DB.Model(&models.Wishlist{}).
		Where("wishlist_id = ?", wishlist.WishlistID).
		Update("restock_alert_enabled", true).Error; err != nil {
		return nil, errors.New("gagal mengaktifkan restock alert")
	}
	wishlist.RestockAlertEnabled = true
	return wishlist, nil
}

// DisableRestockAlertService menonaktifkan restock alert wishlist.
func DisableRestockAlertService(customerID, productID string) error {
	wishlist, err := findWishlist(customerID, productID)
	if err != nil {
		return err
	}
	if err := database.DB.Model(&models.Wishlist{}).
		Where("wishlist_id = ?", wishlist.WishlistID).
		Update("restock_alert_enabled", false).Error; err != nil {
		return errors.New("gagal menonaktifkan restock alert")
	}
	return nil
}
