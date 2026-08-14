package responses

type SmartFilterItemResponse struct {
	ProductID  string   `json:"product_id"`
	NamaProduk string   `json:"nama_produk"`
	Harga      float64  `json:"harga"`
	MatchScore int      `json:"match_score"`
	Alasan     []string `json:"alasan"`
}

type SmartFilterDataResponse struct {
	Items []SmartFilterItemResponse `json:"items"`
}

type HomeProductResponse struct {
	ProductID  string  `json:"product_id"`
	NamaProduk string  `json:"nama_produk"`
	Harga      float64 `json:"harga"`
	ImageURL   string  `json:"image_url"`
}

type HomeSectionResponse struct {
	Type     string                `json:"type"`
	Title    string                `json:"title"`
	Products []HomeProductResponse `json:"products"`
}

type HomePersonalizedDataResponse struct {
	Sections []HomeSectionResponse `json:"sections"`
}

type RecommendationItemResponse struct {
	ProductID  string  `json:"product_id"`
	NamaProduk string  `json:"nama_produk"`
	Harga      float64 `json:"harga"`
	Score      float64 `json:"score"`
	Reason     string  `json:"reason"`
}

type RecommendationDataResponse struct {
	Items []RecommendationItemResponse `json:"items"`
}

type TrendingItemResponse struct {
	ProductID     string `json:"product_id"`
	NamaProduk    string `json:"nama_produk"`
	TrendScore    int    `json:"trend_score"`
	Views         int64  `json:"views"`
	WishlistCount int64  `json:"wishlist_count"`
}

type TrendingDataResponse struct {
	Period string                 `json:"period"`
	Items  []TrendingItemResponse `json:"items"`
}

type BestSellerItemResponse struct {
	Rank         int    `json:"rank"`
	ProductID    string `json:"product_id"`
	NamaProduk   string `json:"nama_produk"`
	TotalTerjual int64  `json:"total_terjual"`
}

type BestSellerDataResponse struct {
	PeriodStart string                   `json:"period_start"`
	PeriodEnd   string                   `json:"period_end"`
	Items       []BestSellerItemResponse `json:"items"`
}
