package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductVariant struct {
	VariantID string    `gorm:"column:variant_id;type:char(36);primaryKey" json:"variant_id"`
	ProductID string    `gorm:"column:product_id;type:char(36);not null" json:"product_id"`
	Ukuran    string    `gorm:"column:ukuran;type:varchar(30);not null" json:"ukuran"`
	Stok      int       `gorm:"column:stok;not null;default:0" json:"stok"`
	Harga     float64   `gorm:"column:harga;type:decimal(15,2);not null" json:"harga"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

	Product Product `gorm:"foreignKey:ProductID;references:ProductID" json:"product,omitempty"`
}

func (p *ProductVariant) BeforeCreate(tx *gorm.DB) error {
	if p.VariantID == "" {
		p.VariantID = uuid.NewString()
	}
	return nil
}
func (ProductVariant) TableName() string { return "product_variants" }
