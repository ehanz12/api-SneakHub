package requests

import (
	"regexp"
	"strings"
)

var kodePosPattern = regexp.MustCompile(`^\d{5}$`)

type CreateSellerRequest struct {
	NamaToko      string  `json:"nama_toko"`
	DeskripsiToko *string `json:"deskripsi_toko"`
	KodePosAsal   *string `json:"kode_pos_asal"`
	KotaAsal      *string `json:"kota_asal"`
	AlamatAsal    *string `json:"alamat_asal"`
}

func (r *CreateSellerRequest) Validation() map[string]string {
	errs := make(map[string]string)

	if strings.TrimSpace(r.NamaToko) == "" {
		errs["nama_toko"] = "nama toko wajib diisi"
	} else if len(strings.TrimSpace(r.NamaToko)) < 3 {
		errs["nama_toko"] = "nama toko minimal 3 karakter"
	}

	if r.DeskripsiToko != nil {
		if len(strings.TrimSpace(*r.DeskripsiToko)) < 3 {
			errs["deskripsi_toko"] = "deskripsi toko minimal 3 karakter"
		}
	}

	validateOriginFields(r.KodePosAsal, r.KotaAsal, r.AlamatAsal, errs)

	return errs
}

// UpdateSellerProfileRequest adalah body update profil toko (partial update).
// Field yang tidak dikirim tidak berubah; field dikirim kosong ("")
// akan dihapus (NULL).
type UpdateSellerProfileRequest struct {
	NamaToko      *string `json:"nama_toko"`
	DeskripsiToko *string `json:"deskripsi_toko"`
	KodePosAsal   *string `json:"kode_pos_asal"`
	KotaAsal      *string `json:"kota_asal"`
	AlamatAsal    *string `json:"alamat_asal"`
}

func (r *UpdateSellerProfileRequest) Validation() map[string]string {
	errs := make(map[string]string)

	if r.NamaToko != nil {
		if strings.TrimSpace(*r.NamaToko) == "" {
			errs["nama_toko"] = "nama toko wajib diisi"
		} else if len(strings.TrimSpace(*r.NamaToko)) < 3 {
			errs["nama_toko"] = "nama toko minimal 3 karakter"
		}
	}

	if r.DeskripsiToko != nil {
		if len(strings.TrimSpace(*r.DeskripsiToko)) < 3 {
			errs["deskripsi_toko"] = "deskripsi toko minimal 3 karakter"
		}
	}

	validateOriginFields(r.KodePosAsal, r.KotaAsal, r.AlamatAsal, errs)

	return errs
}

// validateOriginFields memvalidasi kolom asal pengiriman toko (opsional).
// Kode pos wajib 5 digit; kota & alamat minimal panjang tertentu jika diisi.
func validateOriginFields(kodePos, kota, alamat *string, errs map[string]string) {
	if kodePos != nil {
		trimmed := strings.TrimSpace(*kodePos)
		if trimmed != "" && !kodePosPattern.MatchString(trimmed) {
			errs["kode_pos_asal"] = "kode pos asal harus 5 digit angka"
		}
	}

	if kota != nil {
		trimmed := strings.TrimSpace(*kota)
		if trimmed != "" && len(trimmed) < 2 {
			errs["kota_asal"] = "kota asal minimal 2 karakter"
		}
	}

	if alamat != nil {
		trimmed := strings.TrimSpace(*alamat)
		if trimmed != "" && len(trimmed) < 5 {
			errs["alamat_asal"] = "alamat asal minimal 5 karakter"
		}
	}
}