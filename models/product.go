package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Product struct {
	ProductID       string         `gorm:"column:product_id;type:char(36);primaryKey" json:"product_id"`
	SellerID        string         `gorm:"column:seller_id;type:char(36);not null" json:"seller_id"`
	BrandID         string         `gorm:"column:brand_id;type:char(36);not null" json:"brand_id"`
	CategoryID      string         `gorm:"column:category_id;type:char(36);not null" json:"category_id"`
	NamaProduk      string         `gorm:"column:nama_produk;type:varchar(200);not null" json:"nama_produk"`
	Kondisi         string         `gorm:"column:kondisi;type:enum('new','used','refurbished');not null" json:"kondisi"`
	Deskripsi       *string        `gorm:"column:deskripsi;type:text" json:"deskripsi,omitempty"`
	Harga           float64        `gorm:"column:harga;type:decimal(15,2);not null" json:"harga"`
	Stok            int            `gorm:"column:stok;not null;default:0" json:"stok"`
	Berat           int            `gorm:"column:berat;not null;default:0" json:"berat"`
	UkuranTersedia  datatypes.JSON `gorm:"column:ukuran_tersedia;type:json"`
	ConditionScore  *float64       `gorm:"column:condition_score;type:decimal(5,2)" json:"condition_score,omitempty"`
	StatusPublikasi string         `gorm:"column:status_publikasi;type:enum('aktif','draft','nonaktif');default:draft" json:"status_publikasi"`
	Timestamps

	Seller       Seller           `gorm:"foreignKey:SellerID;references:SellerID" json:"seller,omitempty"`
	Brand        Brand            `gorm:"foreignKey:BrandID;references:BrandID" json:"brand,omitempty"`
	Category     Category         `gorm:"foreignKey:CategoryID;references:CategoryID" json:"category,omitempty"`
	Images       []ProductImage   `gorm:"foreignKey:ProductID;references:ProductID" json:"images,omitempty"`
	Embeddings   []ImageEmbedding `gorm:"foreignKey:ProductID;references:ProductID" json:"embeddings,omitempty"`
	Condition    *ConditionScore  `gorm:"foreignKey:ProductID;references:ProductID" json:"condition,omitempty"`
	PriceHistory []PriceHistory   `gorm:"foreignKey:ProductID;references:ProductID" json:"price_history,omitempty"`
	Wishlists    []Wishlist       `gorm:"foreignKey:ProductID;references:ProductID" json:"wishlists,omitempty"`
	CartItems    []CartItem       `gorm:"foreignKey:ProductID;references:ProductID" json:"cart_items,omitempty"`
	OrderItems   []OrderItem      `gorm:"foreignKey:ProductID;references:ProductID" json:"order_items,omitempty"`
	Reviews      []Review         `gorm:"foreignKey:ProductID;references:ProductID" json:"reviews,omitempty"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ProductID == "" {
		p.ProductID = uuid.NewString()
	}
	return nil
}
func (Product) TableName() string { return "products" }
