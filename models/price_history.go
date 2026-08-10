package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PriceHistory struct {
	PriceHistoryID string    `gorm:"column:price_history_id;type:char(36);primaryKey" json:"price_history_id"`
	ProductID      string    `gorm:"column:product_id;type:char(36);not null" json:"product_id"`
	HargaLama      float64   `gorm:"column:harga_lama;type:decimal(15,2);not null" json:"harga_lama"`
	HargaBaru      float64   `gorm:"column:harga_baru;type:decimal(15,2);not null" json:"harga_baru"`
	WaktuPerubahan time.Time `gorm:"column:waktu_perubahan" json:"waktu_perubahan"`

	Product Product `gorm:"foreignKey:ProductID;references:ProductID" json:"product,omitempty"`
}

func (p *PriceHistory) BeforeCreate(tx *gorm.DB) error {
	if p.PriceHistoryID == "" {
		p.PriceHistoryID = uuid.NewString()
	}
	return nil
}
func (PriceHistory) TableName() string { return "price_history" }
