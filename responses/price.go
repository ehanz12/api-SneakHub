package responses

type PricePredictionResponse struct {
	EstimatedPriceMin float64 `json:"estimated_market_price_min"`
	EstimatedPriceMax float64 `json:"estimated_market_price_max"`
	RecommendedPrice  float64 `json:"recommended_price"`
	Confidence        float64 `json:"confidence"`
}

type PriceInsightResponse struct {
	CurrentPrice           float64 `json:"current_price"`
	MarketPriceMin         float64 `json:"market_price_min"`
	MarketPriceMax         float64 `json:"market_price_max"`
	MarketAverage          float64 `json:"market_average"`
	PriceDifferencePercent float64 `json:"price_difference_percent"`
	Anomaly                bool    `json:"anomaly"`
	AnomalyType            string  `json:"anomaly_type,omitempty"`
	Message                string  `json:"message"`
}
