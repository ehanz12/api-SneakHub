package services

import (
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
