package responses

import "time"

type OrderListItemResponse struct {
	OrderID         string    `json:"order_id"`
	SellerID        string    `json:"seller_id"`
	StatusOrder     string    `json:"status_order"`
	TotalPembayaran float64   `json:"total_pembayaran"`
	CreatedAt       time.Time `json:"created_at"`
}

type OrderAlamatResponse struct {
	NamaPenerima string `json:"nama_penerima"`
	NomorTelepon string `json:"nomor_telepon"`
	Alamat       string `json:"alamat"`
	Kota         string `json:"kota"`
	Provinsi     string `json:"provinsi"`
	KodePos      string `json:"kode_pos"`
}

type OrderItemResponse struct {
	OrderItemID        string  `json:"order_item_id"`
	ProductID          string  `json:"product_id"`
	NamaProduk         string  `json:"nama_produk"`
	Jumlah             int     `json:"jumlah"`
	HargaSaatTransaksi float64 `json:"harga_saat_transaksi"`
}

type OrderDetailResponse struct {
	OrderID          string              `json:"order_id"`
	CustomerID       string              `json:"customer_id"`
	SellerID         string              `json:"seller_id"`
	StatusOrder      string              `json:"status_order"`
	AlamatPengiriman OrderAlamatResponse `json:"alamat_pengiriman"`
	MetodePembayaran string              `json:"metode_pembayaran"`
	Items            []OrderItemResponse `json:"items"`
	Subtotal         float64             `json:"subtotal"`
	BiayaPengiriman  float64             `json:"biaya_pengiriman"`
	TotalPembayaran  float64             `json:"total_pembayaran"`
	CreatedAt        time.Time           `json:"created_at"`
}

type OrderListDataResponse struct {
	Items      []OrderListItemResponse `json:"items"`
	Pagination PaginationResponse      `json:"pagination"`
}
