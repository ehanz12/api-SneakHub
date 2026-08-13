package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Payment struct {
	PaymentID            string     `gorm:"column:payment_id;type:char(36);primaryKey" json:"payment_id"`
	OrderID              string     `gorm:"column:order_id;type:char(36);not null;unique" json:"order_id"`
	MetodePembayaran     string     `gorm:"column:metode_pembayaran;type:varchar(50);not null" json:"metode_pembayaran"`
	Jumlah               float64    `gorm:"column:jumlah;type:decimal(15,2);not null" json:"jumlah"`
	StatusPembayaran     string     `gorm:"column:status_pembayaran;type:enum('pending','paid','failed','expired','refunded');not null;default:pending" json:"status_pembayaran"`
	GatewayReference     *string    `gorm:"column:gateway_reference;type:varchar(255)" json:"gateway_reference,omitempty"`
	PaymentURL           *string    `gorm:"column:payment_url;type:varchar(500)" json:"payment_url,omitempty"`
	TransactionReference *string    `gorm:"column:transaction_reference;type:varchar(150)" json:"transaction_reference,omitempty"`
	PaidAt               *time.Time `gorm:"column:paid_at" json:"paid_at,omitempty"`
	CreatedAt            time.Time  `gorm:"column:created_at" json:"created_at"`

	Order Order `gorm:"foreignKey:OrderID;references:OrderID" json:"order,omitempty"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.PaymentID == "" {
		p.PaymentID = uuid.NewString()
	}
	return nil
}
func (Payment) TableName() string { return "payments" }
