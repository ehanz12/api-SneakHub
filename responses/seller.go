package responses

type SellerResponse struct {
	SellerID         string `json:"seller_id"`
	UserID           string `json:"user_id"`
	NamaToko         string `json:"nama_toko"`
	StatusVerifikasi string `json:"status_verifikasi"`
}
