package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConditionScore struct {
	ConditionID string    `gorm:"column:condition_id;type:char(36);primaryKey" json:"condition_id"`
	ProductID   string    `gorm:"column:product_id;type:char(36);not null;unique" json:"product_id"`
	Upper       *float64  `gorm:"column:upper;type:decimal(5,2)" json:"upper,omitempty"`
	Outsole     *float64  `gorm:"column:outsole;type:decimal(5,2)" json:"outsole,omitempty"`
	Midsole     *float64  `gorm:"column:midsole;type:decimal(5,2)" json:"midsole,omitempty"`
	Insole      *float64  `gorm:"column:insole;type:decimal(5,2)" json:"insole,omitempty"`
	Accessories *float64  `gorm:"column:accessories;type:decimal(5,2)" json:"accessories,omitempty"`
	Box         *float64  `gorm:"column:box;type:decimal(5,2)" json:"box,omitempty"`
	SkorAkhir   *float64  `gorm:"column:skor_akhir;type:decimal(5,2)" json:"skor_akhir,omitempty"`
	DinilaiOleh string    `gorm:"column:dinilai_oleh;type:enum('seller','admin');not null" json:"dinilai_oleh"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`

	Product Product `gorm:"foreignKey:ProductID;references:ProductID" json:"product,omitempty"`
}

func (c *ConditionScore) BeforeCreate(tx *gorm.DB) error {
	if c.ConditionID == "" {
		c.ConditionID = uuid.NewString()
	}
	return nil
}
func (ConditionScore) TableName() string { return "condition_scores" }
