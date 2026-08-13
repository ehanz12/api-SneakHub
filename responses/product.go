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

type SellerInfoResponse struct {
	SellerID         string   `json:"seller_id"`
	NamaToko         string   `json:"nama_toko"`
	SellerTrustScore *float64 `json:"seller_trust_score"`
}

type ProductListItemResponse struct {
	ProductID      string             `json:"product_id"`
	NamaProduk     string             `json:"nama_produk"`
	Harga          float64            `json:"harga"`
	Kondisi        string             `json:"kondisi"`
	Stok           int                `json:"stok"`
	UkuranTersedia []string           `json:"ukuran_tersedia"`
	ConditionScore *float64           `json:"condition_score"`
	ImageURL       string             `json:"image_url"`
	Seller         SellerInfoResponse `json:"seller"`
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type ProductListDataResponse struct {
	Items      []ProductListItemResponse `json:"items"`
	Pagination PaginationResponse        `json:"pagination"`
}

type ProductDetailImageResponse struct {
	ImageID      string `json:"image_id"`
	URL          string `json:"url"`
	UrutanTampil int    `json:"urutan_tampil"`
}

type ProductDetailResponse struct {
	ProductID       string                       `json:"product_id"`
	SellerID        string                       `json:"seller_id"`
	NamaProduk      string                       `json:"nama_produk"`
	BrandID         string                       `json:"brand_id"`
	CategoryID      string                       `json:"category_id"`
	Kondisi         string                       `json:"kondisi"`
	Deskripsi       *string                      `json:"deskripsi,omitempty"`
	Harga           float64                      `json:"harga"`
	Stok            int                          `json:"stok"`
	UkuranTersedia  []string                     `json:"ukuran_tersedia"`
	ConditionScore  *float64                     `json:"condition_score"`
	StatusPublikasi string                       `json:"status_publikasi"`
	Images          []ProductDetailImageResponse `json:"images"`
}

type SearchProductByImageResponse struct {
	ProductID     string  `json:"product_id"`
	NamaProduk    string  `json:"nama_produk"`
	Harga         float64 `json:"harga"`
	ImageURL      string  `json:"image_url"`
	SkorKemiripan float64 `json:"skor_kemiripan"`
}
