package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CartItem struct {
	CartItemID           string    `gorm:"column:cart_item_id;type:char(36);primaryKey" json:"cart_item_id"`
	CartID               string    `gorm:"column:cart_id;type:char(36);not null" json:"cart_id"`
	ProductID            string    `gorm:"column:product_id;type:char(36);not null" json:"product_id"`
	Jumlah               int       `gorm:"column:jumlah;not null;default:1" json:"jumlah"`
	HargaSaatDitambahkan float64   `gorm:"column:harga_saat_ditambahkan;type:decimal(15,2);not null" json:"harga_saat_ditambahkan"`
	CreatedAt            time.Time `gorm:"column:created_at" json:"created_at"`

	Cart    Cart    `gorm:"foreignKey:CartID;references:CartID" json:"cart,omitempty"`
	Product Product `gorm:"foreignKey:ProductID;references:ProductID" json:"product,omitempty"`
}

func (c *CartItem) BeforeCreate(tx *gorm.DB) error {
	if c.CartItemID == "" {
		c.CartItemID = uuid.NewString()
	}
	return nil
}
func (CartItem) TableName() string { return "cart_items" }
