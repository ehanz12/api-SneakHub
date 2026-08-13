package requests

import (
	"strings"
)

type CheckoutRequest struct {
	AddressID        string `json:"address_id"`
	MetodePembayaran string `json:"metode_pembayaran"`
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
