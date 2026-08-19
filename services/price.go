package services

import (
	"encoding/json"
	"errors"
	"math"
	"strings"

	"github.com/ehanz12/api-SneakHub/database"
	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/requests"
	"github.com/ehanz12/api-SneakHub/responses"
)

type marketStats struct {
	Min     float64 `gorm:"column:min_price"`
	Max     float64 `gorm:"column:max_price"`
	Avg     float64 `gorm:"column:avg_price"`
	Count   int64   `gorm:"column:total_count"`
	HasData bool
}

func getMarketStats(brandID, model, kondisi, ukuran string) (*marketStats, error) {
	query := database.DB.Model(&models.MarketPriceData{}).Where("brand_id = ?", brandID)
	if model != "" {
		query = query.Where("nama_model LIKE ?", "%"+model+"%")
	}
	if kondisi != "" {
		query = query.Where("kondisi = ?", kondisi)
	}
	if ukuran != "" {
		query = query.Where("ukuran = ?", ukuran)
	}

	var s marketStats
	if err := query.Select("MIN(harga) AS min_price, MAX(harga) AS max_price, AVG(harga) AS avg_price, COUNT(*) AS total_count").
		Scan(&s).Error; err != nil {
		return nil, errors.New("gagal memuat data pasar")
	}
	if s.Count == 0 {
		return nil, nil
	}
	s.HasData = true
	return &s, nil
}

func productPriceFallback(brandID string) (*marketStats, error) {
	var s marketStats
	if err := database.DB.Model(&models.Product{}).
		Where("brand_id = ? AND status_publikasi = ?", brandID, "aktif").
		Select("MIN(harga) AS min_price, MAX(harga) AS max_price, AVG(harga) AS avg_price, COUNT(*) AS total_count").
		Scan(&s).Error; err != nil {
		return nil, errors.New("gagal memuat harga produk")
	}
	if s.Count == 0 {
		return nil, nil
	}
	s.HasData = true
	return &s, nil
}

func confidenceFor(count int64) float64 {
	switch {
	case count >= 30:
		return 0.9
	case count >= 10:
		return 0.75
	case count >= 3:
		return 0.6
	default:
		return 0.4
	}
}

func PredictPriceService(userID string, r requests.PricePredictionRequest) (*responses.PricePredictionResponse, error) {
	seller, err := findSellerByUserID(database.DB, userID)
	if err != nil {
		return nil, err
	}

	kondisi := normalizeKondisi(r.Kondisi)
	model := strings.TrimSpace(r.Model)

	stats, err := getMarketStats(r.BrandID, model, kondisi, strings.TrimSpace(r.Ukuran))
	if err != nil {
		return nil, err
	}
	if stats == nil {
		stats, err = getMarketStats(r.BrandID, model, "", "")
		if err != nil {
			return nil, err
		}
	}
	if stats == nil {
		stats, err = getMarketStats(r.BrandID, "", "", "")
		if err != nil {
			return nil, err
		}
	}
	if stats == nil {
		stats, err = productPriceFallback(r.BrandID)
		if err != nil {
			return nil, err
		}
	}
	if stats == nil {
		return nil, errors.New("data pasar belum tersedia untuk prediksi")
	}

	scoreFactor := r.ConditionScore / 100
	recommended := stats.Avg * (0.7 + 0.3*scoreFactor)
	if recommended < stats.Min {
		recommended = stats.Min
	}
	if recommended > stats.Max {
		recommended = stats.Max
	}

	inputJSON, err := json.Marshal(r)
	if err != nil {
		return nil, errors.New("gagal menyimpan atribut input")
	}
	prediction := models.PricePrediction{
		SellerID:          seller.SellerID,
		InputAtribut:      inputJSON,
		EstimatedPriceMin: &stats.Min,
		EstimatedPriceMax: &stats.Max,
		RecommendedPrice:  &recommended,
	}
	if err := database.DB.Create(&prediction).Error; err != nil {
		return nil, errors.New("gagal menyimpan prediksi harga")
	}

	return &responses.PricePredictionResponse{
		EstimatedPriceMin: math.Round(stats.Min),
		EstimatedPriceMax: math.Round(stats.Max),
		RecommendedPrice:  math.Round(recommended),
		Confidence:        confidenceFor(stats.Count),
	}, nil
}

func PriceInsightService(productID string) (*responses.PriceInsightResponse, error) {
	var product models.Product
	if err := database.DB.Select("product_id", "brand_id", "nama_produk", "harga").
		Where("product_id = ? AND status_publikasi = ?", productID, "aktif").
		First(&product).Error; err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}

	stats, err := getMarketStats(product.BrandID, "", "", "")
	if err != nil {
		return nil, err
	}

	insight := &responses.PriceInsightResponse{CurrentPrice: product.Harga}
	if stats == nil {
		insight.MarketPriceMin = product.Harga
		insight.MarketPriceMax = product.Harga
		insight.MarketAverage = product.Harga
		insight.Message = "Data pasar belum tersedia, harga produk dijadikan acuan."
		return insight, nil
	}

	diffPercent := (product.Harga - stats.Avg) / stats.Avg * 100
	insight.MarketPriceMin = math.Round(stats.Min)
	insight.MarketPriceMax = math.Round(stats.Max)
	insight.MarketAverage = math.Round(stats.Avg*100) / 100
	insight.PriceDifferencePercent = math.Round(diffPercent*100) / 100

	if math.Abs(diffPercent) >= 10 {
		insight.Anomaly = true
		if diffPercent < 0 {
			insight.AnomalyType = "LOW"
			insight.Message = "Harga produk berada di bawah rata-rata pasar."
		} else {
			insight.AnomalyType = "HIGH"
			insight.Message = "Harga produk berada di atas rata-rata pasar."
		}
	} else {
		insight.Message = "Harga produk sesuai dengan rata-rata pasar."
	}

	return insight, nil
}
