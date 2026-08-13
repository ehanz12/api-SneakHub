package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Order struct {
	OrderID          string         `gorm:"column:order_id;type:char(36);primaryKey" json:"order_id"`
	CustomerID       string         `gorm:"column:customer_id;type:char(36);not null" json:"customer_id"`
	SellerID         string         `gorm:"column:seller_id;type:char(36);not null" json:"seller_id"`
	AddressID        *string        `gorm:"column:address_id;type:char(36)" json:"address_id,omitempty"`
	StatusOrder      string         `gorm:"column:status_order;type:enum('pending','diproses','dikirim','selesai','dibatalkan');default:pending" json:"status_order"`
	AlamatPengiriman datatypes.JSON `gorm:"column:alamat_pengiriman;type:json;not null" json:"alamat_pengiriman"`
	MetodePembayaran string         `gorm:"column:metode_pembayaran;type:varchar(50);not null" json:"metode_pembayaran"`
	Subtotal         float64        `gorm:"column:subtotal;type:decimal(15,2);not null;default:0" json:"subtotal"`
	BiayaPengiriman  float64        `gorm:"column:biaya_pengiriman;type:decimal(15,2);not null;default:0" json:"biaya_pengiriman"`
	TotalPesanan     float64        `gorm:"column:total_pesanan;type:decimal(15,2);not null" json:"total_pesanan"`
	Timestamps

	Customer User        `gorm:"foreignKey:CustomerID;references:UserID" json:"customer,omitempty"`
	Seller   Seller      `gorm:"foreignKey:SellerID;references:SellerID" json:"seller,omitempty"`
	Items    []OrderItem `gorm:"foreignKey:OrderID;references:OrderID" json:"items,omitempty"`
	Payment  *Payment    `gorm:"foreignKey:OrderID;references:OrderID" json:"payment,omitempty"`
	Shipment *Shipment   `gorm:"foreignKey:OrderID;references:OrderID" json:"shipment,omitempty"`
	Reviews  []Review    `gorm:"foreignKey:OrderID;references:OrderID" json:"reviews,omitempty"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.OrderID == "" {
		o.OrderID = uuid.NewString()
	}
	return nil
}
func (Order) TableName() string { return "orders" }
