package mappers

import (
	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/responses"
)

func ToReviewResponse(r models.Review) responses.ReviewResponse {
	return responses.ReviewResponse{
		ReviewID:  r.ReviewID,
		OrderID:   r.OrderID,
		ProductID: r.ProductID,
		Rating:    r.Rating,
		Komentar:  r.Komentar,
	}
}
