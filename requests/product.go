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
	Berat           int      `json:"berat"`
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

	if r.Berat < 0 {
		errs["berat"] = "berat tidak boleh kurang dari 0"
	}

	if strings.TrimSpace(r.StatusPublikasi) != "" {
		switch strings.ToUpper(strings.TrimSpace(r.StatusPublikasi)) {
		case "ACTIVE", "AKTIF", "DRAFT", "INACTIVE", "NONAKTIF", "TIDAK_AKTIF":
		default:
			errs["status_publikasi"] = "status_publikasi harus salah satu dari: ACTIVE, DRAFT, INACTIVE"
		}
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

// UpdateProductRequest untuk update parsial: hanya field yang dikirim
// (tidak nil) yang akan diubah. Deskripsi bisa dikosongkan dengan "".
type UpdateProductRequest struct {
	NamaProduk      *string   `json:"nama_produk"`
	BrandID         *string   `json:"brand_id"`
	CategoryID      *string   `json:"category_id"`
	Kondisi         *string   `json:"kondisi"`
	Deskripsi       *string   `json:"deskripsi"`
	Harga           *float64  `json:"harga"`
	Stok            *int      `json:"stok"`
	Berat           *int      `json:"berat"`
	ConditionScore  *float64  `json:"condition_score"`
	StatusPublikasi *string   `json:"status_publikasi"`
	UkuranTersedia  *[]string `json:"ukuran_tersedia"`
}

func (r *UpdateProductRequest) Validate() map[string]string {
	errs := make(map[string]string)

	if r.NamaProduk != nil && strings.TrimSpace(*r.NamaProduk) == "" {
		errs["nama_produk"] = "nama produk wajib diisi"
	}
	if r.Harga != nil && *r.Harga <= 0 {
		errs["harga"] = "harga harus lebih dari 0"
	}
	if r.Deskripsi != nil && len(strings.TrimSpace(*r.Deskripsi)) < 10 {
		errs["deskripsi"] = "deskripsi harus lebih dari 10 karakter"
	}
	if r.Stok != nil && *r.Stok < 0 {
		errs["stok"] = "stok tidak boleh kurang dari 0"
	}
	if r.Berat != nil && *r.Berat < 0 {
		errs["berat"] = "berat tidak boleh kurang dari 0"
	}
	if r.ConditionScore != nil && (*r.ConditionScore < 0 || *r.ConditionScore > 100) {
		errs["condition_score"] = "condition_score harus antara 0 sampai 100"
	}
	if r.StatusPublikasi != nil {
		switch strings.ToUpper(strings.TrimSpace(*r.StatusPublikasi)) {
		case "ACTIVE", "AKTIF", "DRAFT", "INACTIVE", "NONAKTIF", "TIDAK_AKTIF":
		default:
			errs["status_publikasi"] = "status_publikasi harus salah satu dari: ACTIVE, DRAFT, INACTIVE"
		}
	}
	if r.UkuranTersedia != nil {
		if len(*r.UkuranTersedia) == 0 {
			errs["ukuran_tersedia"] = "ukuran tersedia wajib diisi"
		} else {
			for _, ukuran := range *r.UkuranTersedia {
				if !utils.IsValidSize(ukuran) {
					errs["ukuran_tersedia"] = "terdapat ukuran sepatu yang tidak valid"
					break
				}
			}
		}
	}
	return errs
}
