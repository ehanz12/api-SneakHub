package requests

import (
	"fmt"
	"strings"
)

type CartItemRequest struct {
	ProductID string `json:"product_id"`
	VariantID string `json:"variant_id"`
	Jumlah    int    `json:"jumlah"`
}

type AddCartItemsRequest struct {
	Items []CartItemRequest `json:"items"`
}

func (r *AddCartItemsRequest) Validate() map[string]string {
	errs := make(map[string]string)

	if len(r.Items) == 0 {
		errs["items"] = "minimal satu item wajib diisi"
		return errs
	}

	for i, item := range r.Items {
		if strings.TrimSpace(item.ProductID) == "" {
			errs[fmt.Sprintf("items[%d].product_id", i)] = "product_id wajib diisi"
		}
		if item.Jumlah < 1 {
			errs[fmt.Sprintf("items[%d].jumlah", i)] = "jumlah minimal 1"
		}
	}

	return errs
}

type UpdateCartItemRequest struct {
	Jumlah int `json:"jumlah"`
}

func (r *UpdateCartItemRequest) Validate() map[string]string {
	errs := make(map[string]string)

	if r.Jumlah < 1 {
		errs["jumlah"] = "jumlah minimal 1"
	}

	return errs
}
