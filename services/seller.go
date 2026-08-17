package services

import (
	"errors"
	"strings"

	"github.com/ehanz12/api-SneakHub/database"
	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/requests"
	"gorm.io/gorm"
)

// normalizeProductStatus memetakan alias status publikasi (mis. ACTIVE)
// ke nilai enum status_publikasi di database. Mengembalikan "" bila tidak dikenal.
func normalizeProductStatus(status string) string {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case "ACTIVE", "AKTIF":
		return "aktif"
	case "DRAFT":
		return "draft"
	case "INACTIVE", "NONAKTIF", "TIDAK_AKTIF":
		return "nonaktif"
	case "PENDING":
		return "draft"
	}
	return ""
}

// normalizeUserStatus memetakan alias status akun (mis. SUSPENDED)
// ke nilai enum status_akun di database. Mengembalikan "" bila tidak dikenal.
func normalizeUserStatus(status string) string {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case "ACTIVE", "AKTIF":
		return "aktif"
	case "INACTIVE", "TIDAK_AKTIF":
		return "tidak_aktif"
	case "SUSPENDED", "BLOCKED", "BLOKIR":
		return "blokir"
	}
	return ""
}

// normalizeRole memetakan alias peran (mis. CUSTOMER) ke nilai enum peran
// di database. Mengembalikan "" bila tidak dikenal.
func normalizeRole(role string) string {
	switch strings.ToUpper(strings.TrimSpace(role)) {
	case "CUSTOMER":
		return "customer"
	case "SELLER":
		return "seller"
	case "ADMIN":
		return "admin"
	}
	return ""
}

// normalizeSellerStatus memetakan alias status verifikasi (mis. PENDING)
// ke nilai enum status_verifikasi di database. Mengembalikan "" bila tidak dikenal.
func normalizeSellerStatus(status string) string {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case "PENDING":
		return "pending"
	case "VERIFIED", "VERIFIKASI":
		return "verified"
	case "REJECTED", "DITOLAK":
		return "rejected"
	}
	return ""
}

// ProductWithSales adalah produk beserta jumlah total terjualnya.
type ProductWithSales struct {
	models.Product
	TotalTerjual int64 `gorm:"column:total_terjual"`
}

// SellerTopProduct adalah produk terlaris milik seller.
type SellerTopProduct struct {
	ProductID    string `gorm:"column:product_id"`
	NamaProduk   string `gorm:"column:nama_produk"`
	TotalTerjual int64  `gorm:"column:total_terjual"`
}

// SellerDashboardData adalah hasil agregasi dashboard seller.
type SellerDashboardData struct {
	TotalProduk      int64
	ProdukAktif      int64
	TotalTerjual     int64
	TotalPendapatan  float64
	RatingRataRata   float64
	SellerTrustScore *float64
	ProdukTerlaris   []SellerTopProduct
}

func CreateSellerService(UserID string, req requests.CreateSellerRequest) (*models.Seller, error) {
	tx := database.DB.Begin()
	if tx.Error != nil {
		return nil, errors.New("gagal menyambungkan ke server")
	}
	var user models.User
	err := tx.Select("user_id").Where("user_id = ?", UserID).First(&user).Error
	if err != nil {
		tx.Rollback()
		return nil, errors.New("user tidak ditemukan")
	}
	var exist models.Seller
	if err := tx.Select("seller_id", "status_verifikasi").Where("user_id = ?", UserID).First(&exist).Error; err == nil {
		// Pengajuan sebelumnya ditolak: boleh mengajukan ulang dengan
		// memperbarui nama & deskripsi toko.
		if exist.StatusVerifikasi == "rejected" {
			updates := map[string]interface{}{
				"nama_toko":         req.NamaToko,
				"deskripsi_toko":    req.DeskripsiToko,
				"status_verifikasi": "pending",
			}
			if err := tx.Model(&models.Seller{}).Where("seller_id = ?", exist.SellerID).
				Updates(updates).Error; err != nil {
				tx.Rollback()
				return nil, errors.New("gagal mengajukan ulang seller")
			}
			if err := tx.Commit().Error; err != nil {
				tx.Rollback()
				return nil, errors.New("terjadi kesalahan server")
			}
			exist.NamaToko = req.NamaToko
			exist.DeskripsiToko = req.DeskripsiToko
			exist.StatusVerifikasi = "pending"
			return &exist, nil
		}
		tx.Rollback()
		return nil, errors.New("user sudah mengajukan menjadi seller")
	}
	seller := models.Seller{
		UserID:        UserID,
		NamaToko:      req.NamaToko,
		DeskripsiToko: req.DeskripsiToko,
	}

	if err := tx.Create(&seller).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal membuat seller")
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New("terjadi kesalahan server")
	}
	return &seller, nil
}

// findSellerByUserID mengambil data toko milik user.
func findSellerByUserID(db *gorm.DB, userID string) (*models.Seller, error) {
	var seller models.Seller
	if err := db.Select("seller_id", "user_id", "seller_trust_score").
		Where("user_id = ?", userID).First(&seller).Error; err != nil {
		return nil, errors.New("data toko seller tidak ditemukan")
	}
	return &seller, nil
}

// GetAdminSellersService mengambil daftar toko seller (scope admin) dengan
// filter status verifikasi dan pagination.
func GetAdminSellersService(page, limit int, status string) ([]models.Seller, int64, error) {
	query := database.DB.Model(&models.Seller{})

	if s := normalizeSellerStatus(status); s != "" {
		query = query.Where("status_verifikasi = ?", s)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.New("gagal memuat data toko")
	}

	var sellers []models.Seller
	if err := query.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("user_id", "nama", "email")
	}).Order("created_at desc").
		Offset((page - 1) * limit).Limit(limit).Find(&sellers).Error; err != nil {
		return nil, 0, errors.New("gagal memuat data toko")
	}

	return sellers, total, nil
}

// VerifySellerService memperbarui status verifikasi pengajuan toko seller
// (scope admin). Jika disetujui (verified), peran user otomatis menjadi
// seller. Jika ditolak (rejected), user tetap customer dan boleh mengajukan
// ulang.
func VerifySellerService(sellerID, status string) (*models.Seller, error) {
	s := normalizeSellerStatus(status)
	if s != "verified" && s != "rejected" {
		return nil, errors.New("status tidak valid, gunakan VERIFIED atau REJECTED")
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		return nil, errors.New("terjadi kesalahan server")
	}

	var seller models.Seller
	if err := tx.Select("seller_id", "user_id", "status_verifikasi").
		Where("seller_id = ?", sellerID).First(&seller).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("toko seller tidak ditemukan")
	}
	if seller.StatusVerifikasi != "pending" {
		tx.Rollback()
		return nil, errors.New("hanya pengajuan berstatus pending yang dapat diverifikasi")
	}

	if err := tx.Model(&models.Seller{}).Where("seller_id = ?", sellerID).
		Update("status_verifikasi", s).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal memperbarui status toko")
	}

	isi := ""
	if s == "verified" {
		if err := tx.Model(&models.User{}).Where("user_id = ?", seller.UserID).
			Update("peran", "seller").Error; err != nil {
			tx.Rollback()
			return nil, errors.New("gagal memperbarui peran user")
		}
		isi = "Pengajuan toko Anda telah disetujui. Selamat, Anda sekarang menjadi seller!"
	} else {
		isi = "Pengajuan toko Anda ditolak. Anda dapat mengajukan ulang."
	}
	if err := createNotification(tx, seller.UserID, "dll", isi); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New("perubahan gagal disimpan")
	}

	seller.StatusVerifikasi = s
	return &seller, nil
}

// GetSellerProductsService mengambil daftar produk milik seller
// beserta total terjual, dengan filter status dan pagination.
func GetSellerProductsService(userID string, page, limit int, status string) ([]ProductWithSales, int64, error) {
	seller, err := findSellerByUserID(database.DB, userID)
	if err != nil {
		return nil, 0, err
	}

	salesSub := database.DB.Model(&models.OrderItem{}).
		Select("product_id, SUM(jumlah) AS total_terjual").
		Joins("JOIN orders ON orders.order_id = order_items.order_id").
		Where("orders.status_order <> ?", "dibatalkan").
		Group("product_id")

	query := database.DB.Model(&models.Product{}).
		Select("products.*, COALESCE(s.total_terjual, 0) AS total_terjual").
		Joins("LEFT JOIN (?) s ON s.product_id = products.product_id", salesSub).
		Where("products.seller_id = ?", seller.SellerID)

	if s := normalizeProductStatus(status); s != "" {
		query = query.Where("products.status_publikasi = ?", s)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.New("gagal menghitung produk")
	}

	var items []ProductWithSales
	if err := query.Order("products.created_at desc").
		Offset((page - 1) * limit).Limit(limit).Scan(&items).Error; err != nil {
		return nil, 0, errors.New("gagal memuat produk")
	}

	return items, total, nil
}

// GetSellerOrdersService mengambil daftar order milik toko seller
// beserta info customer, dengan filter status dan pagination.
func GetSellerOrdersService(userID string, page, limit int, status string) ([]models.Order, int64, error) {
	seller, err := findSellerByUserID(database.DB, userID)
	if err != nil {
		return nil, 0, err
	}

	query := database.DB.Model(&models.Order{}).
		Preload("Customer", func(db *gorm.DB) *gorm.DB {
			return db.Select("user_id", "nama")
		}).
		Where("seller_id = ?", seller.SellerID)

	if s := normalizeOrderStatus(status); s != "" {
		query = query.Where("status_order = ?", s)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.New("gagal memuat order")
	}

	var orders []models.Order
	if err := query.Order("created_at desc").
		Offset((page - 1) * limit).Limit(limit).Find(&orders).Error; err != nil {
		return nil, 0, errors.New("gagal memuat order")
	}

	return orders, total, nil
}

// GetSellerDashboardService menghitung statistik dashboard toko seller.
func GetSellerDashboardService(userID string) (*SellerDashboardData, error) {
	seller, err := findSellerByUserID(database.DB, userID)
	if err != nil {
		return nil, err
	}

	var totalProduk, produkAktif int64
	if err := database.DB.Model(&models.Product{}).
		Where("seller_id = ?", seller.SellerID).Count(&totalProduk).Error; err != nil {
		return nil, errors.New("gagal menghitung produk")
	}
	if err := database.DB.Model(&models.Product{}).
		Where("seller_id = ? AND status_publikasi = ?", seller.SellerID, "aktif").
		Count(&produkAktif).Error; err != nil {
		return nil, errors.New("gagal menghitung produk aktif")
	}

	var totalTerjual int64
	if err := database.DB.Model(&models.OrderItem{}).
		Select("COALESCE(SUM(order_items.jumlah), 0)").
		Joins("JOIN orders ON orders.order_id = order_items.order_id").
		Where("orders.seller_id = ? AND orders.status_order <> ?", seller.SellerID, "dibatalkan").
		Scan(&totalTerjual).Error; err != nil {
		return nil, errors.New("gagal menghitung total terjual")
	}

	var totalPendapatan float64
	if err := database.DB.Model(&models.Order{}).
		Select("COALESCE(SUM(total_pesanan), 0)").
		Where("seller_id = ? AND status_order = ?", seller.SellerID, "selesai").
		Scan(&totalPendapatan).Error; err != nil {
		return nil, errors.New("gagal menghitung total pendapatan")
	}

	var ratingRataRata float64
	if err := database.DB.Model(&models.Review{}).
		Select("COALESCE(AVG(reviews.rating), 0)").
		Joins("JOIN products ON products.product_id = reviews.product_id").
		Where("products.seller_id = ?", seller.SellerID).
		Scan(&ratingRataRata).Error; err != nil {
		return nil, errors.New("gagal menghitung rating")
	}

	var topProducts []SellerTopProduct
	if err := database.DB.Model(&models.OrderItem{}).
		Select("products.product_id, products.nama_produk, SUM(order_items.jumlah) AS total_terjual").
		Joins("JOIN products ON products.product_id = order_items.product_id").
		Joins("JOIN orders ON orders.order_id = order_items.order_id").
		Where("products.seller_id = ? AND orders.status_order <> ?", seller.SellerID, "dibatalkan").
		Group("products.product_id, products.nama_produk").
		Order("total_terjual desc").
		Limit(5).
		Scan(&topProducts).Error; err != nil {
		return nil, errors.New("gagal memuat produk terlaris")
	}

	return &SellerDashboardData{
		TotalProduk:      totalProduk,
		ProdukAktif:      produkAktif,
		TotalTerjual:     totalTerjual,
		TotalPendapatan:  totalPendapatan,
		RatingRataRata:   ratingRataRata,
		SellerTrustScore: seller.SellerTrustScore,
		ProdukTerlaris:   topProducts,
	}, nil
}
