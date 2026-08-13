package requests

import (
	"strings"
)

type CreateReviewRequest struct {
	ProductID string  `json:"product_id"`
	Rating    float64 `json:"rating"`
	Komentar  *string `json:"komentar"`
}

func (r *CreateReviewRequest) Validate() map[string]string {
	errs := make(map[string]string)

	if strings.TrimSpace(r.ProductID) == "" {
		errs["product_id"] = "product_id wajib diisi"
	}
	if r.Rating < 1 || r.Rating > 5 {
		errs["rating"] = "rating harus antara 1 sampai 5"
	}

	return errs
}
