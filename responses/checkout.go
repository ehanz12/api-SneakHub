package responses

type CheckoutResponse struct {
	OrderID          string  `json:"order_id"`
	StatusOrder      string  `json:"status_order"`
	Subtotal         float64 `json:"subtotal"`
	BiayaPengiriman  float64 `json:"biaya_pengiriman"`
	TotalPembayaran  float64 `json:"total_pembayaran"`
	MetodePembayaran string  `json:"metode_pembayaran"`
	PaymentURL       string  `json:"payment_url"`
}
