package models

import (
	"time"

	"github.com/google/uuid"
)

// UUIDModel provides a CHAR(36) UUID primary key.
// Every model below uses its own named ID field because the SQL schema
// uses user_id, product_id, order_id, etc.
type UUIDModel struct{}

func NewUUID() string {
	return uuid.NewString()
}

func setUUID(id *string) {
	if id != nil && *id == "" {
		*id = NewUUID()
	}
}

// Timestamps can be embedded when a table has created_at/updated_at.
type Timestamps struct {
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func AutoUUID(id *string) error {
	if *id == "" {
		*id = NewUUID()
	}
	return nil
}
