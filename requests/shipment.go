package requests

// ShipOrderRequest adalah body untuk pengiriman pesanan oleh seller.
// NomorResi opsional: jika kosong, resi dibuat otomatis via booking Biteship.
type ShipOrderRequest struct {
	NomorResi *string `json:"nomor_resi"`
}
