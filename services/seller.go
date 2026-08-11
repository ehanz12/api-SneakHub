package services

import (
	"errors"

	"github.com/ehanz12/api-SneakHub/database"
	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/requests"
)

func CreateSellerService(UserID string, req requests.CreateSellerRequest) (*models.Seller, error) {
	tx := database.DB.Begin()
	if tx.Error != nil {
		return nil, errors.New("gagal menyambungkan ke server")
	}
	var user models.User
	err := tx.Select("user_id").Where("user_id = ?", UserID).First(&user).Error
	if err != nil {
		tx.Rollback()
		return nil, errors.New("user tidak ditemukan")
	}
	var exist models.Seller
	if err := tx.Select("user_id").Where("user_id = ?", UserID).First(&exist).Error; err == nil {
		return nil, errors.New("user sudah mengajukan menjadi seller")
	}
	seller := models.Seller{
		UserID:        UserID,
		NamaToko:      req.NamaToko,
		DeskripsiToko: req.DeskripsiToko,
	}

	if err := tx.Create(&seller).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal membuat seller")
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New("terjadi kesalahan server")
	}
	return &seller, nil
}
