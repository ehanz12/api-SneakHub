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

func ToReviewListResponse(reviews []models.Review) []responses.ReviewListItemResponse {
	out := make([]responses.ReviewListItemResponse, 0, len(reviews))
	for _, r := range reviews {
		item := responses.ReviewListItemResponse{
			ReviewID:  r.ReviewID,
			ProductID: r.ProductID,
			Rating:    r.Rating,
			Komentar:  r.Komentar,
			CreatedAt: r.CreatedAt,
		}
		if r.Customer.UserID != "" {
			item.Customer = responses.SellerCustomerResponse{
				UserID: r.Customer.UserID,
				Nama:   r.Customer.Nama,
			}
		}
		out = append(out, item)
	}
	return out
}
