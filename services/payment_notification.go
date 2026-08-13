package services

import (
	"errors"
	"time"

	"github.com/ehanz12/api-SneakHub/database"
	"github.com/ehanz12/api-SneakHub/models"
)

// HandlePaymentNotificationService memperbarui status pembayaran & pesanan
// berdasarkan callback dari Tripay.
func HandlePaymentNotificationService(orderID, transactionStatus, transactionID string, grossAmount float64) error {
	tx := database.DB.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return errors.New("koneksi database gagal")
	}

	var payment models.Payment
	if err := tx.Where("order_id = ?", orderID).First(&payment).Error; err != nil {
		tx.Rollback()
		return errors.New("data pembayaran tidak ditemukan")
	}

	if grossAmount > 0 && grossAmount != payment.Jumlah {
		tx.Rollback()
		return errors.New("jumlah pembayaran tidak sesuai")
	}

	updates := make(map[string]interface{})

	switch transactionStatus {
	case "PAID":
		updates["status_pembayaran"] = "paid"
		updates["paid_at"] = time.Now()
		if err := tx.Model(&models.Order{}).Where("order_id = ?", orderID).Update("status_order", "diproses").Error; err != nil {
			tx.Rollback()
			return errors.New("gagal memperbarui status pesanan")
		}
	case "EXPIRED":
		updates["status_pembayaran"] = "expired"
	case "FAILED":
		updates["status_pembayaran"] = "failed"
	case "REFUND":
		updates["status_pembayaran"] = "refunded"
	default:
		tx.Rollback()
		return nil
	}

	if transactionID != "" {
		updates["transaction_reference"] = transactionID
	}

	if err := tx.Model(&payment).Updates(updates).Error; err != nil {
		tx.Rollback()
		return errors.New("gagal memperbarui status pembayaran")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errors.New("status pembayaran gagal disimpan")
	}

	return nil
}
