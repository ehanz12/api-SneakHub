package mappers

import (
	"encoding/json"

	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/responses"
)

func ToOrderAlamatResponse(raw []byte) responses.OrderAlamatResponse {
	var alamat responses.OrderAlamatResponse
	_ = json.Unmarshal(raw, &alamat)
	return alamat
}

func ToOrderListItemResponse(o models.Order) responses.OrderListItemResponse {
	return responses.OrderListItemResponse{
		OrderID:         o.OrderID,
		SellerID:        o.SellerID,
		StatusOrder:     o.StatusOrder,
		TotalPembayaran: o.TotalPesanan,
		CreatedAt:       o.CreatedAt,
	}
}

func ToOrderListResponse(orders []models.Order) []responses.OrderListItemResponse {
	out := make([]responses.OrderListItemResponse, 0, len(orders))
	for _, o := range orders {
		out = append(out, ToOrderListItemResponse(o))
	}
	return out
}

func ToOrderDetailResponse(o models.Order) responses.OrderDetailResponse {
	items := make([]responses.OrderItemResponse, 0, len(o.Items))
	for _, item := range o.Items {
		items = append(items, responses.OrderItemResponse{
			OrderItemID:        item.OrderItemID,
			ProductID:          item.ProductID,
			NamaProduk:         item.Product.NamaProduk,
			Jumlah:             item.Jumlah,
			HargaSaatTransaksi: item.HargaSaatTransaksi,
		})
	}

	return responses.OrderDetailResponse{
		OrderID:          o.OrderID,
		CustomerID:       o.CustomerID,
		SellerID:         o.SellerID,
		StatusOrder:      o.StatusOrder,
		AlamatPengiriman: ToOrderAlamatResponse(o.AlamatPengiriman),
		MetodePembayaran: o.MetodePembayaran,
		Items:            items,
		Subtotal:         o.Subtotal,
		BiayaPengiriman:  o.BiayaPengiriman,
		TotalPembayaran:  o.TotalPesanan,
		CreatedAt:        o.CreatedAt,
	}
}
