package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MarketPriceData struct {
	MarketPriceID string    `gorm:"column:market_price_id;type:char(36);primaryKey" json:"market_price_id"`
	BrandID       string    `gorm:"column:brand_id;type:char(36);not null" json:"brand_id"`
	CategoryID    *string   `gorm:"column:category_id;type:char(36)" json:"category_id,omitempty"`
	NamaModel     string    `gorm:"column:nama_model;type:varchar(200);not null" json:"nama_model"`
	Kondisi       string    `gorm:"column:kondisi;type:enum('new','used','refurbished');not null" json:"kondisi"`
	Ukuran        *string   `gorm:"column:ukuran;type:varchar(30)" json:"ukuran,omitempty"`
	Harga         float64   `gorm:"column:harga;type:decimal(15,2);not null" json:"harga"`
	Sumber        *string   `gorm:"column:sumber;type:varchar(255)" json:"sumber,omitempty"`
	RecordedAt    time.Time `gorm:"column:recorded_at" json:"recorded_at"`

	Brand    Brand     `gorm:"foreignKey:BrandID;references:BrandID" json:"brand,omitempty"`
	Category *Category `gorm:"foreignKey:CategoryID;references:CategoryID" json:"category,omitempty"`
}

func (m *MarketPriceData) BeforeCreate(tx *gorm.DB) error {
	if m.MarketPriceID == "" {
		m.MarketPriceID = uuid.NewString()
	}
	return nil
}
func (MarketPriceData) TableName() string { return "market_price_data" }
