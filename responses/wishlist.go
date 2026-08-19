package responses

type WishlistPriceAlertResponse struct {
	Enabled     bool     `json:"enabled"`
	TargetPrice *float64 `json:"target_price,omitempty"`
}

type WishlistRestockAlertResponse struct {
	Enabled bool `json:"enabled"`
}

type WishlistItemResponse struct {
	WishlistID         string                       `json:"wishlist_id"`
	ProductID          string                       `json:"product_id"`
	NamaProduk         string                       `json:"nama_produk"`
	Harga              float64                      `json:"harga"`
	ImageURL           string                       `json:"image_url"`
	StatusStokTerakhir string                       `json:"status_stok_terakhir"`
	PriceAlert         WishlistPriceAlertResponse   `json:"price_alert"`
	RestockAlert       WishlistRestockAlertResponse `json:"restock_alert"`
}

type CreateWishlistResponse struct {
	WishlistID string `json:"wishlist_id"`
	ProductID  string `json:"product_id"`
}

type PriceAlertResponse struct {
	ProductID   string   `json:"product_id"`
	Enabled     bool     `json:"enabled"`
	TargetPrice *float64 `json:"target_price,omitempty"`
}

type RestockAlertResponse struct {
	ProductID string `json:"product_id"`
	Enabled   bool   `json:"enabled"`
}
