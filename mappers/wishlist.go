package mappers

import (
	"strings"

	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/responses"
	"github.com/ehanz12/api-SneakHub/services"
)

func displayStokStatus(status string) string {
	switch status {
	case "available":
		return "AVAILABLE"
	case "out_of_stock":
		return "OUT_OF_STOCK"
	}
	return strings.ToUpper(status)
}

func ToWishlistListResponse(wishlists []models.Wishlist) []responses.WishlistItemResponse {
	out := make([]responses.WishlistItemResponse, 0, len(wishlists))
	for _, w := range wishlists {
		out = append(out, responses.WishlistItemResponse{
			WishlistID:         w.WishlistID,
			ProductID:          w.ProductID,
			NamaProduk:         w.Product.NamaProduk,
			Harga:              w.Product.Harga,
			ImageURL:           services.FirstImageURL(w.Product.Images),
			StatusStokTerakhir: displayStokStatus(w.StatusStok),
			PriceAlert: responses.WishlistPriceAlertResponse{
				Enabled:     w.PriceAlertEnabled,
				TargetPrice: w.TargetPrice,
			},
			RestockAlert: responses.WishlistRestockAlertResponse{
				Enabled: w.RestockAlertEnabled,
			},
		})
	}
	return out
}

func ToCreateWishlistResponse(w models.Wishlist) responses.CreateWishlistResponse {
	return responses.CreateWishlistResponse{
		WishlistID: w.WishlistID,
		ProductID:  w.ProductID,
	}
}

func ToPriceAlertResponse(w models.Wishlist) responses.PriceAlertResponse {
	return responses.PriceAlertResponse{
		ProductID:   w.ProductID,
		Enabled:     w.PriceAlertEnabled,
		TargetPrice: w.TargetPrice,
	}
}

func ToRestockAlertResponse(w models.Wishlist) responses.RestockAlertResponse {
	return responses.RestockAlertResponse{
		ProductID: w.ProductID,
		Enabled:   w.RestockAlertEnabled,
	}
}
