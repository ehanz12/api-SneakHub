package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ImageEmbedding struct {
	EmbeddingID string         `gorm:"column:embedding_id;type:char(36);primaryKey" json:"embedding_id"`
	ProductID   string         `gorm:"column:product_id;type:char(36);not null" json:"product_id"`
	ImageID     string         `gorm:"column:image_id;type:char(36);not null;unique" json:"image_id"`
	Vector      datatypes.JSON `gorm:"column:vector;type:json;not null" json:"vector"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`

	Product Product      `gorm:"foreignKey:ProductID;references:ProductID" json:"product,omitempty"`
	Image   ProductImage `gorm:"foreignKey:ImageID;references:ImageID" json:"image,omitempty"`
}

func (i *ImageEmbedding) BeforeCreate(tx *gorm.DB) error {
	if i.EmbeddingID == "" {
		i.EmbeddingID = uuid.NewString()
	}
	return nil
}
func (ImageEmbedding) TableName() string { return "image_embeddings" }
