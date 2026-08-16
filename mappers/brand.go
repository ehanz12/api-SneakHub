package mappers

import (
	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/responses"
)

func ToBrandResponse(b *models.Brand) *responses.BrandResponse {
	if b == nil {
		return nil
	}
	return &responses.BrandResponse{
		BrandID:   b.BrandID,
		NamaBrand: b.NamaBrand,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}

func ListToBrandResponse(bs []models.Brand) []responses.BrandResponse {
	res := make([]responses.BrandResponse, 0, len(bs))
	for _, b := range bs {
		res = append(res, *ToBrandResponse(&b))
	}
	return res
}
