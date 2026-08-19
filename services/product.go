package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ehanz12/api-SneakHub/database"
	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/requests"
	"gorm.io/gorm"
)

func CreateProductService(userID string, r requests.CreateProduct) (*models.Product, error) {
	var sellerID models.Seller
	tx := database.DB.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return nil, errors.New("gagal koneksi ke server")
	}
	err := tx.Select("seller_id", "user_id", "status_verifikasi").Where("user_id = ? AND status_verifikasi = ?", userID, "verified").First(&sellerID).Error
	if err != nil {
		tx.Rollback()
		return nil, errors.New("user bukan Seller")
	}

	var brandID models.Brand
	if err := tx.Select("brand_id").Where("brand_id = ? ", r.BrandID).First(&brandID).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("brand tidak ditemukan")
	}
	var CategoryID models.Category
	if err := tx.Select("category_id").Where("category_id = ? ", r.CategoryID).First(&CategoryID).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("category tidak ditemukan")
	}

	ukuranTersedia, err := json.Marshal(r.UkuranTersedia)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("ukuran tersedia tidak valid")
	}
	product := models.Product{
		SellerID:        sellerID.SellerID,
		NamaProduk:      r.NamaProduk,
		BrandID:         r.BrandID,
		CategoryID:      r.CategoryID,
		Kondisi:         normalizeKondisi(r.Kondisi),
		Deskripsi:       r.Deskripsi,
		Harga:           float64(r.Harga),
		Stok:            r.Stok,
		Berat:           r.Berat,
		StatusPublikasi: normalizeProductStatus(r.StatusPublikasi),
		ConditionScore:  r.ConditionScore,
		UkuranTersedia:  ukuranTersedia,
	}

	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal membuat product")
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New("data produk gagal disimpan")
	}

	return &product, nil
}

// RatingSummary adalah ringkasan rating sebuah produk.
type RatingSummary struct {
	AvgRating   float64
	TotalReview int64
}

// FirstImageURL mengembalikan URL gambar pertama (urutan tampil terkecil)
// dari sebuah produk, atau string kosong bila tidak ada gambar.
func FirstImageURL(images []models.ProductImage) string {
	if len(images) == 0 {
		return ""
	}
	return images[0].URLObjectStorage
}

// GetRatingSummaries menghitung rata-rata rating dan jumlah review untuk
// daftar product_id dalam satu query GROUP BY.
func GetRatingSummaries(productIDs []string) map[string]RatingSummary {
	summaries := make(map[string]RatingSummary)
	if len(productIDs) == 0 {
		return summaries
	}

	var rows []struct {
		ProductID   string  `gorm:"column:product_id"`
		AvgRating   float64 `gorm:"column:avg_rating"`
		TotalReview int64   `gorm:"column:total_review"`
	}
	if err := database.DB.Model(&models.Review{}).
		Select("product_id, AVG(rating) AS avg_rating, COUNT(*) AS total_review").
		Where("product_id IN ?", productIDs).
		Group("product_id").
		Scan(&rows).Error; err != nil {
		return summaries
	}
	for _, row := range rows {
		summaries[row.ProductID] = RatingSummary{AvgRating: row.AvgRating, TotalReview: row.TotalReview}
	}
	return summaries
}

func GetProductsService(page, limit int, search, brandID, categoryID, kondisi string, minPrice, maxPrice float64, size, sort string) ([]models.Product, int64, error) {
	query := database.DB.Model(&models.Product{}).Where("status_publikasi = ?", "aktif")

	if search != "" {
		query = query.Where("nama_produk LIKE ?", "%"+search+"%")
	}
	if brandID != "" {
		query = query.Where("brand_id = ?", brandID)
	}
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if kondisi != "" {
		query = query.Where("kondisi = ?", strings.ToLower(kondisi))
	}
	if minPrice > 0 {
		query = query.Where("harga >= ?", minPrice)
	}
	if maxPrice > 0 {
		query = query.Where("harga <= ?", maxPrice)
	}
	if size != "" {
		query = query.Where("JSON_SEARCH(ukuran_tersedia, 'one', ?) IS NOT NULL", size)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.New("gagal menghitung product")
	}

	switch sort {
	case "price_asc":
		query = query.Order("harga asc")
	case "price_desc":
		query = query.Order("harga desc")
	case "name_asc":
		query = query.Order("nama_produk asc")
	case "name_desc":
		query = query.Order("nama_produk desc")
	case "oldest":
		query = query.Order("created_at asc")
	default:
		query = query.Order("created_at desc")
	}

	query = query.Preload("Seller").Preload("Images", func(db *gorm.DB) *gorm.DB {
		return db.Order("urutan_tampil asc")
	})

	var products []models.Product
	if err := query.Offset((page - 1) * limit).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, errors.New("gagal memuat product")
	}
	return products, total, nil
}

func GetProductByIDService(productID string) (*models.Product, error) {
	var product models.Product
	err := database.DB.
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("urutan_tampil asc")
		}).
		Where("product_id = ? AND status_publikasi = ?", productID, "aktif").
		First(&product).Error
	if err != nil {
		return nil, errors.New("product tidak ditemukan")
	}
	return &product, nil
}

func UpdateProductService(userID string, productID string, r requests.UpdateProductRequest) (*models.Product, error) {
	var sellerID models.Seller
	tx := database.DB.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return nil, errors.New("gagal koneksi ke server")
	}
	err := tx.Select("seller_id", "user_id", "status_verifikasi").Where("user_id = ? AND status_verifikasi = ?", userID, "verified").First(&sellerID).Error
	if err != nil {
		tx.Rollback()
		return nil, errors.New("user bukan Seller")
	}

	var product models.Product
	if err := tx.Where("product_id = ? AND seller_id = ?", productID, sellerID.SellerID).First(&product).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("product tidak ditemukan")
	}
	oldPrice := product.Harga
	oldStok := product.Stok

	if r.BrandID != nil {
		var brandID models.Brand
		if err := tx.Select("brand_id").Where("brand_id = ? ", *r.BrandID).First(&brandID).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("brand tidak ditemukan")
		}
		product.BrandID = *r.BrandID
	}
	if r.CategoryID != nil {
		var CategoryID models.Category
		if err := tx.Select("category_id").Where("category_id = ? ", *r.CategoryID).First(&CategoryID).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("category tidak ditemukan")
		}
		product.CategoryID = *r.CategoryID
	}

	if r.Kondisi != nil {
		product.Kondisi = normalizeKondisi(*r.Kondisi)
		if product.Kondisi == "" {
			tx.Rollback()
			return nil, errors.New("kondisi tidak valid, gunakan NEW, USED, atau REFURBISHED")
		}
	}

	if r.Deskripsi != nil {
		product.Deskripsi = r.Deskripsi
	}

	if r.NamaProduk != nil {
		product.NamaProduk = *r.NamaProduk
	}

	if r.Harga != nil {
		product.Harga = *r.Harga
	}

	if r.Stok != nil {
		product.Stok = *r.Stok
	}

	if r.Berat != nil {
		product.Berat = *r.Berat
	}

	if r.StatusPublikasi != nil {
		product.StatusPublikasi = normalizeProductStatus(*r.StatusPublikasi)
	}

	if r.ConditionScore != nil {
		product.ConditionScore = r.ConditionScore
	}

	if r.UkuranTersedia != nil {
		ukuranTersedia, err := json.Marshal(*r.UkuranTersedia)
		if err != nil {
			tx.Rollback()
			return nil, errors.New("ukuran tersedia tidak valid")
		}
		product.UkuranTersedia = ukuranTersedia
	}

	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal update product")
	}

	priceChanged := r.Harga != nil && *r.Harga != oldPrice
	restocked := r.Stok != nil && oldStok == 0 && *r.Stok > 0
	if priceChanged || restocked {
		if err := triggerProductAlerts(tx, product, priceChanged, restocked); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New("data produk gagal disimpan")
	}
	return &product, nil
}

// triggerProductAlerts membangkitkan notifikasi price_alert/restock_alert
// kepada user yang memasang alert pada produk, lalu menonaktifkan alert-nya.
func triggerProductAlerts(tx *gorm.DB, product models.Product, priceChanged, restocked bool) error {
	if priceChanged {
		var wishlists []models.Wishlist
		if err := tx.Where("product_id = ? AND price_alert_enabled = ? AND target_price >= ?",
			product.ProductID, true, product.Harga).Find(&wishlists).Error; err != nil {
			return errors.New("gagal memuat wishlist")
		}
		for _, w := range wishlists {
			isi := fmt.Sprintf("Harga produk %s turun ke Rp %.0f, sesuai target kamu.", product.NamaProduk, product.Harga)
			if err := createNotification(tx, w.CustomerID, "price_alert", isi); err != nil {
				return err
			}
			if err := tx.Model(&models.Wishlist{}).Where("wishlist_id = ?", w.WishlistID).
				Updates(map[string]interface{}{"price_alert_enabled": false, "target_price": nil}).Error; err != nil {
				return errors.New("gagal menonaktifkan price alert")
			}
		}
	}

	if restocked {
		var wishlists []models.Wishlist
		if err := tx.Where("product_id = ? AND restock_alert_enabled = ?",
			product.ProductID, true).Find(&wishlists).Error; err != nil {
			return errors.New("gagal memuat wishlist")
		}
		for _, w := range wishlists {
			isi := fmt.Sprintf("Produk %s sudah tersedia kembali (stok %d).", product.NamaProduk, product.Stok)
			if err := createNotification(tx, w.CustomerID, "restock_alert", isi); err != nil {
				return err
			}
			if err := tx.Model(&models.Wishlist{}).Where("wishlist_id = ?", w.WishlistID).
				Update("restock_alert_enabled", false).Error; err != nil {
				return errors.New("gagal menonaktifkan restock alert")
			}
		}
	}

	return nil
}

// DeleteProductService menghapus produk beserta data terkait.
// Diblokir jika produk masih memiliki pesanan aktif.
func DeleteProductService(userID, productID string) error {
	tx := database.DB.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return errors.New("gagal koneksi ke server")
	}

	var sellerID models.Seller
	err := tx.Select("seller_id", "user_id", "status_verifikasi").Where("user_id = ? AND status_verifikasi = ?", userID, "verified").First(&sellerID).Error
	if err != nil {
		tx.Rollback()
		return errors.New("user bukan Seller")
	}

	var product models.Product
	if err := tx.Where("product_id = ? AND seller_id = ?", productID, sellerID.SellerID).First(&product).Error; err != nil {
		tx.Rollback()
		return errors.New("product tidak ditemukan")
	}

	var activeOrderCount int64
	if err := tx.Model(&models.OrderItem{}).
		Joins("JOIN orders ON orders.order_id = order_items.order_id").
		Where("order_items.product_id = ? AND LOWER(orders.status_order) IN ?", productID, []string{"pending", "diproses", "dikirim"}).
		Count(&activeOrderCount).Error; err != nil {
		tx.Rollback()
		return errors.New("gagal memeriksa pesanan aktif")
	}
	if activeOrderCount > 0 {
		tx.Rollback()
		return errors.New("produk masih memiliki pesanan aktif")
	}

	var images []models.ProductImage
	if err := tx.Where("product_id = ?", productID).Find(&images).Error; err != nil {
		tx.Rollback()
		return errors.New("gagal memuat gambar produk")
	}

	if err := tx.Where("product_id = ?", productID).Delete(&models.ImageEmbedding{}).Error; err != nil {
		tx.Rollback()
		return errors.New("gagal menghapus data terkait produk")
	}
	if err := tx.Where("product_id = ?", productID).Delete(&models.ProductImage{}).Error; err != nil {
		tx.Rollback()
		return errors.New("gagal menghapus data terkait produk")
	}
	if err := tx.Where("product_id = ?", productID).Delete(&models.ConditionScore{}).Error; err != nil {
		tx.Rollback()
		return errors.New("gagal menghapus data terkait produk")
	}
	if err := tx.Where("product_id = ?", productID).Delete(&models.PriceHistory{}).Error; err != nil {
		tx.Rollback()
		return errors.New("gagal menghapus data terkait produk")
	}
	if err := tx.Model(&models.RecommendationData{}).
		Where("JSON_SEARCH(daftar_product_id, 'one', ?) IS NOT NULL", productID).
		Update("daftar_product_id", gorm.Expr("JSON_REMOVE(daftar_product_id, JSON_UNQUOTE(JSON_SEARCH(daftar_product_id, 'one', ?)))", productID)).Error; err != nil {
		tx.Rollback()
		return errors.New("gagal menghapus data terkait produk")
	}
	if err := tx.Where("product_id = ?", productID).Delete(&models.Wishlist{}).Error; err != nil {
		tx.Rollback()
		return errors.New("gagal menghapus data terkait produk")
	}
	if err := tx.Where("product_id = ?", productID).Delete(&models.CartItem{}).Error; err != nil {
		tx.Rollback()
		return errors.New("gagal menghapus data terkait produk")
	}

	if err := tx.Delete(&product).Error; err != nil {
		tx.Rollback()
		return errors.New("gagal menghapus product")
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errors.New("gagal menghapus product")
	}

	for _, image := range images {
		if fileName := filepath.Base(image.URLObjectStorage); fileName != "." && fileName != "/" {
			os.Remove(filepath.Join(UploadDir, fileName))
		}
	}

	return nil
}
