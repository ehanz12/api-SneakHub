package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	UserID           string         `gorm:"column:user_id;type:char(36);primaryKey" json:"user_id"`
	Nama             string         `gorm:"column:nama;type:varchar(100);not null" json:"nama"`
	Email            string         `gorm:"column:email;type:varchar(150);unique;not null" json:"email"`
	NomorTelepon     *string        `gorm:"column:nomor_telepon;type:varchar(30)" json:"nomor_telepon,omitempty"`
	KataSandiHash    string         `gorm:"column:kata_sandi_hash;type:varchar(255);not null" json:"-"`
	Peran            string         `gorm:"column:peran;type:enum('customer','seller','admin');default:customer" json:"peran"`
	StatusAkun       string         `gorm:"column:status_akun;type:enum('aktif','tidak_aktif','blokir');default:aktif" json:"status_akun"`
	PreferensiUkuran datatypes.JSON `gorm:"column:preferensi_ukuran;type:json" json:"preferensi_ukuran,omitempty"`
	BrandFavorit     datatypes.JSON `gorm:"column:brand_favorit;type:json" json:"brand_favorit,omitempty"`
	Timestamps

	Seller          *Seller              `gorm:"foreignKey:UserID;references:UserID" json:"seller,omitempty"`
	Addresses       []Address            `gorm:"foreignKey:UserID;references:UserID" json:"addresses,omitempty"`
	Activities      []UserActivity       `gorm:"foreignKey:UserID;references:UserID" json:"activities,omitempty"`
	Orders          []Order              `gorm:"foreignKey:CustomerID;references:UserID" json:"orders,omitempty"`
	Wishlists       []Wishlist           `gorm:"foreignKey:CustomerID;references:UserID" json:"wishlists,omitempty"`
	Cart            *Cart                `gorm:"foreignKey:CustomerID;references:UserID" json:"cart,omitempty"`
	Notifications   []Notification       `gorm:"foreignKey:UserID;references:UserID" json:"notifications,omitempty"`
	Recommendations []RecommendationData `gorm:"foreignKey:UserID;references:UserID" json:"recommendations,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.UserID == "" {
		u.UserID = uuid.NewString()
	}
	return nil
}

func (User) TableName() string { return "users" }
