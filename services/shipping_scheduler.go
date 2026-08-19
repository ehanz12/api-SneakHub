package services

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/ehanz12/api-SneakHub/config"
	"github.com/ehanz12/api-SneakHub/database"
	"github.com/ehanz12/api-SneakHub/models"
	"gorm.io/datatypes"
)

func StartShippingScheduler() {
	interval := config.AppConfig.ShippingPollMinutes
	if interval <= 0 {
		interval = 30
	}
	go func() {
		ticker := time.NewTicker(time.Duration(interval) * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			pollShipmentTracking()
			autoCompleteOrders()
		}
	}()
	log.Println("👀 SHIPPING SCHEDULER STARTED")
}

func pollShipmentTracking() {
	var shipments []models.Shipment
	if err := database.DB.
		Where("tracking_id IS NOT NULL AND tracking_id <> '' AND status_pengiriman <> ?", "sampai").
		Find(&shipments).Error; err != nil {
		log.Println("[shipping-scheduler] gagal memuat shipment:", err)
		return
	}
	if len(shipments) == 0 {
		return
	}

	for _, shipment := range shipments {
		if shipment.TrackingID == nil || strings.TrimSpace(*shipment.TrackingID) == "" {
			continue
		}

		status, history, err := TrackBiteshipShipment(*shipment.TrackingID)
		if err != nil {
			log.Printf("[shipping-scheduler] gagal tracking %s: %v", *shipment.TrackingID, err)
			continue
		}

		statusPengiriman := mapBiteshipCourierStatus(status)
		if statusPengiriman == shipment.StatusPengiriman {
			continue
		}

		updates := map[string]interface{}{
			"status_pengiriman": statusPengiriman,
		}
		if len(history) > 0 {
			raw, err := json.Marshal(history)
			if err == nil {
				updates["tracking_history"] = datatypes.JSON(raw)
			}
		}
		if statusPengiriman == "sampai" {
			updates["delivered_at"] = time.Now()
		}

		if err := database.DB.Model(&models.Shipment{}).
			Where("shipment_id = ?", shipment.ShipmentID).
			Updates(updates).Error; err != nil {
			log.Printf("[shipping-scheduler] gagal update shipment %s: %v", shipment.ShipmentID, err)
			continue
		}

		var order models.Order
		if err := database.DB.Select("customer_id", "order_id").
			Where("order_id = ?", shipment.OrderID).First(&order).Error; err != nil {
			continue
		}

		if statusPengiriman == "sampai" {
			_ = createNotification(database.DB, order.CustomerID, "order_update",
				"Paket pesanan "+order.OrderID+" sudah sampai. Jangan lupa konfirmasi pesanan selesai.")
		}
		log.Printf("[shipping-scheduler] shipment %s (%s) -> %s", shipment.ShipmentID, *shipment.TrackingID, statusPengiriman)
	}
}

func autoCompleteOrders() {
	days := config.AppConfig.AutoCompleteDays
	if days <= 0 {
		days = 7
	}
	deadline := time.Now().Add(-time.Duration(days) * 24 * time.Hour)

	var orderIDs []string
	if err := database.DB.Model(&models.Order{}).
		Joins("JOIN shipments ON shipments.order_id = orders.order_id").
		Where("orders.status_order = ?", "dikirim").
		Where("shipments.status_pengiriman = ?", "sampai").
		Where("shipments.delivered_at IS NOT NULL").
		Where("shipments.delivered_at < ?", deadline).
		Pluck("orders.order_id", &orderIDs).Error; err != nil {
		log.Println("[shipping-scheduler] gagal memuat order auto-complete:", err)
		return
	}
	if len(orderIDs) == 0 {
		return
	}

	res := database.DB.Model(&models.Order{}).
		Where("order_id IN ? AND status_order = ?", orderIDs, "dikirim").
		Update("status_order", "selesai")
	if res.Error != nil {
		log.Println("[shipping-scheduler] gagal auto-complete:", res.Error)
		return
	}
	if res.RowsAffected > 0 {
		log.Printf("[shipping-scheduler] %d pesanan otomatis selesai", res.RowsAffected)
	}
}
