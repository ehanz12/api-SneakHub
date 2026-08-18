package services

import (
	"errors"
	"strings"
	"time"

	"github.com/ehanz12/api-SneakHub/database"
	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/requests"
	"gorm.io/gorm"
)

// ShipOrderService memproses pengiriman pesanan oleh seller. Nomor resi
// didapat otomatis dari booking Biteship; jika booking gagal atau request
// berisi nomor_resi, dipakai resi manual sebagai fallback.
func ShipOrderService(userID, orderID string, r requests.ShipOrderRequest) (*models.Order, error) {
	tx := database.DB.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return nil, errors.New("transaksi database gagal dimulai")
	}

	var store models.Seller
	if err := tx.Select("seller_id", "user_id", "nama_toko", "alamat_asal", "kode_pos_asal", "kota_asal").
		Where("user_id = ?", userID).First(&store).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("data toko seller tidak ditemukan")
	}

	var order models.Order
	if err := tx.
		Preload("Shipment").
		Preload("Items.Product", func(db *gorm.DB) *gorm.DB {
			return db.Select("product_id", "nama_produk", "berat")
		}).
		Where("order_id = ? AND seller_id = ?", orderID, store.SellerID).
		First(&order).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order tidak ditemukan")
		}
		return nil, errors.New("gagal memuat order")
	}

	if order.StatusOrder != "diproses" {
		tx.Rollback()
		return nil, errors.New("pesanan hanya dapat dikirim dari status diproses")
	}
	if order.Shipment == nil {
		tx.Rollback()
		return nil, errors.New("data pengiriman tidak ditemukan")
	}
	if order.Shipment.NomorResi != nil && strings.TrimSpace(*order.Shipment.NomorResi) != "" {
		tx.Rollback()
		return nil, errors.New("pesanan sudah dikirim")
	}

	nomorResi := ""
	if r.NomorResi != nil {
		nomorResi = strings.TrimSpace(*r.NomorResi)
	}
	trackingID := ""

	if nomorResi == "" {
		var sellerUser models.User
		if err := tx.Select("nama", "email", "nomor_telepon").Where("user_id = ?", store.UserID).First(&sellerUser).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("gagal memuat data seller")
		}

		resi, waybill, err := CreateBiteshipOrder(&order, order.Shipment, order.Items, store, sellerUser)
		if err != nil {
			tx.Rollback()
			return nil, errors.New(err.Error() + "; gunakan nomor_resi manual sebagai pengganti")
		}
		nomorResi = resi
		trackingID = waybill
	}

	now := time.Now()
	updates := map[string]interface{}{
		"nomor_resi":        nomorResi,
		"status_pengiriman": "dikirim",
		"shipped_at":        now,
	}
	if trackingID != "" {
		updates["tracking_id"] = trackingID
	}
	if err := tx.Model(&models.Shipment{}).
		Where("shipment_id = ?", order.Shipment.ShipmentID).
		Updates(updates).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal memperbarui data pengiriman")
	}

	if err := tx.Model(&models.Order{}).
		Where("order_id = ?", orderID).
		Update("status_order", "dikirim").Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal memperbarui status pesanan")
	}

	if err := createNotification(tx, order.CustomerID, "order_update",
		"Pesanan "+orderID+" sudah dikirim. Nomor resi: "+nomorResi); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New("perubahan pengiriman gagal disimpan")
	}

	order.StatusOrder = "dikirim"
	order.Shipment.NomorResi = &nomorResi
	order.Shipment.StatusPengiriman = "dikirim"
	order.Shipment.ShippedAt = &now
	return &order, nil
}

// ConfirmReceivedService menandai pesanan selesai oleh customer
// (konfirmasi barang sudah diterima).
func ConfirmReceivedService(userID, orderID string) (*models.Order, error) {
	tx := database.DB.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return nil, errors.New("transaksi database gagal dimulai")
	}

	var order models.Order
	if err := tx.Preload("Shipment").
		Where("order_id = ? AND customer_id = ?", orderID, userID).
		First(&order).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order tidak ditemukan")
		}
		return nil, errors.New("gagal memuat order")
	}

	if order.StatusOrder != "dikirim" {
		tx.Rollback()
		return nil, errors.New("pesanan hanya dapat dikonfirmasi setelah dikirim")
	}

	if err := tx.Model(&models.Order{}).
		Where("order_id = ?", orderID).
		Update("status_order", "selesai").Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal memperbarui status pesanan")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New("perubahan pesanan gagal disimpan")
	}

	order.StatusOrder = "selesai"
	return &order, nil
}
