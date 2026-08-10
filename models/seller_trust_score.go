package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SellerTrustScore struct {
	ScoreID             string    `gorm:"column:score_id;type:char(36);primaryKey" json:"score_id"`
	SellerID            string    `gorm:"column:seller_id;type:char(36);not null;unique" json:"seller_id"`
	OrderCompletionRate *float64  `gorm:"column:order_completion_rate;type:decimal(5,2)" json:"order_completion_rate,omitempty"`
	AverageRating       *float64  `gorm:"column:average_rating;type:decimal(3,2)" json:"average_rating,omitempty"`
	CancellationRate    *float64  `gorm:"column:cancellation_rate;type:decimal(5,2)" json:"cancellation_rate,omitempty"`
	ResponseRate        *float64  `gorm:"column:response_rate;type:decimal(5,2)" json:"response_rate,omitempty"`
	SkorAkhir           *float64  `gorm:"column:skor_akhir;type:decimal(5,2)" json:"skor_akhir,omitempty"`
	CalculatedAt        time.Time `gorm:"column:calculated_at" json:"calculated_at"`
	Seller              Seller    `gorm:"foreignKey:SellerID;references:SellerID" json:"seller,omitempty"`
}

func (s *SellerTrustScore) BeforeCreate(tx *gorm.DB) error {
	if s.ScoreID == "" {
		s.ScoreID = uuid.NewString()
	}
	return nil
}
func (SellerTrustScore) TableName() string { return "seller_trust_scores" }
