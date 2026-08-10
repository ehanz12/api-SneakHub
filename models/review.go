package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Review struct {
	ReviewID   string    `gorm:"column:review_id;type:char(36);primaryKey" json:"review_id"`
	OrderID    string    `gorm:"column:order_id;type:char(36);not null" json:"order_id"`
	CustomerID string    `gorm:"column:customer_id;type:char(36);not null" json:"customer_id"`
	ProductID  string    `gorm:"column:product_id;type:char(36);not null" json:"product_id"`
	Rating     float64   `gorm:"column:rating;type:decimal(3,2);not null" json:"rating"`
	Komentar   *string   `gorm:"column:komentar;type:text" json:"komentar,omitempty"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`

	Order    Order   `gorm:"foreignKey:OrderID;references:OrderID" json:"order,omitempty"`
	Customer User    `gorm:"foreignKey:CustomerID;references:UserID" json:"customer,omitempty"`
	Product  Product `gorm:"foreignKey:ProductID;references:ProductID" json:"product,omitempty"`
}

func (r *Review) BeforeCreate(tx *gorm.DB) error {
	if r.ReviewID == "" {
		r.ReviewID = uuid.NewString()
	}
	return nil
}
func (Review) TableName() string { return "reviews" }
