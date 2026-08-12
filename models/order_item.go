package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderItem struct {
	OrderItemID        string  `gorm:"column:order_item_id;type:char(36);primaryKey" json:"order_item_id"`
	OrderID            string  `gorm:"column:order_id;type:char(36);not null" json:"order_id"`
	ProductID          string  `gorm:"column:product_id;type:char(36);not null" json:"product_id"`
	Jumlah             int     `gorm:"column:jumlah;not null" json:"jumlah"`
	HargaSaatTransaksi float64 `gorm:"column:harga_saat_transaksi;type:decimal(15,2);not null" json:"harga_saat_transaksi"`

	Order   Order   `gorm:"foreignKey:OrderID;references:OrderID" json:"order,omitempty"`
	Product Product `gorm:"foreignKey:ProductID;references:ProductID" json:"product,omitempty"`
}

func (o *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if o.OrderItemID == "" {
		o.OrderItemID = uuid.NewString()
	}
	return nil
}
func (OrderItem) TableName() string { return "order_items" }
