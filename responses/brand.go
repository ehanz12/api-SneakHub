package responses

import "time"

type BrandResponse struct {
	BrandID   string    `json:"brand_id"`
	NamaBrand string    `json:"nama_brand"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
