package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Address struct {
	AddressID    string `gorm:"column:address_id;type:char(36);primaryKey" json:"address_id"`
	UserID       string `gorm:"column:user_id;type:char(36);not null" json:"user_id"`
	NamaPenerima string `gorm:"column:nama_penerima;type:varchar(100);not null" json:"nama_penerima"`
	NomorTelepon string `gorm:"column:nomor_telepon;type:varchar(30);not null" json:"nomor_telepon"`
	Alamat       string `gorm:"column:alamat;type:text;not null" json:"alamat"`
	Kota         string `gorm:"column:kota;type:varchar(100);not null" json:"kota"`
	Provinsi     string `gorm:"column:provinsi;type:varchar(100);not null" json:"provinsi"`
	KodePos      string `gorm:"column:kode_pos;type:varchar(10);not null" json:"kode_pos"`
	IsDefault    bool   `gorm:"column:is_default;not null;default:false" json:"is_default"`
	Timestamps

	User User `gorm:"foreignKey:UserID;references:UserID" json:"user,omitempty"`
}

func (a *Address) BeforeCreate(tx *gorm.DB) error {
	if a.AddressID == "" {
		a.AddressID = uuid.NewString()
	}
	return nil
}
