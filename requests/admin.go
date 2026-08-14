package requests

import (
	"strings"
)

type UpdateUserStatusRequest struct {
	StatusAkun string `json:"status_akun"`
	Alasan     string `json:"alasan"`
}

func (r *UpdateUserStatusRequest) Validate() map[string]string {
	errs := make(map[string]string)

	status := strings.ToUpper(strings.TrimSpace(r.StatusAkun))
	if status == "" {
		errs["status_akun"] = "status_akun wajib diisi"
	} else {
		switch status {
		case "ACTIVE", "INACTIVE", "SUSPENDED", "BLOCKED":
		default:
			errs["status_akun"] = "status_akun harus salah satu dari: ACTIVE, INACTIVE, SUSPENDED"
		}
	}

	if strings.TrimSpace(r.Alasan) == "" {
		errs["alasan"] = "alasan wajib diisi"
	}

	return errs
}

type UpdateProductStatusRequest struct {
	StatusPublikasi string `json:"status_publikasi"`
	Catatan         string `json:"catatan"`
}

func (r *UpdateProductStatusRequest) Validate() map[string]string {
	errs := make(map[string]string)

	status := strings.ToUpper(strings.TrimSpace(r.StatusPublikasi))
	if status == "" {
		errs["status_publikasi"] = "status_publikasi wajib diisi"
	} else {
		switch status {
		case "ACTIVE", "DRAFT", "INACTIVE", "PENDING":
		default:
			errs["status_publikasi"] = "status_publikasi harus salah satu dari: ACTIVE, DRAFT, INACTIVE, PENDING"
		}
	}

	if strings.TrimSpace(r.Catatan) == "" {
		errs["catatan"] = "catatan wajib diisi"
	}

	return errs
}
