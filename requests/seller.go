package requests

import "strings"

type CreateSellerRequest struct {
	NamaToko      string  `json:"nama_toko"`
	DeskripsiToko *string `json:"deskripsi_toko"`
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

	return errs
}
