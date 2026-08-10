package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductImage struct {
	ImageID          string    `gorm:"column:image_id;type:char(36);primaryKey" json:"image_id"`
	ProductID        string    `gorm:"column:product_id;type:char(36);not null" json:"product_id"`
	URLObjectStorage string    `gorm:"column:url_object_storage;type:text;not null" json:"url_object_storage"`
	UrutanTampil     int       `gorm:"column:urutan_tampil;not null;default:0" json:"urutan_tampil"`
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at"`

	Product   Product         `gorm:"foreignKey:ProductID;references:ProductID" json:"product,omitempty"`
	Embedding *ImageEmbedding `gorm:"foreignKey:ImageID;references:ImageID" json:"embedding,omitempty"`
}

func (p *ProductImage) BeforeCreate(tx *gorm.DB) error {
	if p.ImageID == "" {
		p.ImageID = uuid.NewString()
	}
	return nil
}
func (ProductImage) TableName() string { return "product_images" }
