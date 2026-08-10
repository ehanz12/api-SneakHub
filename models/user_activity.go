package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type UserActivity struct {
	ActivityID     string         `gorm:"column:activity_id;type:char(36);primaryKey" json:"activity_id"`
	UserID         string         `gorm:"column:user_id;type:char(36);not null" json:"user_id"`
	JenisAktivitas string         `gorm:"column:jenis_aktivitas;type:enum('search','view','add_to_wishlist','add_to_cart','purchase','review','dll');not null" json:"jenis_aktivitas"`
	ProductID      *string        `gorm:"column:product_id;type:char(36)" json:"product_id,omitempty"`
	Metadata       datatypes.JSON `gorm:"column:metadata;type:json" json:"metadata,omitempty"`
	CreatedAt      time.Time      `gorm:"column:created_at" json:"created_at"`

	User    User     `gorm:"foreignKey:UserID;references:UserID" json:"user,omitempty"`
	Product *Product `gorm:"foreignKey:ProductID;references:ProductID" json:"product,omitempty"`
}

func (u *UserActivity) BeforeCreate(tx *gorm.DB) error {
	if u.ActivityID == "" {
		u.ActivityID = uuid.NewString()
	}
	return nil
}
func (UserActivity) TableName() string { return "user_activities" }
