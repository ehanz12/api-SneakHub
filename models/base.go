package models

import "time"

// Timestamps dapat di-embed ketika tabel memiliki created_at/updated_at.
type Timestamps struct {
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}
