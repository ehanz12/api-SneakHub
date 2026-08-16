package requests

import (
	"strings"
)

type ShippingRatesRequest struct {
	AddressID string `json:"address_id"`
}

func (r *ShippingRatesRequest) Validate() map[string]string {
	errs := make(map[string]string)

	if strings.TrimSpace(r.AddressID) == "" {
		errs["address_id"] = "address_id wajib diisi"
	}
	return errs
}

type CheckoutShippingRequest struct {
	SellerID string `json:"seller_id"`
	Kurir    string `json:"kurir"`
}

type CheckoutRequest struct {
	AddressID        string                    `json:"address_id"`
	MetodePembayaran string                    `json:"metode_pembayaran"`
	Pengiriman       []CheckoutShippingRequest `json:"pengiriman"`
}

func (r *CheckoutRequest) Validate() map[string]string {
	errs := make(map[string]string)

	if strings.TrimSpace(r.AddressID) == "" {
		errs["address_id"] = "address_id wajib diisi"
	}
	if strings.TrimSpace(r.MetodePembayaran) == "" {
		errs["metode_pembayaran"] = "metode_pembayaran wajib diisi"
	}
	return errs
}
