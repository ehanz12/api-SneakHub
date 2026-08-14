package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Wishlist struct {
	WishlistID          string    `gorm:"column:wishlist_id;type:char(36);primaryKey" json:"wishlist_id"`
	CustomerID          string    `gorm:"column:customer_id;type:char(36);not null" json:"customer_id"`
	ProductID           string    `gorm:"column:product_id;type:char(36);not null" json:"product_id"`
	StatusStok          string    `gorm:"column:status_stok;type:enum('available','out_of_stock');default:available" json:"status_stok"`
	PriceAlertEnabled   bool      `gorm:"column:price_alert_enabled;not null;default:false" json:"price_alert_enabled"`
	TargetPrice         *float64  `gorm:"column:target_price;type:decimal(15,2)" json:"target_price,omitempty"`
	RestockAlertEnabled bool      `gorm:"column:restock_alert_enabled;not null;default:false" json:"restock_alert_enabled"`
	CreatedAt           time.Time `gorm:"column:created_at" json:"created_at"`

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
