package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	CategoryID   string `gorm:"column:category_id;type:char(36);primaryKey" json:"category_id"`
	NamaKategori string `gorm:"column:nama_kategori;type:varchar(100);unique;not null" json:"nama_kategori"`
	Timestamps

	Products        []Product         `gorm:"foreignKey:CategoryID;references:CategoryID" json:"products,omitempty"`
	MarketPriceData []MarketPriceData `gorm:"foreignKey:CategoryID;references:CategoryID" json:"market_price_data,omitempty"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.CategoryID == "" {
		c.CategoryID = uuid.NewString()
	}
	return nil
}
func (Category) TableName() string { return "categories" }
