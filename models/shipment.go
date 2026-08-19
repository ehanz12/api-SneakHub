package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ShipmentTrackingEvent struct {
	Note      string `json:"note"`
	Status    string `json:"status"`
	UpdatedAt string `json:"updated_at"`
}

type Shipment struct {
	ShipmentID       string         `gorm:"column:shipment_id;type:char(36);primaryKey" json:"shipment_id"`
	OrderID          string         `gorm:"column:order_id;type:char(36);not null;unique" json:"order_id"`
	Kurir            string         `gorm:"column:kurir;type:varchar(100);not null" json:"kurir"`
	Service          *string        `gorm:"column:service;type:varchar(100)" json:"service,omitempty"`
	NomorResi        *string        `gorm:"column:nomor_resi;type:varchar(100)" json:"nomor_resi,omitempty"`
	TrackingID       *string        `gorm:"column:tracking_id;type:varchar(100)" json:"tracking_id,omitempty"`
	TrackingHistory  datatypes.JSON `gorm:"column:tracking_history;type:json" json:"tracking_history,omitempty"`
	StatusPengiriman string         `gorm:"column:status_pengiriman;type:enum('menunggu','dikirim','dalam_perjalanan','sampai');default:menunggu" json:"status_pengiriman"`
	ShippedAt        *time.Time     `gorm:"column:shipped_at" json:"shipped_at,omitempty"`
	DeliveredAt      *time.Time     `gorm:"column:delivered_at" json:"delivered_at,omitempty"`
	CreatedAt        time.Time      `gorm:"column:created_at" json:"created_at"`

	Order Order `gorm:"foreignKey:OrderID;references:OrderID" json:"order,omitempty"`
}

func (s *Shipment) BeforeCreate(tx *gorm.DB) error {
	if s.ShipmentID == "" {
		s.ShipmentID = uuid.NewString()
	}
	return nil
}
func (Shipment) TableName() string { return "shipments" }
