package mappers

import (
	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/responses"
	"github.com/ehanz12/api-SneakHub/utils"
)

func ToUserRes(u *models.User) *responses.UserResponse {
	nomor_telp := ""
	if u.NomorTelepon != nil {
		nomor_telp = *u.NomorTelepon
	}
	return &responses.UserResponse{
		UserID:       u.UserID,
		Nama:         u.Nama,
		Email:        u.Email,
		NomorTelepon: &nomor_telp,
		Peran:        u.Peran,
		StatusAkun:   u.StatusAkun,
	}
}

func ToLoginRes(u *models.User) *responses.LoginRes {
	return &responses.LoginRes{
		UserID: u.UserID,
		Nama:   u.Nama,
		Email:  u.Email,
		Peran:  u.Peran,
	}
}

func ToUserBigRes(u *models.User) *responses.UserBigResponse {
	var nomorTelp string
	if u.NomorTelepon != nil {
		nomorTelp = *u.NomorTelepon
	}

	return &responses.UserBigResponse{
		UserID:           u.UserID,
		Nama:             u.Nama,
		Email:            u.Email,
		NomorTelepon:     &nomorTelp,
		Peran:            u.Peran,
		StatusAkun:       u.StatusAkun,
		PreferensiUkuran: utils.MapJSONToStringSlice(u.PreferensiUkuran),
		BrandFavorit:     utils.MapJSONToStringSlice(u.BrandFavorit),
	}
}

func ToUserUpdateRes(u *models.User) *responses.UpdateUserResponse {
	var nomorTelp string
	if u.NomorTelepon != nil {
		nomorTelp = *u.NomorTelepon
	}

	return &responses.UpdateUserResponse{
		UserID:           u.UserID,
		Nama:             u.Nama,
		NomorTelepon:     &nomorTelp,
		PreferensiUkuran: utils.MapJSONToStringSlice(u.PreferensiUkuran),
		BrandFavorit:     utils.MapJSONToStringSlice(u.BrandFavorit),
	}
}
