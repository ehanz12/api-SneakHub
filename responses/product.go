package responses

import "time"

type CreateProductResponse struct {
	ProductID       string `json:"product_id"`
	StatusPublikasi string `json:"status_publikasi"`
}

type UpdateProductResponse struct {
	ProductID string    `json:"product_id"`
	Harga     float64   `json:"harga"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductImageResponse struct {
	ImageID      string    `json:"image_id"`
	ProductID    string    `json:"product_id"`
	URL          string    `json:"url"`
	UrutanTampil int       `json:"urutan_tampil"`
	Embedding    []float64 `json:"embedding,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type SearchProductByImageResponse struct {
	ProductID     string  `json:"product_id"`
	NamaProduk    string  `json:"nama_produk"`
	Harga         float64 `json:"harga"`
	ImageURL      string  `json:"image_url"`
	SkorKemiripan float64 `json:"skor_kemiripan"`
}
