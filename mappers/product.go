package mappers

import (
	"encoding/json"
	"strings"

	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/responses"
	"github.com/ehanz12/api-SneakHub/services"
)

func ToProductResponse(p models.Product) responses.CreateProductResponse {
	return responses.CreateProductResponse{
		ProductID:       p.ProductID,
		StatusPublikasi: p.StatusPublikasi,
	}
}

func ToProductUpdateResponse(p models.Product) responses.UpdateProductResponse {
	return responses.UpdateProductResponse{
		ProductID:       p.ProductID,
		SellerID:        p.SellerID,
		NamaProduk:      p.NamaProduk,
		BrandID:         p.BrandID,
		CategoryID:      p.CategoryID,
		Kondisi:         strings.ToUpper(p.Kondisi),
		Deskripsi:       p.Deskripsi,
		Harga:           p.Harga,
		Stok:            p.Stok,
		Berat:           p.Berat,
		ConditionScore:  p.ConditionScore,
		StatusPublikasi: displayProductStatus(p.StatusPublikasi),
		UkuranTersedia:  unmarshalUkuran(p.UkuranTersedia),
		UpdatedAt:       p.UpdatedAt,
	}
}

func ToProductImageResponse(pi models.ProductImage) responses.ProductImageResponse {
	resp := responses.ProductImageResponse{
		ImageID:      pi.ImageID,
		ProductID:    pi.ProductID,
		URL:          pi.URLObjectStorage,
		UrutanTampil: pi.UrutanTampil,
		CreatedAt:    pi.CreatedAt,
	}
	if pi.Embedding != nil {
		var vec []float64
		if err := json.Unmarshal(pi.Embedding.Vector, &vec); err == nil && len(vec) > 0 {
			resp.Embedding = vec
		}
	}
	return resp
}

func ToProductImageListResponse(images []models.ProductImage) []responses.ProductImageResponse {
	out := make([]responses.ProductImageResponse, 0, len(images))
	for _, img := range images {
		out = append(out, ToProductImageResponse(img))
	}
	return out
}

func unmarshalUkuran(data []byte) []string {
	var sizes []string
	if len(data) > 0 {
		_ = json.Unmarshal(data, &sizes)
	}
	return sizes
}

func firstImageURL(images []models.ProductImage) string {
	if len(images) == 0 {
		return ""
	}
	return images[0].URLObjectStorage
}

func ToProductListItemResponse(p models.Product, rs services.RatingSummary) responses.ProductListItemResponse {
	item := responses.ProductListItemResponse{
		ProductID:      p.ProductID,
		NamaProduk:     p.NamaProduk,
		Harga:          p.Harga,
		Kondisi:        strings.ToUpper(p.Kondisi),
		Stok:           p.Stok,
		UkuranTersedia: unmarshalUkuran(p.UkuranTersedia),
		ConditionScore: p.ConditionScore,
		AvgRating:      rs.AvgRating,
		TotalReview:    rs.TotalReview,
		ImageURL:       firstImageURL(p.Images),
	}
	if p.Seller.SellerID != "" {
		item.Seller = responses.SellerInfoResponse{
			SellerID:         p.Seller.SellerID,
			NamaToko:         p.Seller.NamaToko,
			SellerTrustScore: p.Seller.SellerTrustScore,
		}
	}
	return item
}

func ToProductListResponse(products []models.Product, summaries map[string]services.RatingSummary) []responses.ProductListItemResponse {
	out := make([]responses.ProductListItemResponse, 0, len(products))
	for _, p := range products {
		out = append(out, ToProductListItemResponse(p, summaries[p.ProductID]))
	}
	return out
}

func ToProductDetailResponse(p models.Product, rs services.RatingSummary) responses.ProductDetailResponse {
	detail := responses.ProductDetailResponse{
		ProductID:       p.ProductID,
		SellerID:        p.SellerID,
		NamaProduk:      p.NamaProduk,
		BrandID:         p.BrandID,
		CategoryID:      p.CategoryID,
		Kondisi:         strings.ToUpper(p.Kondisi),
		Deskripsi:       p.Deskripsi,
		Harga:           p.Harga,
		Stok:            p.Stok,
		UkuranTersedia:  unmarshalUkuran(p.UkuranTersedia),
		ConditionScore:  p.ConditionScore,
		AvgRating:       rs.AvgRating,
		TotalReview:     rs.TotalReview,
		StatusPublikasi: strings.ToUpper(p.StatusPublikasi),
		Images:          make([]responses.ProductDetailImageResponse, 0, len(p.Images)),
	}
	for _, img := range p.Images {
		detail.Images = append(detail.Images, responses.ProductDetailImageResponse{
			ImageID:      img.ImageID,
			URL:          img.URLObjectStorage,
			UrutanTampil: img.UrutanTampil,
		})
	}
	return detail
}
