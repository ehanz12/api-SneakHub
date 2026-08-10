package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type PricePrediction struct {
	PredictionID      string         `gorm:"column:prediction_id;type:char(36);primaryKey" json:"prediction_id"`
	SellerID          string         `gorm:"column:seller_id;type:char(36);not null" json:"seller_id"`
	InputAtribut      datatypes.JSON `gorm:"column:input_atribut;type:json" json:"input_atribut,omitempty"`
	EstimatedPriceMin *float64       `gorm:"column:estimated_price_min;type:decimal(15,2)" json:"estimated_price_min,omitempty"`
	EstimatedPriceMax *float64       `gorm:"column:estimated_price_max;type:decimal(15,2)" json:"estimated_price_max,omitempty"`
	RecommendedPrice  *float64       `gorm:"column:recommended_price;type:decimal(15,2)" json:"recommended_price,omitempty"`
	CreatedAt         time.Time      `gorm:"column:created_at" json:"created_at"`

	Seller Seller `gorm:"foreignKey:SellerID;references:SellerID" json:"seller,omitempty"`
}

func (p *PricePrediction) BeforeCreate(tx *gorm.DB) error {
	if p.PredictionID == "" {
		p.PredictionID = uuid.NewString()
	}
	return nil
}
func (PricePrediction) TableName() string { return "price_predictions" }
