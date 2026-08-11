package services

import (
	"encoding/json"
	"errors"
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
		return nil, errors.New("gagal menyambung server")
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
		Kondisi:         r.Kondisi,
		Deskripsi:       r.Deskripsi,
		Harga:           float64(r.Harga),
		Stok:            r.Stok,
		StatusPublikasi: r.StatusPublikasi,
		ConditionScore:  r.ConditionScore,
		UkuranTersedia:  ukuranTersedia,
	}

	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal membuat product")
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal menyimpan data")
	}

	return &product, nil
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

func UpdateProductService(userID string, productID string, r requests.CreateProduct) (*models.Product, error) {
	var sellerID models.Seller
	tx := database.DB.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return nil, errors.New("gagal menyambung server")
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

	if r.BrandID != "" {
		var brandID models.Brand
		if err := tx.Select("brand_id").Where("brand_id = ? ", r.BrandID).First(&brandID).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("brand tidak ditemukan")
		}
		product.BrandID = r.BrandID
	}
	if r.CategoryID != "" {
		var CategoryID models.Category
		if err := tx.Select("category_id").Where("category_id = ? ", r.CategoryID).First(&CategoryID).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("category tidak ditemukan")
		}
		product.CategoryID = r.CategoryID
	}

	if r.Kondisi != "" {
		product.Kondisi = r.Kondisi
	}
	if r.Stok <= 0 {
		product.Stok = r.Stok
	}

	if r.Deskripsi != nil {
		product.Deskripsi = r.Deskripsi
	}

	if r.NamaProduk != "" {
		product.NamaProduk = r.NamaProduk
	}

	if r.Harga <= 0 {
		product.Harga = float64(r.Harga)
	}

	if r.StatusPublikasi != "" {
		product.StatusPublikasi = r.StatusPublikasi
	}
	product.ConditionScore = r.ConditionScore

	if r.UkuranTersedia != nil {
		ukuranTersedia, err := json.Marshal(r.UkuranTersedia)
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
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal menyimpan data")
	}
	return &product, nil
}
