package responses

import "time"

type CategoryResponse struct {
	CategoryID   string    `json:"cateogry_id"`
	NamaKategori string    `json:"nama_kategori"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
