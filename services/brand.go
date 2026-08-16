package services

import (
	"errors"
	"strings"

	"github.com/ehanz12/api-SneakHub/database"
	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/requests"
)

func CreateBrandService(r requests.BrandRequest) (*models.Brand, error) {
	tx := database.DB.Begin()
	if tx.Error != nil {
		return nil, errors.New("gagal menyambungkan ke server")
	}

	b := models.Brand{
		NamaBrand: strings.TrimSpace(r.NamaBrand),
	}
	if err := tx.Create(&b).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal membuat brand, nama brand mungkin sudah ada")
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New("terjadi kesalahan server")
	}
	return &b, nil
}

func GetBrandService() ([]models.Brand, error) {
	var brands []models.Brand
	if err := database.DB.Order("nama_brand asc").Find(&brands).Error; err != nil {
		return nil, errors.New("gagal memuat brand")
	}
	return brands, nil
}

func UpdateBrandService(brandID string, r requests.BrandRequest) (*models.Brand, error) {
	var brand models.Brand
	if err := database.DB.Where("brand_id = ?", brandID).First(&brand).Error; err != nil {
		return nil, errors.New("brand tidak ditemukan")
	}

	name := strings.TrimSpace(r.NamaBrand)
	if name != "" {
		brand.NamaBrand = name
	}
	if err := database.DB.Save(&brand).Error; err != nil {
		return nil, errors.New("gagal update brand, nama brand mungkin sudah ada")
	}
	return &brand, nil
}

func DeleteBrandService(brandID string) error {
	var brand models.Brand
	if err := database.DB.Where("brand_id = ?", brandID).First(&brand).Error; err != nil {
		return errors.New("brand tidak ditemukan")
	}

	var totalProduk int64
	if err := database.DB.Model(&models.Product{}).
		Where("brand_id = ?", brandID).Count(&totalProduk).Error; err != nil {
		return errors.New("gagal memeriksa produk brand")
	}
	if totalProduk > 0 {
		return errors.New("brand tidak dapat dihapus karena masih digunakan oleh produk")
	}

	if err := database.DB.Delete(&brand).Error; err != nil {
		return errors.New("gagal menghapus brand")
	}
	return nil
}
