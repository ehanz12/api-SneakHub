package requests

import (
	"strings"

	"github.com/ehanz12/api-SneakHub/utils"
)

type RegisterRequest struct {
	Nama         string  `json:"nama"`
	Email        string  `json:"email"`
	NomorTelepon *string `json:"nomor_telepon"`
	Password     string  `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Nama             string   `json:"nama"`
	NomorTelepon     *string  `json:"nomor_telepon"`
	PreferensiUkuran []string `json:"preferensi_ukuran"`
	BrandFavorit     []string `json:"brand_favorit"`
}

func (r *RegisterRequest) Validate() map[string]string {
	errs := make(map[string]string)

	if strings.TrimSpace(r.Nama) == "" {
		errs["name"] = "nama harus diisi"
	}
	if strings.TrimSpace(r.Email) == "" || !strings.Contains(r.Email, "@") {
		errs["email"] = "email harus diisi dan formatnya benar"
	}
	if len(r.Password) < 6 {
		errs["password"] = "password harus lebih 6 karakter"
	}

	if r.NomorTelepon != nil {
		nomor := strings.TrimSpace(*r.NomorTelepon)

		if nomor != "" {
			if ok := utils.IsValidPhone(nomor); !ok {
				errs["nomor_telepon"] = "format hp tidak valid"
			}
		}
	}

	return errs
}

func (r *LoginRequest) ValidateLogin() map[string]string {
	errs := make(map[string]string)

	if strings.TrimSpace(r.Email) == "" || !strings.Contains(r.Email, "@") {
		errs["email"] = "email harus diisi dan formatnya benar"
	}
	if len(r.Password) < 6 {
		errs["password"] = "password harus lebih 6 karakter"
	}
	return errs
}

func (r *UpdateUserRequest) ValidateUpdateUser() map[string]string {
	errs := make(map[string]string)

	if r.NomorTelepon != nil {
		nomor := strings.TrimSpace(*r.NomorTelepon)

		if nomor != "" {
			if ok := utils.IsValidPhone(nomor); !ok {
				errs["nomor_telepon"] = "format hp tidak valid"
			}
		}
	}
	return errs
}
