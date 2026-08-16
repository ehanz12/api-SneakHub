package responses

import "time"

type ReviewResponse struct {
	ReviewID  string  `json:"review_id"`
	OrderID   string  `json:"order_id"`
	ProductID string  `json:"product_id"`
	Rating    float64 `json:"rating"`
	Komentar  *string `json:"komentar,omitempty"`
}

type ReviewListItemResponse struct {
	ReviewID  string                 `json:"review_id"`
	ProductID string                 `json:"product_id"`
	Customer  SellerCustomerResponse `json:"customer"`
	Rating    float64                `json:"rating"`
	Komentar  *string                `json:"komentar,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
}

type ReviewListDataResponse struct {
	Items          []ReviewListItemResponse `json:"items"`
	RatingRataRata float64                  `json:"rating_rata_rata"`
	TotalReview    int64                    `json:"total_review"`
	Pagination     PaginationResponse       `json:"pagination"`
}
