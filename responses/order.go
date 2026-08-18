package responses

import "time"

type OrderListItemResponse struct {
	OrderID          string    `json:"order_id"`
	SellerID         string    `json:"seller_id"`
	StatusOrder      string    `json:"status_order"`
	TotalPembayaran  float64   `json:"total_pembayaran"`
	StatusPembayaran string    `json:"status_pembayaran"`
	CreatedAt        time.Time `json:"created_at"`
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

type OrderPaymentResponse struct {
	PaymentID            string     `json:"payment_id"`
	MetodePembayaran     string     `json:"metode_pembayaran"`
	Jumlah               float64    `json:"jumlah"`
	StatusPembayaran     string     `json:"status_pembayaran"`
	PaymentURL           *string    `json:"payment_url,omitempty"`
	GatewayReference     *string    `json:"gateway_reference,omitempty"`
	TransactionReference *string    `json:"transaction_reference,omitempty"`
	PaidAt               *time.Time `json:"paid_at,omitempty"`
}

type ShipmentTrackingEventResponse struct {
	Note      string `json:"note"`
	Status    string `json:"status"`
	UpdatedAt string `json:"updated_at"`
}

type OrderShipmentResponse struct {
	ShipmentID       string                         `json:"shipment_id"`
	Kurir            string                         `json:"kurir"`
	Service          *string                        `json:"service,omitempty"`
	NomorResi        *string                        `json:"nomor_resi,omitempty"`
	TrackingID       *string                        `json:"tracking_id,omitempty"`
	TrackingHistory  []ShipmentTrackingEventResponse `json:"tracking_history,omitempty"`
	StatusPengiriman string                         `json:"status_pengiriman"`
	ShippedAt        *time.Time                     `json:"shipped_at,omitempty"`
	DeliveredAt      *time.Time                     `json:"delivered_at,omitempty"`
}

type OrderDetailResponse struct {
	OrderID          string                 `json:"order_id"`
	CustomerID       string                 `json:"customer_id"`
	SellerID         string                 `json:"seller_id"`
	StatusOrder      string                 `json:"status_order"`
	AlamatPengiriman OrderAlamatResponse    `json:"alamat_pengiriman"`
	MetodePembayaran string                 `json:"metode_pembayaran"`
	Items            []OrderItemResponse    `json:"items"`
	Subtotal         float64                `json:"subtotal"`
	BiayaPengiriman  float64                `json:"biaya_pengiriman"`
	TotalPembayaran  float64                `json:"total_pembayaran"`
	Payment          *OrderPaymentResponse  `json:"payment,omitempty"`
	Shipment         *OrderShipmentResponse `json:"shipment,omitempty"`
	CreatedAt        time.Time              `json:"created_at"`
}

type OrderListDataResponse struct {
	Items      []OrderListItemResponse `json:"items"`
	Pagination PaginationResponse      `json:"pagination"`
}
