package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Notification struct {
	NotificationID  string    `gorm:"column:notification_id;type:char(36);primaryKey" json:"notification_id"`
	UserID          string    `gorm:"column:user_id;type:char(36);not null" json:"user_id"`
	JenisNotifikasi string    `gorm:"column:jenis_notifikasi;type:enum('price_alert','restock_alert','order_update','promo','dll');not null" json:"jenis_notifikasi"`
	IsiNotifikasi   string    `gorm:"column:isi_notifikasi;type:text;not null" json:"isi_notifikasi"`
	StatusBaca      bool      `gorm:"column:status_baca;not null;default:false" json:"status_baca"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`

	User User `gorm:"foreignKey:UserID;references:UserID" json:"user,omitempty"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.NotificationID == "" {
		n.NotificationID = uuid.NewString()
	}
	return nil
}
func (Notification) TableName() string { return "notifications" }
