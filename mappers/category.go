package mappers

import (
	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/responses"
)

func ToCategoryRes(c *models.Category) *responses.CategoryResponse {
	return &responses.CategoryResponse{
		CategoryID:   c.CategoryID,
		NamaKategori: c.NamaKategori,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}

func ToCategoryRes2(c models.Category) responses.CategoryResponse {
	return responses.CategoryResponse{
		CategoryID:   c.CategoryID,
		NamaKategori: c.NamaKategori,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}

func ListToCategoryResponse(cs []models.Category) []responses.CategoryResponse {
	res := make([]responses.CategoryResponse, 0, len(cs))
	for _, c := range cs {
		res = append(res, ToCategoryRes2(c))
	}

	return res
}
