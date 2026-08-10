package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Brand struct {
	BrandID   string `gorm:"column:brand_id;type:char(36);primaryKey" json:"brand_id"`
	NamaBrand string `gorm:"column:nama_brand;type:varchar(100);unique;not null" json:"nama_brand"`
	Timestamps

	Products        []Product         `gorm:"foreignKey:BrandID;references:BrandID" json:"products,omitempty"`
	MarketPriceData []MarketPriceData `gorm:"foreignKey:BrandID;references:BrandID" json:"market_price_data,omitempty"`
}

func (b *Brand) BeforeCreate(tx *gorm.DB) error {
	if b.BrandID == "" {
		b.BrandID = uuid.NewString()
	}
	return nil
}

func (Brand) TableName() string { return "brands" }
