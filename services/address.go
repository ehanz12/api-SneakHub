package services

import (
	"errors"

	"github.com/ehanz12/api-SneakHub/database"
	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/requests"
	"gorm.io/gorm"
)

func unsetOtherDefaults(tx *gorm.DB, userID, keepID string) error {
	return tx.Model(&models.Address{}).
		Where("user_id = ? AND address_id <> ? AND is_default = ?", userID, keepID, true).
		Update("is_default", false).Error
}

func CreateAddressService(userID string, r requests.AddressCreateRequest) (*models.Address, error) {
	tx := database.DB.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return nil, errors.New("koneksi ke database gagal")
	}

	if r.IsDefault {
		if err := unsetOtherDefaults(tx, userID, ""); err != nil {
			tx.Rollback()
			return nil, errors.New("gagal memperbarui alamat default")
		}
	}

	address := models.Address{
		UserID:       userID,
		NamaPenerima: r.NamaPenerima,
		NomorTelepon: r.NomorTelepon,
		Alamat:       r.Alamat,
		Kota:         r.Kota,
		Provinsi:     r.Provinsi,
		KodePos:      r.KodePos,
		IsDefault:    r.IsDefault,
	}
	if err := tx.Create(&address).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal menambahkan alamat")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New("perubahan alamat gagal disimpan")
	}

	return &address, nil
}

func GetAddressesService(userID string) ([]models.Address, error) {
	var addresses []models.Address
	err := database.DB.
		Where("user_id = ?", userID).
		Order("is_default desc, created_at desc").
		Find(&addresses).Error
	if err != nil {
		return nil, errors.New("gagal memuat alamat")
	}
	return addresses, nil
}

func GetAddressService(userID, addressID string) (*models.Address, error) {
	var address models.Address
	err := database.DB.
		Where("address_id = ? AND user_id = ?", addressID, userID).
		First(&address).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("alamat tidak ditemukan")
	}
	if err != nil {
		return nil, errors.New("gagal memuat alamat")
	}
	return &address, nil
}

func UpdateAddressService(userID, addressID string, r requests.AddressUpdateRequest) (*models.Address, error) {
	tx := database.DB.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return nil, errors.New("koneksi ke database gagal")
	}

	var address models.Address
	if err := tx.Where("address_id = ? AND user_id = ?", addressID, userID).First(&address).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("alamat tidak ditemukan")
		}
		return nil, errors.New("gagal memuat alamat")
	}

	updates := make(map[string]interface{})
	if r.NamaPenerima != nil {
		updates["nama_penerima"] = *r.NamaPenerima
	}
	if r.NomorTelepon != nil {
		updates["nomor_telepon"] = *r.NomorTelepon
	}
	if r.Alamat != nil {
		updates["alamat"] = *r.Alamat
	}
	if r.Kota != nil {
		updates["kota"] = *r.Kota
	}
	if r.Provinsi != nil {
		updates["provinsi"] = *r.Provinsi
	}
	if r.KodePos != nil {
		updates["kode_pos"] = *r.KodePos
	}
	if r.IsDefault != nil {
		updates["is_default"] = *r.IsDefault
		if *r.IsDefault {
			if err := unsetOtherDefaults(tx, userID, addressID); err != nil {
				tx.Rollback()
				return nil, errors.New("gagal memperbarui alamat default")
			}
		}
	}

	if err := tx.Model(&address).Updates(updates).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal memperbarui alamat")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New("perubahan alamat gagal disimpan")
	}

	return &address, nil
}

func DeleteAddressService(userID, addressID string) error {
	tx := database.DB.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return errors.New("koneksi ke database gagal")
	}

	res := tx.Where("address_id = ? AND user_id = ?", addressID, userID).Delete(&models.Address{})
	if res.Error != nil {
		tx.Rollback()
		return errors.New("gagal menghapus alamat")
	}
	if res.RowsAffected == 0 {
		tx.Rollback()
		return errors.New("alamat tidak ditemukan")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errors.New("perubahan alamat gagal disimpan")
	}

	return nil
}
