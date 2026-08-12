package responses

type CartItemResponse struct {
	CartItemID string  `json:"cart_item_id"`
	ProductID  string  `json:"product_id"`
	Jumlah     int     `json:"jumlah"`
	Subtotal   float64 `json:"subtotal"`
}

type CartItemsResponse struct {
	Items []CartItemResponse `json:"items"`
	Total float64            `json:"total"`
}

type CartItemDetailResponse struct {
	CartItemResponse
	NamaProduk string  `json:"nama_produk"`
	Harga      float64 `json:"harga"`
	ImageURL   string  `json:"image_url"`
}

type CartResponse struct {
	CartID    string                   `json:"cart_id"`
	Items     []CartItemDetailResponse `json:"items"`
	TotalItem int                      `json:"total_item"`
	Total     float64                  `json:"total"`
}
