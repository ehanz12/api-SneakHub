package services

import (
	"errors"

	"github.com/ehanz12/api-SneakHub/database"
	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/requests"
)

func CreateCategoryService(r requests.CategoryRequest) (*models.Category, error) {
	tx := database.DB.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return nil, errors.New("gagal menyambungan ke server")
	}
	c := models.Category{
		NamaKategori: r.NamaKategori,
	}
	if err := tx.Create(&c).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal membuat category")
	}
	tx.Commit()
	return &c, nil
}

func GetCategoryService() ([]models.Category, error) {
	var category []models.Category
	if err := database.DB.Find(&category).Error; err != nil {
		return nil, errors.New("gagal memuat category")
	}
	return category, nil
}

func UpdateCategoryService(cID string, r requests.CategoryRequest) (*models.Category, error) {
	var category models.Category
	if err := database.DB.Where("category_id = ?", cID).First(&category).Error; err != nil {
		return nil, errors.New("category tidak ditemukan")
	}
	if r.NamaKategori != "" {
		category.NamaKategori = r.NamaKategori
	}
	if err := database.DB.Save(&category).Error; err != nil {
		return nil, errors.New("gagal update category")
	}
	return &category, nil
}

func DeleteCategoryService(cID string) error {
	var c models.Category
	if err := database.DB.Where("category_id = ? ", cID).First(&c).Error; err != nil {
		return errors.New("cateogry tidak ditemukan")
	}
	if err := database.DB.Delete(&c).Error; err != nil {
		return errors.New("gagal menghapus category")
	}
	return nil
}
