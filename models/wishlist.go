package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Wishlist struct {
	WishlistID string    `gorm:"column:wishlist_id;type:char(36);primaryKey" json:"wishlist_id"`
	CustomerID string    `gorm:"column:customer_id;type:char(36);not null" json:"customer_id"`
	ProductID  string    `gorm:"column:product_id;type:char(36);not null" json:"product_id"`
	StatusStok string    `gorm:"column:status_stok;type:enum('available','out_of_stock');default:available" json:"status_stok"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`

	Customer User    `gorm:"foreignKey:CustomerID;references:UserID" json:"customer,omitempty"`
	Product  Product `gorm:"foreignKey:ProductID;references:ProductID" json:"product,omitempty"`
}

func (w *Wishlist) BeforeCreate(tx *gorm.DB) error {
	if w.WishlistID == "" {
		w.WishlistID = uuid.NewString()
	}
	return nil
}
func (Wishlist) TableName() string { return "wishlists" }
