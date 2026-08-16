package requests

import "strings"

type BrandRequest struct {
	NamaBrand string `json:"nama_brand"`
}

func (r *BrandRequest) Validate() map[string]string {
	errs := make(map[string]string)

	name := strings.TrimSpace(r.NamaBrand)
	if name == "" {
		errs["nama_brand"] = "nama brand wajib diisi"
	} else if len(name) < 3 {
		errs["nama_brand"] = "nama brand minimal 3 karakter"
	}

	return errs
}
