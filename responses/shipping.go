package responses

type ShippingOptionResponse struct {
	Kurir      string  `json:"kurir"`
	Service    string  `json:"service"`
	Biaya      float64 `json:"biaya"`
	Estimasi   string  `json:"estimasi"`
	IsFallback bool    `json:"is_fallback"`
}

type SellerShippingRatesResponse struct {
	SellerID string                   `json:"seller_id"`
	NamaToko string                   `json:"nama_toko"`
	Berat    int                      `json:"berat"`
	KotaAsal *string                  `json:"kota_asal,omitempty"`
	Options  []ShippingOptionResponse `json:"options"`
}
