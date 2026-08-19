package services

import (
	"errors"

	"github.com/ehanz12/api-SneakHub/database"
	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/requests"
	"gorm.io/gorm"
)

func getOrCreateCart(tx *gorm.DB, customerID string) (*models.Cart, error) {
	var cart models.Cart
	err := tx.Where("customer_id = ?", customerID).First(&cart).Error
	if err == nil {
		return &cart, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("gagal memuat cart")
	}

	cart = models.Cart{CustomerID: customerID}
	if err := tx.Create(&cart).Error; err != nil {
		return nil, errors.New("gagal membuat cart")
	}
	return &cart, nil
}

func AddCartItemsService(customerID string, r requests.AddCartItemsRequest) ([]models.CartItem, float64, error) {
	tx := database.DB.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return nil, 0, errors.New("gagal koneksi ke database")
	}

	cart, err := getOrCreateCart(tx, customerID)
	if err != nil {
		tx.Rollback()
		return nil, 0, err
	}

	added := make([]models.CartItem, 0, len(r.Items))
	var total float64

	for _, item := range r.Items {
		var product models.Product
		if err := tx.Select("product_id", "harga", "stok", "status_publikasi").
			Where("product_id = ? AND status_publikasi = ?", item.ProductID, "aktif").
			First(&product).Error; err != nil {
			tx.Rollback()
			return nil, 0, errors.New("produk tidak ditemukan atau tidak aktif: " + item.ProductID)
		}

		var existing models.CartItem
		err := tx.Where("cart_id = ? AND product_id = ?", cart.CartID, item.ProductID).First(&existing).Error
		if err == nil {
			newJumlah := existing.Jumlah + item.Jumlah
			if newJumlah > product.Stok {
				tx.Rollback()
				return nil, 0, errors.New("stok produk tidak mencukupi: " + item.ProductID)
			}
			if err := tx.Model(&existing).Update("jumlah", newJumlah).Error; err != nil {
				tx.Rollback()
				return nil, 0, errors.New("gagal memperbarui item cart")
			}
			existing.Jumlah = newJumlah
			added = append(added, existing)
			total += existing.HargaSaatDitambahkan * float64(existing.Jumlah)
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return nil, 0, errors.New("gagal memuat item cart")
		}

		if item.Jumlah > product.Stok {
			tx.Rollback()
			return nil, 0, errors.New("stok produk tidak mencukupi: " + item.ProductID)
		}

		cartItem := models.CartItem{
			CartID:               cart.CartID,
			ProductID:            item.ProductID,
			Jumlah:               item.Jumlah,
			HargaSaatDitambahkan: product.Harga,
		}
		if err := tx.Create(&cartItem).Error; err != nil {
			tx.Rollback()
			return nil, 0, errors.New("gagal menambahkan item ke cart")
		}
		added = append(added, cartItem)
		total += product.Harga * float64(item.Jumlah)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, 0, errors.New("data cart gagal tersimpan")
	}

	return added, total, nil
}

func GetCartService(customerID string) (*models.Cart, error) {
	var cart models.Cart
	err := database.DB.
		Preload("Items.Product.Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("urutan_tampil asc")
		}).
		Where("customer_id = ?", customerID).
		First(&cart).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &cart, nil
	}
	if err != nil {
		return nil, errors.New("gagal memuat cart")
	}
	return &cart, nil
}

func UpdateCartItemService(customerID, cartItemID string, jumlah int) (*models.CartItem, error) {
	tx := database.DB.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return nil, errors.New("gagal koneksi ke database")
	}

	var cart models.Cart
	if err := tx.Select("cart_id").Where("customer_id = ?", customerID).First(&cart).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("cart tidak ditemukan")
	}

	var item models.CartItem
	if err := tx.Where("cart_item_id = ? AND cart_id = ?", cartItemID, cart.CartID).First(&item).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("item cart tidak ditemukan")
	}

	var product models.Product
	if err := tx.Select("stok").Where("product_id = ?", item.ProductID).First(&product).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("produk tidak ditemukan")
	}
	if jumlah > product.Stok {
		tx.Rollback()
		return nil, errors.New("stok produk tidak mencukupi")
	}

	if err := tx.Model(&item).Update("jumlah", jumlah).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal memperbarui item cart")
	}
	item.Jumlah = jumlah

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New("data cart gagal tersimpan")
	}

	return &item, nil
}

func DeleteCartItemService(customerID, cartItemID string) error {
	tx := database.DB.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return errors.New("gagal koneksi ke database")
	}

	var cart models.Cart
	if err := tx.Select("cart_id").Where("customer_id = ?", customerID).First(&cart).Error; err != nil {
		tx.Rollback()
		return errors.New("cart tidak ditemukan")
	}

	var item models.CartItem
	if err := tx.Where("cart_item_id = ? AND cart_id = ?", cartItemID, cart.CartID).First(&item).Error; err != nil {
		tx.Rollback()
		return errors.New("item cart tidak ditemukan")
	}

	if err := tx.Delete(&item).Error; err != nil {
		tx.Rollback()
		return errors.New("gagal menghapus item cart")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errors.New("data cart gagal tersimpan")
	}

	return nil
}
