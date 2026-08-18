package services

import (
	"errors"
	"strings"

	"github.com/ehanz12/api-SneakHub/database"
	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/requests"
	"gorm.io/gorm"
)

// orderScopeQuery membatasi query order berdasarkan peran pemanggil:
// customer hanya melihat order miliknya, seller hanya order tokonya,
// admin melihat semua.
func orderScopeQuery(db *gorm.DB, userID, role string) (*gorm.DB, error) {
	switch role {
	case "customer":
		return db.Where("customer_id = ?", userID), nil
	case "seller":
		var seller models.Seller
		if err := database.DB.Select("seller_id").Where("user_id = ?", userID).First(&seller).Error; err != nil {
			return nil, errors.New("data toko seller tidak ditemukan")
		}
		return db.Where("seller_id = ?", seller.SellerID), nil
	default:
		return db, nil
	}
}

// restoreOrderStock mengembalikan stok seluruh produk pada sebuah pesanan.
// Dipakai saat pesanan dibatalkan atau pembayaran expired/failed.
func restoreOrderStock(tx *gorm.DB, orderID string) error {
	var items []models.OrderItem
	if err := tx.Select("product_id", "jumlah").Where("order_id = ?", orderID).Find(&items).Error; err != nil {
		return errors.New("gagal memuat item pesanan")
	}
	for _, item := range items {
		if err := tx.Model(&models.Product{}).
			Where("product_id = ?", item.ProductID).
			Update("stok", gorm.Expr("stok + ?", item.Jumlah)).Error; err != nil {
			return errors.New("gagal mengembalikan stok produk")
		}
	}
	return nil
}

// normalizeOrderStatus memetakan alias status (mis. COMPLETED) ke nilai
// enum status_order di database. Mengembalikan "" bila tidak dikenal.
func normalizeOrderStatus(status string) string {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case "PENDING":
		return "pending"
	case "PROCESSING", "DIPROSES":
		return "diproses"
	case "SHIPPED", "DIKIRIM":
		return "dikirim"
	case "COMPLETED", "SELESAI":
		return "selesai"
	case "CANCELLED", "CANCELED", "DIBATALKAN":
		return "dibatalkan"
	}
	return ""
}

func GetOrdersService(userID, role string, page, limit int, status string) ([]models.Order, int64, error) {
	query := database.DB.Model(&models.Order{})
	query, err := orderScopeQuery(query, userID, role)
	if err != nil {
		return nil, 0, err
	}
	if s := normalizeOrderStatus(status); s != "" {
		query = query.Where("status_order = ?", s)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.New("gagal memuat order")
	}

	var orders []models.Order
	if err := query.Preload("Payment").Order("created_at desc").Offset((page - 1) * limit).Limit(limit).Find(&orders).Error; err != nil {
		return nil, 0, errors.New("gagal memuat order")
	}

	return orders, total, nil
}

func GetOrderService(userID, role, orderID string) (*models.Order, error) {
	query := database.DB.
		Preload("Items.Product", func(db *gorm.DB) *gorm.DB {
			return db.Select("product_id", "nama_produk")
		}).
		Preload("Payment").
		Preload("Shipment").
		Where("order_id = ?", orderID)
	query, err := orderScopeQuery(query, userID, role)
	if err != nil {
		return nil, err
	}

	var order models.Order
	if err := query.First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order tidak ditemukan")
		}
		return nil, errors.New("gagal memuat order")
	}

	return &order, nil
}

func CreateOrderService(userID, role string, r requests.CreateOrderRequest) (*models.Order, error) {
	tx := database.DB.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return nil, errors.New("transaksi database gagal dimulai")
	}

	customerID := r.CustomerID
	sellerID := r.SellerID
	switch role {
	case "customer":
		customerID = userID
	case "seller":
		var seller models.Seller
		if err := tx.Select("seller_id").Where("user_id = ?", userID).First(&seller).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("data toko seller tidak ditemukan")
		}
		sellerID = seller.SellerID
	default:
		if strings.TrimSpace(customerID) == "" {
			tx.Rollback()
			return nil, errors.New("customer_id wajib diisi")
		}
		if strings.TrimSpace(sellerID) == "" {
			tx.Rollback()
			return nil, errors.New("seller_id wajib diisi")
		}
	}

	var sellerCheck int64
	if err := tx.Model(&models.Seller{}).Where("seller_id = ?", sellerID).Count(&sellerCheck).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal memvalidasi toko")
	}
	if sellerCheck == 0 {
		tx.Rollback()
		return nil, errors.New("seller tidak ditemukan")
	}

	var customerCheck int64
	if err := tx.Model(&models.User{}).Where("user_id = ?", customerID).Count(&customerCheck).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal memvalidasi pelanggan")
	}
	if customerCheck == 0 {
		tx.Rollback()
		return nil, errors.New("customer tidak ditemukan")
	}

	subtotal := 0.0
	for _, item := range r.Items {
		var product models.Product
		if err := tx.Select("product_id", "nama_produk", "harga", "stok", "status_publikasi").
			Where("product_id = ?", item.ProductID).First(&product).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("produk tidak ditemukan: " + item.ProductID)
		}
		if product.StatusPublikasi != "aktif" {
			tx.Rollback()
			return nil, errors.New("produk tidak aktif: " + item.ProductID)
		}
		if product.Stok < item.Jumlah {
			tx.Rollback()
			return nil, errors.New("stok produk tidak mencukupi: " + item.ProductID)
		}
		res := tx.Model(&models.Product{}).
			Where("product_id = ? AND stok >= ?", item.ProductID, item.Jumlah).
			Update("stok", gorm.Expr("stok - ?", item.Jumlah))
		if res.Error != nil {
			tx.Rollback()
			return nil, errors.New("gagal memperbarui stok produk")
		}
		if res.RowsAffected == 0 {
			tx.Rollback()
			return nil, errors.New("stok produk sudah habis: " + item.ProductID)
		}
		subtotal += product.Harga * float64(item.Jumlah)
	}

	snapshot, err := buildAlamatSnapshot(
		r.AlamatPengiriman.NamaPenerima,
		r.AlamatPengiriman.NomorTelepon,
		r.AlamatPengiriman.Alamat,
		r.AlamatPengiriman.Kota,
		r.AlamatPengiriman.Provinsi,
		r.AlamatPengiriman.KodePos,
	)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("gagal menyimpan data alamat")
	}

	order := models.Order{
		CustomerID:       customerID,
		SellerID:         sellerID,
		StatusOrder:      "pending",
		AlamatPengiriman: snapshot,
		MetodePembayaran: r.MetodePembayaran,
		Subtotal:         subtotal,
		BiayaPengiriman:  BiayaPengirimanFlat,
		TotalPesanan:     subtotal + BiayaPengirimanFlat,
	}
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal membuat pesanan")
	}

	for _, item := range r.Items {
		var product models.Product
		if err := tx.Select("harga").Where("product_id = ?", item.ProductID).First(&product).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("produk tidak ditemukan: " + item.ProductID)
		}
		if err := tx.Create(&models.OrderItem{
			OrderID:            order.OrderID,
			ProductID:          item.ProductID,
			Jumlah:             item.Jumlah,
			HargaSaatTransaksi: product.Harga,
		}).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("gagal menyimpan item pesanan")
		}
	}

	if err := tx.Create(&models.Shipment{
		OrderID: order.OrderID,
		Kurir:   "flat",
	}).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal menyimpan data pengiriman")
	}

	var customer models.User
	if err := tx.Select("nama", "email", "nomor_telepon").Where("user_id = ?", customerID).First(&customer).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal memuat data pelanggan")
	}
	if err := createOrderPayment(tx, &order, customer); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New("perubahan order gagal disimpan")
	}

	return &order, nil
}

func UpdateOrderStatusService(userID, role, orderID, status string) (*models.Order, error) {
	tx := database.DB.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return nil, errors.New("transaksi database gagal dimulai")
	}

	query, err := orderScopeQuery(tx.Model(&models.Order{}), userID, role)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var order models.Order
	if err := query.Where("order_id = ?", orderID).First(&order).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order tidak ditemukan")
		}
		return nil, errors.New("gagal memuat order")
	}

	normalized := normalizeOrderStatus(status)
	if normalized == "" {
		tx.Rollback()
		return nil, errors.New("status_order tidak valid")
	}
	if normalized == "pending" {
		tx.Rollback()
		return nil, errors.New("status pending hanya berlaku saat pesanan dibuat")
	}

	if role == "customer" && normalized != "dibatalkan" {
		tx.Rollback()
		return nil, errors.New("customer hanya dapat membatalkan pesanan")
	}
	if role == "seller" {
		switch normalized {
		case "diproses":
			if order.StatusOrder != "pending" {
				tx.Rollback()
				return nil, errors.New("pesanan hanya dapat diproses dari status pending")
			}
		case "dibatalkan":
			if order.StatusOrder != "pending" && order.StatusOrder != "diproses" {
				tx.Rollback()
				return nil, errors.New("pesanan hanya dapat dibatalkan sebelum dikirim")
			}
		case "dikirim", "selesai":
			tx.Rollback()
			return nil, errors.New("status dikirim/selesai diatur lewat endpoint pengiriman atau konfirmasi")
		}
	}
	if normalized == "dibatalkan" {
		if order.StatusOrder != "pending" && order.StatusOrder != "diproses" {
			tx.Rollback()
			return nil, errors.New("pesanan hanya dapat dibatalkan sebelum dikirim")
		}
		if err := restoreOrderStock(tx, orderID); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Model(&order).Update("status_order", normalized).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal memperbarui status order")
	}
	order.StatusOrder = normalized

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New("perubahan order gagal disimpan")
	}

	return &order, nil
}

func DeleteOrderService(userID, role, orderID string) error {
	if role == "customer" {
		return errors.New("customer tidak memiliki akses menghapus order")
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return errors.New("transaksi database gagal dimulai")
	}

	query, err := orderScopeQuery(tx.Model(&models.Order{}), userID, role)
	if err != nil {
		tx.Rollback()
		return err
	}

	var order models.Order
	if err := query.Where("order_id = ?", orderID).First(&order).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("order tidak ditemukan")
		}
		return errors.New("gagal memuat order")
	}

	if err := tx.Where("order_id = ?", orderID).Delete(&models.Payment{}).Error; err != nil {
		tx.Rollback()
		return errors.New("gagal menghapus data pembayaran")
	}
	if err := tx.Where("order_id = ?", orderID).Delete(&models.OrderItem{}).Error; err != nil {
		tx.Rollback()
		return errors.New("gagal menghapus item pesanan")
	}
	if err := tx.Delete(&order).Error; err != nil {
		tx.Rollback()
		return errors.New("gagal menghapus order")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errors.New("perubahan order gagal disimpan")
	}

	return nil
}
