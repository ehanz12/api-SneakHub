package mappers

import (
	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/responses"
)

func subtotal(item models.CartItem) float64 {
	return item.HargaSaatDitambahkan * float64(item.Jumlah)
}

func ToCartItemResponse(ci models.CartItem) responses.CartItemResponse {
	return responses.CartItemResponse{
		CartItemID: ci.CartItemID,
		ProductID:  ci.ProductID,
		VariantID:  ci.VariantID,
		Jumlah:     ci.Jumlah,
		Subtotal:   subtotal(ci),
	}
}

func ToCartItemListResponse(items []models.CartItem) []responses.CartItemResponse {
	out := make([]responses.CartItemResponse, 0, len(items))
	for _, ci := range items {
		out = append(out, ToCartItemResponse(ci))
	}
	return out
}

func ToCartItemDetailResponse(ci models.CartItem) responses.CartItemDetailResponse {
	resp := responses.CartItemDetailResponse{
		CartItemResponse: ToCartItemResponse(ci),
		NamaProduk:       ci.Product.NamaProduk,
		Harga:            ci.Product.Harga,
	}
	if len(ci.Product.Images) > 0 {
		resp.ImageURL = ci.Product.Images[0].URLObjectStorage
	}
	return resp
}

func ToCartResponse(cart models.Cart) responses.CartResponse {
	resp := responses.CartResponse{
		CartID:    cart.CartID,
		Items:     make([]responses.CartItemDetailResponse, 0, len(cart.Items)),
		TotalItem: 0,
		Total:     0,
	}
	for _, ci := range cart.Items {
		resp.Total += subtotal(ci)
		resp.TotalItem += ci.Jumlah
		resp.Items = append(resp.Items, ToCartItemDetailResponse(ci))
	}
	return resp
}
