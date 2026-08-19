package mappers

import (
	"strings"

	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/responses"
	"github.com/ehanz12/api-SneakHub/services"
)

func displayProductStatus(status string) string {
	switch status {
	case "aktif":
		return "ACTIVE"
	case "draft":
		return "DRAFT"
	case "nonaktif":
		return "INACTIVE"
	case "pending":
		return "PENDING"
	}
	return strings.ToUpper(status)
}

func displayUserStatus(status string) string {
	switch status {
	case "aktif":
		return "ACTIVE"
	case "tidak_aktif":
		return "INACTIVE"
	case "blokir":
		return "SUSPENDED"
	}
	return strings.ToUpper(status)
}

func displayOrderStatus(status string) string {
	switch status {
	case "pending":
		return "PENDING"
	case "diproses":
		return "PROCESSING"
	case "dikirim":
		return "SHIPPED"
	case "selesai":
		return "COMPLETED"
	case "dibatalkan":
		return "CANCELLED"
	}
	return strings.ToUpper(status)
}

func displaySellerStatus(status string) string {
	switch status {
	case "pending":
		return "PENDING"
	case "verified":
		return "VERIFIED"
	case "rejected":
		return "REJECTED"
	}
	return strings.ToUpper(status)
}

func ToSellerCreate(u *models.Seller) *responses.SellerResponse {
	return &responses.SellerResponse{
		SellerID:         u.SellerID,
		UserID:           u.UserID,
		NamaToko:         u.NamaToko,
		DeskripsiToko:    u.DeskripsiToko,
		StatusVerifikasi: displaySellerStatus(u.StatusVerifikasi),
		KodePosAsal:      u.KodePosAsal,
		KotaAsal:         u.KotaAsal,
		AlamatAsal:       u.AlamatAsal,
	}
}

func ToSellerProductListResponse(products []services.ProductWithSales) []responses.SellerProductListItemResponse {
	out := make([]responses.SellerProductListItemResponse, 0, len(products))
	for _, p := range products {
		out = append(out, responses.SellerProductListItemResponse{
			ProductID:       p.ProductID,
			NamaProduk:      p.NamaProduk,
			Harga:           p.Harga,
			ImageURL:        p.ImageURL,
			Stok:            p.Stok,
			StatusPublikasi: displayProductStatus(p.StatusPublikasi),
			TotalTerjual:    p.TotalTerjual,
		})
	}
	return out
}

func ToSellerOrderListResponse(orders []models.Order) []responses.SellerOrderListItemResponse {
	out := make([]responses.SellerOrderListItemResponse, 0, len(orders))
	for _, o := range orders {
		item := responses.SellerOrderListItemResponse{
			OrderID:         o.OrderID,
			StatusOrder:     displayOrderStatus(o.StatusOrder),
			TotalPembayaran: o.TotalPesanan,
			CreatedAt:       o.CreatedAt,
		}
		if o.Customer.UserID != "" {
			item.Customer = responses.SellerCustomerResponse{
				UserID: o.Customer.UserID,
				Nama:   o.Customer.Nama,
			}
		}
		out = append(out, item)
	}
	return out
}

func ToSellerDashboardResponse(d *services.SellerDashboardData) responses.SellerDashboardResponse {
	top := make([]responses.SellerTopProductResponse, 0, len(d.ProdukTerlaris))
	for _, p := range d.ProdukTerlaris {
		top = append(top, responses.SellerTopProductResponse{
			ProductID:    p.ProductID,
			NamaProduk:   p.NamaProduk,
			ImageURL:     p.ImageURL,
			TotalTerjual: p.TotalTerjual,
		})
	}
	return responses.SellerDashboardResponse{
		TotalProduk:      d.TotalProduk,
		ProdukAktif:      d.ProdukAktif,
		TotalTerjual:     d.TotalTerjual,
		TotalPendapatan:  d.TotalPendapatan,
		RatingRataRata:   d.RatingRataRata,
		SellerTrustScore: d.SellerTrustScore,
		ProdukTerlaris:   top,
	}
}
