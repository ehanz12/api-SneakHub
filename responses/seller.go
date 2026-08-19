package responses

import "time"

type SellerResponse struct {
	SellerID         string  `json:"seller_id"`
	UserID           string  `json:"user_id"`
	NamaToko         string  `json:"nama_toko"`
	DeskripsiToko    *string `json:"deskripsi_toko,omitempty"`
	StatusVerifikasi string  `json:"status_verifikasi"`
	KodePosAsal      *string `json:"kode_pos_asal,omitempty"`
	KotaAsal         *string `json:"kota_asal,omitempty"`
	AlamatAsal       *string `json:"alamat_asal,omitempty"`
}

type SellerProductListItemResponse struct {
	ProductID       string  `json:"product_id"`
	NamaProduk      string  `json:"nama_produk"`
	Harga           float64 `json:"harga"`
	ImageURL        string  `json:"image_url"`
	Stok            int     `json:"stok"`
	StatusPublikasi string  `json:"status_publikasi"`
	TotalTerjual    int64   `json:"total_terjual"`
}

type SellerProductListDataResponse struct {
	Items      []SellerProductListItemResponse `json:"items"`
	Pagination PaginationResponse              `json:"pagination"`
}

type SellerCustomerResponse struct {
	UserID string `json:"user_id"`
	Nama   string `json:"nama"`
}

type SellerOrderListItemResponse struct {
	OrderID         string                 `json:"order_id"`
	Customer        SellerCustomerResponse `json:"customer"`
	StatusOrder     string                 `json:"status_order"`
	TotalPembayaran float64                `json:"total_pembayaran"`
	CreatedAt       time.Time              `json:"created_at"`
}

type SellerOrderListDataResponse struct {
	Items      []SellerOrderListItemResponse `json:"items"`
	Pagination PaginationResponse            `json:"pagination"`
}

type SellerTopProductResponse struct {
	ProductID    string `json:"product_id"`
	NamaProduk   string `json:"nama_produk"`
	ImageURL     string `json:"image_url"`
	TotalTerjual int64  `json:"total_terjual"`
}

type SellerDashboardResponse struct {
	TotalProduk      int64                      `json:"total_produk"`
	ProdukAktif      int64                      `json:"produk_aktif"`
	TotalTerjual     int64                      `json:"total_terjual"`
	TotalPendapatan  float64                    `json:"total_pendapatan"`
	RatingRataRata   float64                    `json:"rating_rata_rata"`
	SellerTrustScore *float64                   `json:"seller_trust_score"`
	ProdukTerlaris   []SellerTopProductResponse `json:"produk_terlaris"`
}
