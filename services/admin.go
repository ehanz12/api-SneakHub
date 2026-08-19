package services

import (
	"errors"
	"strings"
	"time"

	"github.com/ehanz12/api-SneakHub/database"
	"github.com/ehanz12/api-SneakHub/models"
	"gorm.io/gorm"
)

// AdminReportData adalah hasil agregasi laporan admin.
type AdminReportData struct {
	Period        string
	TotalUsers    int64
	TotalSellers  int64
	TotalProducts int64
	TotalOrders   int64
	TotalRevenue  float64
}

// GetAdminUsersService mengambil daftar pengguna dengan filter
// status akun, peran, dan pagination.
func GetAdminUsersService(page, limit int, status, role string) ([]models.User, int64, error) {
	query := database.DB.Model(&models.User{})

	if s := normalizeUserStatus(status); s != "" {
		query = query.Where("status_akun = ?", s)
	}
	if r := normalizeRole(role); r != "" {
		query = query.Where("peran = ?", r)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.New("gagal memuat data pengguna")
	}

	var users []models.User
	if err := query.Order("created_at desc").
		Offset((page - 1) * limit).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, errors.New("gagal memuat data pengguna")
	}

	return users, total, nil
}

// UpdateUserStatusService memperbarui status akun pengguna.
func UpdateUserStatusService(userID, status string) (*models.User, error) {
	s := normalizeUserStatus(status)
	if s == "" {
		return nil, errors.New("status_akun tidak valid")
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		return nil, errors.New("terjadi kesalahan server")
	}

	var user models.User
	if err := tx.Select("user_id").Where("user_id = ?", userID).First(&user).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("user tidak ditemukan")
	}

	if err := tx.Model(&models.User{}).Where("user_id = ?", userID).
		Update("status_akun", s).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal memperbarui status akun")
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New("perubahan gagal disimpan")
	}

	user.StatusAkun = s
	return &user, nil
}

// UpdateUserRoleService mengubah peran user (customer <-> seller) oleh admin.
// Saat seller diturunkan menjadi customer, toko ditandai rejected dan seluruh
// produknya dinonaktifkan agar tidak bisa dibeli lagi.
func UpdateUserRoleService(userID, role string) (*models.User, error) {
	r := normalizeRole(role)
	if r != "seller" && r != "customer" {
		return nil, errors.New("peran hanya bisa diubah ke SELLER atau CUSTOMER")
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		return nil, errors.New("terjadi kesalahan server")
	}

	var user models.User
	if err := tx.Select("user_id", "nama", "email", "peran").
		Where("user_id = ?", userID).First(&user).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("user tidak ditemukan")
	}
	if user.Peran == "admin" {
		tx.Rollback()
		return nil, errors.New("peran admin tidak dapat diubah")
	}
	if user.Peran == r {
		tx.Rollback()
		return nil, errors.New("peran user sudah " + r)
	}

	var seller models.Seller
	sellerFound := true
	if err := tx.Select("seller_id", "user_id", "status_verifikasi").
		Where("user_id = ?", userID).First(&seller).Error; err != nil {
		if r == "seller" {
			tx.Rollback()
			return nil, errors.New("user belum mengajukan menjadi seller")
		}
		sellerFound = false
	}

	isi := ""
	switch r {
	case "seller":
		// Naikkan customer menjadi seller: toko disetujui.
		if err := tx.Model(&models.Seller{}).Where("seller_id = ?", seller.SellerID).
			Update("status_verifikasi", "verified").Error; err != nil {
			tx.Rollback()
			return nil, errors.New("gagal memperbarui status toko")
		}
		isi = "Peran Anda telah diubah menjadi seller oleh admin."
	case "customer":
		// Turunkan seller menjadi customer: toko ditolak & produk dinonaktifkan.
		if sellerFound {
			if err := tx.Model(&models.Seller{}).Where("seller_id = ?", seller.SellerID).
				Update("status_verifikasi", "rejected").Error; err != nil {
				tx.Rollback()
				return nil, errors.New("gagal menonaktifkan toko")
			}
			if err := tx.Model(&models.Product{}).Where("seller_id = ?", seller.SellerID).
				Update("status_publikasi", "nonaktif").Error; err != nil {
				tx.Rollback()
				return nil, errors.New("gagal menonaktifkan produk")
			}
		}
		isi = "Peran seller Anda telah dinonaktifkan oleh admin. Toko dan produk Anda dinonaktifkan."
	}

	if err := tx.Model(&models.User{}).Where("user_id = ?", userID).
		Update("peran", r).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal memperbarui peran user")
	}
	if err := createNotification(tx, userID, "dll", isi); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New("perubahan gagal disimpan")
	}

	user.Peran = r
	return &user, nil
}

// GetAdminProductsService mengambil daftar semua produk dengan filter
// status publikasi dan pagination.
func GetAdminProductsService(page, limit int, status string) ([]models.Product, int64, error) {
	query := database.DB.Model(&models.Product{})

	if s := normalizeProductStatus(status); s != "" {
		query = query.Where("status_publikasi = ?", s)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.New("gagal memuat data produk")
	}

	var products []models.Product
	if err := query.Preload("Images", func(db *gorm.DB) *gorm.DB {
		return db.Order("urutan_tampil asc")
	}).Order("created_at desc").
		Offset((page - 1) * limit).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, errors.New("gagal memuat data produk")
	}

	return products, total, nil
}

// UpdateProductStatusService memperbarui status publikasi produk (moderasi).
func UpdateProductStatusService(productID, status string) (*models.Product, error) {
	s := normalizeProductStatus(status)
	if s == "" {
		return nil, errors.New("status_publikasi tidak valid")
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		return nil, errors.New("terjadi kesalahan server")
	}

	var product models.Product
	if err := tx.Select("product_id").Where("product_id = ?", productID).First(&product).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("produk tidak ditemukan")
	}

	if err := tx.Model(&models.Product{}).Where("product_id = ?", productID).
		Update("status_publikasi", s).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal memperbarui status produk")
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New("perubahan gagal disimpan")
	}

	product.StatusPublikasi = s
	return &product, nil
}

// GetAdminOrdersService mengambil daftar semua order (scope admin)
// dengan filter status dan pagination.
func GetAdminOrdersService(page, limit int, status string) ([]models.Order, int64, error) {
	return GetOrdersService("", "admin", page, limit, status)
}

// GetAdminReportsService menghitung laporan agregat platform dalam
// rentang tanggal (opsional).
func GetAdminReportsService(period, startDate, endDate string) (*AdminReportData, error) {
	period = strings.ToUpper(strings.TrimSpace(period))
	if period == "" {
		period = "MONTHLY"
	}
	switch period {
	case "DAILY", "WEEKLY", "MONTHLY", "YEARLY":
	default:
		return nil, errors.New("period harus salah satu dari: DAILY, WEEKLY, MONTHLY, YEARLY")
	}

	start, err := parseReportDate(startDate)
	if err != nil {
		return nil, err
	}
	end, err := parseReportDate(endDate)
	if err != nil {
		return nil, err
	}

	userQuery := database.DB.Model(&models.User{})
	sellerQuery := database.DB.Model(&models.Seller{})
	productQuery := database.DB.Model(&models.Product{})
	orderQuery := database.DB.Model(&models.Order{})

	if start != nil {
		userQuery = userQuery.Where("created_at >= ?", *start)
		sellerQuery = sellerQuery.Where("created_at >= ?", *start)
		productQuery = productQuery.Where("created_at >= ?", *start)
		orderQuery = orderQuery.Where("created_at >= ?", *start)
	}
	if end != nil {
		userQuery = userQuery.Where("created_at <= ?", *end)
		sellerQuery = sellerQuery.Where("created_at <= ?", *end)
		productQuery = productQuery.Where("created_at <= ?", *end)
		orderQuery = orderQuery.Where("created_at <= ?", *end)
	}

	report := &AdminReportData{Period: period}

	if err := userQuery.Count(&report.TotalUsers).Error; err != nil {
		return nil, errors.New("gagal menghitung pengguna")
	}
	if err := sellerQuery.Count(&report.TotalSellers).Error; err != nil {
		return nil, errors.New("gagal menghitung seller")
	}
	if err := productQuery.Count(&report.TotalProducts).Error; err != nil {
		return nil, errors.New("gagal menghitung produk")
	}
	if err := orderQuery.Count(&report.TotalOrders).Error; err != nil {
		return nil, errors.New("gagal menghitung order")
	}

	revenueQuery := database.DB.Model(&models.Order{})
	if start != nil {
		revenueQuery = revenueQuery.Where("created_at >= ?", *start)
	}
	if end != nil {
		revenueQuery = revenueQuery.Where("created_at <= ?", *end)
	}
	if err := revenueQuery.Where("status_order = ?", "selesai").
		Select("COALESCE(SUM(total_pesanan), 0)").
		Scan(&report.TotalRevenue).Error; err != nil {
		return nil, errors.New("gagal menghitung pendapatan")
	}

	return report, nil
}

// parseReportDate mem-parse tanggal YYYY-MM-DD. Mengembalikan nil bila kosong.
func parseReportDate(date string) (*time.Time, error) {
	date = strings.TrimSpace(date)
	if date == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, errors.New("format tanggal harus YYYY-MM-DD")
	}
	return &t, nil
}
