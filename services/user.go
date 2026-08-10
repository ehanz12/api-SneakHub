package services

import (
	"encoding/json"
	"errors"

	"github.com/ehanz12/api-SneakHub/database"
	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/requests"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUserService(req requests.RegisterRequest) (*models.User, error) {
	var exits models.User
	tx := database.DB.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return &models.User{}, errors.New("keasalahan server")
	}
	err := tx.Select("user_id", "email").Where("email = ?", req.Email).First(&exits).Error
	if err == nil {
		tx.Rollback()
		return &models.User{}, errors.New("email sudah digunakan")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return nil, errors.New("keasalahan server")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return &models.User{}, errors.New("password gagal dibuat")
	}
	user := models.User{
		Nama:          req.Nama,
		Email:         req.Email,
		NomorTelepon:  req.NomorTelepon,
		KataSandiHash: string(hash),
	}
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return &models.User{}, errors.New("gagal register")
	}

	tx.Commit()
	return &user, nil
}

func LoginUserService(req requests.LoginRequest) (*models.User, error) {
	var exits models.User
	err := database.DB.Select("user_id", "nama", "email", "peran", "kata_sandi_hash", "peran").Where("email = ?", req.Email).First(&exits).Error
	if err != nil {
		return nil, errors.New("email tidak ditemukan")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(exits.KataSandiHash), []byte(req.Password)); err != nil {
		return nil, errors.New("kesalahan password")
	}
	return &exits, nil
}

func MeUserService(userID string) (*models.User, error) {
	var user models.User
	err := database.DB.Select("user_id", "nama", "email", "nomor_telepon", "peran", "status_akun",
		"preferensi_ukuran", "brand_favorit").Where("user_id = ?", userID).Find(&user).Error
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}
	return &user, nil
}

func UpdateUserService(userID string, req requests.UpdateUserRequest) (*models.User, error) {
	tx := database.DB.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return nil, errors.New("terjadi kesalahan server")
	}
	var user models.User
	if err := tx.Where("user_id = ?", userID).First(&user).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("user gagal ditemukan")
	}
	updates := make(map[string]interface{})

	if req.Nama != "" {
		updates["nama"] = req.Nama
	}

	if req.NomorTelepon != nil {
		updates["nomor_telepon"] = req.NomorTelepon
	}

	if req.PreferensiUkuran != nil {
		data, err := json.Marshal(req.PreferensiUkuran)
		if err != nil {
			tx.Rollback()
			return nil, errors.New("preferensi ukuran tidak valid")
		}

		updates["preferensi_ukuran"] = data
	}

	if req.BrandFavorit != nil {
		data, err := json.Marshal(req.BrandFavorit)
		if err != nil {
			tx.Rollback()
			return nil, errors.New("brand favorit tidak valid")
		}

		updates["brand_favorit"] = data
	}

	if len(updates) == 0 {
		tx.Rollback()
		return nil, errors.New("tidak ada data yang diubah")
	}

	if err := tx.
		Model(&models.User{}).
		Where("user_id = ?", userID).
		Updates(updates).Error; err != nil {

		tx.Rollback()
		return nil, errors.New("gagal update user")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("gagal commit")
	}

	// Ambil data terbaru setelah update.
	if err := database.DB.Select("user_id", "nomor_telepon", "nama", "preferensi_ukuran", "brand_favorit").
		Where("user_id = ?", userID).
		First(&user).Error; err != nil {
		return nil, errors.New("gagal mengambil data user terbaru")
	}

	return &user, nil
}
