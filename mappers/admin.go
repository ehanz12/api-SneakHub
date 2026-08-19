package mappers

import (
	"strings"

	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/responses"
	"github.com/ehanz12/api-SneakHub/services"
)

func ToAdminUserListResponse(users []models.User) []responses.AdminUserListItemResponse {
	out := make([]responses.AdminUserListItemResponse, 0, len(users))
	for _, u := range users {
		out = append(out, responses.AdminUserListItemResponse{
			UserID:     u.UserID,
			Nama:       u.Nama,
			Email:      u.Email,
			Peran:      strings.ToUpper(u.Peran),
			StatusAkun: displayUserStatus(u.StatusAkun),
		})
	}
	return out
}

func ToAdminUserStatusResponse(u *models.User) responses.AdminUserStatusResponse {
	return responses.AdminUserStatusResponse{
		UserID:     u.UserID,
		StatusAkun: displayUserStatus(u.StatusAkun),
	}
}

func ToAdminSellerListResponse(sellers []models.Seller) []responses.AdminSellerListItemResponse {
	out := make([]responses.AdminSellerListItemResponse, 0, len(sellers))
	for _, s := range sellers {
		item := responses.AdminSellerListItemResponse{
			SellerID:         s.SellerID,
			UserID:           s.UserID,
			NamaToko:         s.NamaToko,
			StatusVerifikasi: displaySellerStatus(s.StatusVerifikasi),
			CreatedAt:        s.CreatedAt,
		}
		if s.User.UserID != "" {
			item.NamaUser = s.User.Nama
			item.Email = s.User.Email
		}
		out = append(out, item)
	}
	return out
}

func ToAdminSellerVerificationResponse(s *models.Seller) responses.AdminSellerVerificationResponse {
	return responses.AdminSellerVerificationResponse{
		SellerID:         s.SellerID,
		StatusVerifikasi: displaySellerStatus(s.StatusVerifikasi),
	}
}

func ToAdminUserRoleResponse(u *models.User) responses.AdminUserRoleResponse {
	return responses.AdminUserRoleResponse{
		UserID: u.UserID,
		Peran:  strings.ToUpper(u.Peran),
	}
}

func ToAdminProductListResponse(products []models.Product) []responses.AdminProductListItemResponse {
	out := make([]responses.AdminProductListItemResponse, 0, len(products))
	for _, p := range products {
		out = append(out, responses.AdminProductListItemResponse{
			ProductID:       p.ProductID,
			NamaProduk:      p.NamaProduk,
			SellerID:        p.SellerID,
			Harga:           p.Harga,
			ImageURL:        services.FirstImageURL(p.Images),
			StatusPublikasi: displayProductStatus(p.StatusPublikasi),
		})
	}
	return out
}

func ToAdminProductStatusResponse(p *models.Product) responses.AdminProductStatusResponse {
	return responses.AdminProductStatusResponse{
		ProductID:       p.ProductID,
		StatusPublikasi: displayProductStatus(p.StatusPublikasi),
	}
}

func ToAdminOrderListResponse(orders []models.Order) []responses.AdminOrderListItemResponse {
	out := make([]responses.AdminOrderListItemResponse, 0, len(orders))
	for _, o := range orders {
		out = append(out, responses.AdminOrderListItemResponse{
			OrderID:         o.OrderID,
			CustomerID:      o.CustomerID,
			SellerID:        o.SellerID,
			StatusOrder:     displayOrderStatus(o.StatusOrder),
			TotalPembayaran: o.TotalPesanan,
			CreatedAt:       o.CreatedAt,
		})
	}
	return out
}

func ToAdminReportResponse(r *services.AdminReportData) responses.AdminReportResponse {
	return responses.AdminReportResponse{
		Period:        r.Period,
		TotalUsers:    r.TotalUsers,
		TotalSellers:  r.TotalSellers,
		TotalProducts: r.TotalProducts,
		TotalOrders:   r.TotalOrders,
		TotalRevenue:  r.TotalRevenue,
	}
}
