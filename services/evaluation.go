package services

import (
	"errors"
	"math"
	"strings"

	"github.com/ehanz12/api-SneakHub/database"
	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/requests"
	"github.com/ehanz12/api-SneakHub/responses"
)

func fptr(f float64) *float64 {
	return &f
}

func fval(v *float64, fallback float64) float64 {
	if v == nil {
		return fallback
	}
	return *v
}

func CreateConditionScoreService(userID, role, productID string, r requests.ConditionScoreRequest) (*models.ConditionScore, error) {
	tx := database.DB.Begin()
	if tx.Error != nil {
		return nil, errors.New("transaksi database gagal dimulai")
	}

	var product models.Product
	if err := tx.Select("product_id", "seller_id").
		Where("product_id = ?", productID).First(&product).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("produk tidak ditemukan")
	}

	if role == "seller" {
		var seller models.Seller
		if err := tx.Select("seller_id").Where("user_id = ?", userID).First(&seller).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("data toko seller tidak ditemukan")
		}
		if seller.SellerID != product.SellerID {
			tx.Rollback()
			return nil, errors.New("produk bukan milik seller")
		}
	}

	dinilaiOleh := strings.ToLower(strings.TrimSpace(r.DinilaiOleh))
	skorAkhir := r.Upper*0.2 + r.Outsole*0.15 + r.Midsole*0.15 +
		r.Insole*0.2 + r.Accessories*0.15 + r.Box*0.15

	var cs models.ConditionScore
	err := tx.Where("product_id = ?", productID).First(&cs).Error
	if err == nil {
		cs.Upper = fptr(r.Upper)
		cs.Outsole = fptr(r.Outsole)
		cs.Midsole = fptr(r.Midsole)
		cs.Insole = fptr(r.Insole)
		cs.Accessories = fptr(r.Accessories)
		cs.Box = fptr(r.Box)
		cs.SkorAkhir = fptr(skorAkhir)
		cs.DinilaiOleh = dinilaiOleh
		if err := tx.Save(&cs).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("gagal memperbarui condition score")
		}
	} else {
		cs = models.ConditionScore{
			ProductID:   productID,
			Upper:       fptr(r.Upper),
			Outsole:     fptr(r.Outsole),
			Midsole:     fptr(r.Midsole),
			Insole:      fptr(r.Insole),
			Accessories: fptr(r.Accessories),
			Box:         fptr(r.Box),
			SkorAkhir:   fptr(skorAkhir),
			DinilaiOleh: dinilaiOleh,
		}
		if err := tx.Create(&cs).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("gagal menyimpan condition score")
		}
	}

	if err := tx.Model(&models.Product{}).Where("product_id = ?", productID).
		Update("condition_score", skorAkhir).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal memperbarui skor produk")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New("perubahan gagal disimpan")
	}

	return &cs, nil
}

func GetConditionScoreService(productID string) (*models.ConditionScore, error) {
	var cs models.ConditionScore
	if err := database.DB.Where("product_id = ?", productID).First(&cs).Error; err != nil {
		return nil, errors.New("condition score belum tersedia untuk produk ini")
	}
	return &cs, nil
}

func GetSellerTrustScoreService(sellerID string) (*responses.SellerTrustScoreResponse, error) {
	var seller models.Seller
	if err := database.DB.Select("seller_id", "seller_trust_score").
		Where("seller_id = ?", sellerID).First(&seller).Error; err != nil {
		return nil, errors.New("seller tidak ditemukan")
	}

	var ts models.SellerTrustScore
	if err := database.DB.Where("seller_id = ?", sellerID).First(&ts).Error; err == nil {
		return &responses.SellerTrustScoreResponse{
			SellerID:            sellerID,
			SkorAkhir:           fval(ts.SkorAkhir, fval(seller.SellerTrustScore, 0)),
			OrderCompletionRate: fval(ts.OrderCompletionRate, 0),
			AverageRating:       fval(ts.AverageRating, 0),
			CancellationRate:    fval(ts.CancellationRate, 0),
			ResponseRate:        fval(ts.ResponseRate, 0),
		}, nil
	}

	var selesai, dibatalkan int64
	if err := database.DB.Model(&models.Order{}).
		Where("seller_id = ? AND status_order = ?", sellerID, "selesai").
		Count(&selesai).Error; err != nil {
		return nil, errors.New("gagal memuat data order")
	}
	if err := database.DB.Model(&models.Order{}).
		Where("seller_id = ? AND status_order = ?", sellerID, "dibatalkan").
		Count(&dibatalkan).Error; err != nil {
		return nil, errors.New("gagal memuat data order")
	}

	completionRate, cancellationRate := 0.0, 0.0
	total := selesai + dibatalkan
	if total > 0 {
		completionRate = float64(selesai) / float64(total) * 100
		cancellationRate = float64(dibatalkan) / float64(total) * 100
	}

	var avgRating float64
	if err := database.DB.Model(&models.Review{}).
		Select("COALESCE(AVG(reviews.rating), 0)").
		Joins("JOIN products ON products.product_id = reviews.product_id").
		Where("products.seller_id = ?", sellerID).
		Scan(&avgRating).Error; err != nil {
		return nil, errors.New("gagal memuat rating")
	}

	return &responses.SellerTrustScoreResponse{
		SellerID:            sellerID,
		SkorAkhir:           fval(seller.SellerTrustScore, 0),
		OrderCompletionRate: math.Round(completionRate*100) / 100,
		AverageRating:       math.Round(avgRating*100) / 100,
		CancellationRate:    math.Round(cancellationRate*100) / 100,
		ResponseRate:        0,
	}, nil
}
