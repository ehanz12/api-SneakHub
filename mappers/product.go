package mappers

import (
	"encoding/json"

	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/responses"
)

func ToProductResponse(p models.Product) responses.CreateProductResponse {
	return responses.CreateProductResponse{
		ProductID:       p.ProductID,
		StatusPublikasi: p.StatusPublikasi,
	}
}

func ToProductUpdateResponse(p models.Product) responses.UpdateProductResponse {
	return responses.UpdateProductResponse{
		ProductID: p.ProductID,
		Harga:     p.Harga,
		UpdatedAt: p.UpdatedAt,
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