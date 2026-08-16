package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Seller struct {
	SellerID         string   `gorm:"column:seller_id;type:char(36);primaryKey" json:"seller_id"`
	UserID           string   `gorm:"column:user_id;type:char(36);not null;unique" json:"user_id"`
	NamaToko         string   `gorm:"column:nama_toko;type:varchar(150);not null" json:"nama_toko"`
	DeskripsiToko    *string  `gorm:"column:deskripsi_toko;type:text" json:"deskripsi_toko,omitempty"`
	StatusVerifikasi string   `gorm:"column:status_verifikasi;type:enum('pending','verified','rejected');default:pending" json:"status_verifikasi"`
	KodePosAsal      *string  `gorm:"column:kode_pos_asal;type:varchar(10)" json:"kode_pos_asal,omitempty"`
	KotaAsal         *string  `gorm:"column:kota_asal;type:varchar(100)" json:"kota_asal,omitempty"`
	SellerTrustScore *float64 `gorm:"column:seller_trust_score;type:decimal(5,2)" json:"seller_trust_score,omitempty"`
	Timestamps

	User             User              `gorm:"foreignKey:UserID;references:UserID" json:"user,omitempty"`
	TrustScore       *SellerTrustScore `gorm:"foreignKey:SellerID;references:SellerID" json:"trust_score,omitempty"`
	Products         []Product         `gorm:"foreignKey:SellerID;references:SellerID" json:"products,omitempty"`
	Orders           []Order           `gorm:"foreignKey:SellerID;references:SellerID" json:"orders,omitempty"`
	PricePredictions []PricePrediction `gorm:"foreignKey:SellerID;references:SellerID" json:"price_predictions,omitempty"`
}

func (s *Seller) BeforeCreate(tx *gorm.DB) error {
	if s.SellerID == "" {
		s.SellerID = uuid.NewString()
	}
	return nil
}

func (Seller) TableName() string { return "sellers" }
