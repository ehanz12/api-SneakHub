package responses

type ConditionScoreDetailResponse struct {
	Upper       float64 `json:"upper"`
	Outsole     float64 `json:"outsole"`
	Midsole     float64 `json:"midsole"`
	Insole      float64 `json:"insole"`
	Accessories float64 `json:"accessories"`
	Box         float64 `json:"box"`
}

type ConditionScoreCreateResponse struct {
	ProductID   string  `json:"product_id"`
	SkorAkhir   float64 `json:"skor_akhir"`
	Upper       float64 `json:"upper"`
	Outsole     float64 `json:"outsole"`
	Midsole     float64 `json:"midsole"`
	Insole      float64 `json:"insole"`
	Accessories float64 `json:"accessories"`
	Box         float64 `json:"box"`
	DinilaiOleh string  `json:"dinilai_oleh"`
}

type ConditionScoreGetResponse struct {
	ProductID   string                       `json:"product_id"`
	SkorAkhir   float64                      `json:"skor_akhir"`
	Detail      ConditionScoreDetailResponse `json:"detail"`
	DinilaiOleh string                       `json:"dinilai_oleh"`
}

type SellerTrustScoreResponse struct {
	SellerID            string  `json:"seller_id"`
	SkorAkhir           float64 `json:"skor_akhir"`
	OrderCompletionRate float64 `json:"order_completion_rate"`
	AverageRating       float64 `json:"average_rating"`
	CancellationRate    float64 `json:"cancellation_rate"`
	ResponseRate        float64 `json:"response_rate"`
}
