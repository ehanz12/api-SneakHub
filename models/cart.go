package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Cart struct {
	CartID     string    `gorm:"column:cart_id;type:char(36);primaryKey" json:"cart_id"`
	CustomerID string    `gorm:"column:customer_id;type:char(36);not null;unique" json:"customer_id"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`

	Customer User       `gorm:"foreignKey:CustomerID;references:UserID" json:"customer,omitempty"`
	Items    []CartItem `gorm:"foreignKey:CartID;references:CartID" json:"items,omitempty"`
}

func (c *Cart) BeforeCreate(tx *gorm.DB) error {
	if c.CartID == "" {
		c.CartID = uuid.NewString()
	}
	return nil
}
func (Cart) TableName() string { return "carts" }
