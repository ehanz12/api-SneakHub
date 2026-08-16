package mappers

import (
	"github.com/ehanz12/api-SneakHub/responses"
	"github.com/ehanz12/api-SneakHub/services"
)

func ToShippingRatesResponse(rates []services.SellerShippingRates) []responses.SellerShippingRatesResponse {
	out := make([]responses.SellerShippingRatesResponse, 0, len(rates))
	for _, r := range rates {
		options := make([]responses.ShippingOptionResponse, 0, len(r.Options))
		for _, opt := range r.Options {
			options = append(options, responses.ShippingOptionResponse{
				Kurir:      opt.Kurir,
				Service:    opt.Service,
				Biaya:      opt.Biaya,
				Estimasi:   opt.Estimasi,
				IsFallback: opt.IsFallback,
			})
		}
		out = append(out, responses.SellerShippingRatesResponse{
			SellerID: r.SellerID,
			NamaToko: r.NamaToko,
			Berat:    r.Berat,
			KotaAsal: r.KotaAsal,
			Options:  options,
		})
	}
	return out
}
