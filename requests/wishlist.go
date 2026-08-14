package requests

import "strings"

type CreateWishlistRequest struct {
	ProductID string `json:"product_id"`
}

func (r *CreateWishlistRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if strings.TrimSpace(r.ProductID) == "" {
		errs["product_id"] = "product_id wajib diisi"
	}
	return errs
}

type PriceAlertRequest struct {
	TargetPrice *float64 `json:"target_price"`
}

func (r *PriceAlertRequest) Validate() map[string]string {
	errs := make(map[string]string)
	if r.TargetPrice == nil {
		errs["target_price"] = "target_price wajib diisi"
	} else if *r.TargetPrice <= 0 {
		errs["target_price"] = "target_price harus lebih dari 0"
	}
	return errs
}
