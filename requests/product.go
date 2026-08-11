package requests

import (
	"strings"

	"github.com/ehanz12/api-SneakHub/utils"
)

type CreateProduct struct {
	NamaProduk      string   `json:"nama_produk"`
	BrandID         string   `json:"brand_id"`
	CategoryID      string   `json:"category_id"`
	Kondisi         string   `json:"kondisi"`
	Deskripsi       *string  `json:"deskripsi"`
	Harga           int      `json:"harga"`
	Stok            int      `json:"stok"`
	ConditionScore  *float64 `json:"condition_score"`
	StatusPublikasi string   `json:"status_publikasi"`
	UkuranTersedia  []string `json:"ukuran_tersedia"`
}

func (r *CreateProduct) Validate() map[string]string {
	errs := make(map[string]string)

	if strings.TrimSpace(r.NamaProduk) == "" {
		errs["nama_produk"] = "nama produk wajib diisi"
	}
	if strings.TrimSpace(r.BrandID) == "" {
		errs["brand_id"] = "brand wajib diisi"
	}
	if strings.TrimSpace(r.CategoryID) == "" {
		errs["category_id"] = "category wajib diisi"
	}
	if r.Harga <= 0 {
		errs["harga"] = "harga harus lebih dari 0"
	}
	if r.Deskripsi != nil {
		if len(strings.TrimSpace(*r.Deskripsi)) < 10 {
			errs["deskripsi"] = "deskripsi harus lebih dari 10 karakter"
		}
	}

	if r.Stok < 0 {
		errs["stok"] = "stok tidak boleh kurang dari 0"
	}

	if len(r.UkuranTersedia) == 0 {
		errs["ukuran_tersedia"] = "ukuran tersedia wajib diisi"
	} else {
		for _, ukuran := range r.UkuranTersedia {
			if !utils.IsValidSize(ukuran) {
				errs["ukuran_tersedia"] = "terdapat ukuran sepatu yang tidak valid"
				break
			}
		}
	}
	return errs
}
