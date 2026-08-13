package responses

type ReviewResponse struct {
	ReviewID  string  `json:"review_id"`
	OrderID   string  `json:"order_id"`
	ProductID string  `json:"product_id"`
	Rating    float64 `json:"rating"`
	Komentar  *string `json:"komentar,omitempty"`
}
