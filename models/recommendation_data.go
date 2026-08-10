package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type RecommendationData struct {
	RecommendationID string         `gorm:"column:recommendation_id;type:char(36);primaryKey" json:"recommendation_id"`
	UserID           string         `gorm:"column:user_id;type:char(36);not null" json:"user_id"`
	Sumber           string         `gorm:"column:sumber;type:enum('personalized','trending','best_seller','dll');not null" json:"sumber"`
	DaftarProductID  datatypes.JSON `gorm:"column:daftar_product_id;type:json;not null" json:"daftar_product_id"`
	GeneratedAt      time.Time      `gorm:"column:generated_at" json:"generated_at"`

	User User `gorm:"foreignKey:UserID;references:UserID" json:"user,omitempty"`
}

func (r *RecommendationData) BeforeCreate(tx *gorm.DB) error {
	if r.RecommendationID == "" {
		r.RecommendationID = uuid.NewString()
	}
	return nil
}
func (RecommendationData) TableName() string { return "recommendation_data" }
