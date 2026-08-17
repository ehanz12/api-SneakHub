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
	statusPembayaran := ""
	if o.Payment != nil {
		statusPembayaran = o.Payment.StatusPembayaran
	}
	return responses.OrderListItemResponse{
		OrderID:          o.OrderID,
		SellerID:         o.SellerID,
		StatusOrder:      o.StatusOrder,
		TotalPembayaran:  o.TotalPesanan,
		StatusPembayaran: statusPembayaran,
		CreatedAt:        o.CreatedAt,
	}
}

func ToOrderListResponse(orders []models.Order) []responses.OrderListItemResponse {
	out := make([]responses.OrderListItemResponse, 0, len(orders))
	for _, o := range orders {
		out = append(out, ToOrderListItemResponse(o))
	}
	return out
}

func ToOrderPaymentResponse(p *models.Payment) *responses.OrderPaymentResponse {
	if p == nil || p.PaymentID == "" {
		return nil
	}
	return &responses.OrderPaymentResponse{
		PaymentID:            p.PaymentID,
		MetodePembayaran:     p.MetodePembayaran,
		Jumlah:               p.Jumlah,
		StatusPembayaran:     p.StatusPembayaran,
		PaymentURL:           p.PaymentURL,
		GatewayReference:     p.GatewayReference,
		TransactionReference: p.TransactionReference,
		PaidAt:               p.PaidAt,
	}
}

func ToOrderShipmentResponse(s *models.Shipment) *responses.OrderShipmentResponse {
	if s == nil || s.ShipmentID == "" {
		return nil
	}
	return &responses.OrderShipmentResponse{
		ShipmentID:       s.ShipmentID,
		Kurir:            s.Kurir,
		NomorResi:        s.NomorResi,
		StatusPengiriman: s.StatusPengiriman,
		ShippedAt:        s.ShippedAt,
		DeliveredAt:      s.DeliveredAt,
	}
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
		Payment:          ToOrderPaymentResponse(o.Payment),
		Shipment:         ToOrderShipmentResponse(o.Shipment),
		CreatedAt:        o.CreatedAt,
	}
}
