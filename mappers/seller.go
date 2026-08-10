package mappers

import (
	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/responses"
)

func ToSellerCreate(u *models.Seller) *responses.SellerResponse {
	return &responses.SellerResponse{
		SellerID:         u.SellerID,
		UserID:           u.UserID,
		NamaToko:         u.NamaToko,
		StatusVerifikasi: u.StatusVerifikasi,
	}
}
