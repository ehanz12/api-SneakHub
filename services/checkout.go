package services

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/ehanz12/api-SneakHub/database"
	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/requests"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

const BiayaPengirimanFlat = 15000.0

func tripayChannel(metode string) string {
	switch strings.ToUpper(strings.TrimSpace(metode)) {
	case "EWALLET":
		return "QRIS2"
	case "BANK_TRANSFER", "VA":
		return "BCAVA"
	default:
		return "QRIS2"
	}
}

func groupCartItemsBySeller(items []models.CartItem) map[string][]models.CartItem {
	groups := make(map[string][]models.CartItem)
	for _, item := range items {
		sellerID := item.Product.SellerID
		groups[sellerID] = append(groups[sellerID], item)
	}
	return groups
}

func buildAlamatSnapshot(namaPenerima, nomorTelepon, alamat, kota, provinsi, kodePos string) ([]byte, error) {
	payload := struct {
		NamaPenerima string `json:"nama_penerima"`
		NomorTelepon string `json:"nomor_telepon"`
		Alamat       string `json:"alamat"`
		Kota         string `json:"kota"`
		Provinsi     string `json:"provinsi"`
		KodePos      string `json:"kode_pos"`
	}{
		NamaPenerima: namaPenerima,
		NomorTelepon: nomorTelepon,
		Alamat:       alamat,
		Kota:         kota,
		Provinsi:     provinsi,
		KodePos:      kodePos,
	}
	return json.Marshal(payload)
}

func createOrderForSeller(tx *gorm.DB, customerID string, address models.Address, addressID string, metode string, items []models.CartItem) (*models.Order, error) {
	subtotal := 0.0
	for _, item := range items {
		subtotal += item.HargaSaatDitambahkan * float64(item.Jumlah)
	}
	total := subtotal + BiayaPengirimanFlat

	snapshot, err := buildAlamatSnapshot(address.NamaPenerima, address.NomorTelepon, address.Alamat, address.Kota, address.Provinsi, address.KodePos)
	if err != nil {
		return nil, errors.New("gagal menyimpan data alamat")
	}

	order := models.Order{
		CustomerID:       customerID,
		SellerID:         items[0].Product.SellerID,
		AddressID:        &addressID,
		StatusOrder:      "pending",
		AlamatPengiriman: datatypes.JSON(snapshot),
		MetodePembayaran: metode,
		Subtotal:         subtotal,
		BiayaPengiriman:  BiayaPengirimanFlat,
		TotalPesanan:     total,
	}
	if err := tx.Create(&order).Error; err != nil {
		return nil, errors.New("gagal membuat pesanan")
	}

	for _, item := range items {
		if err := tx.Create(&models.OrderItem{
			OrderID:            order.OrderID,
			ProductID:          item.ProductID,
			Jumlah:             item.Jumlah,
			HargaSaatTransaksi: item.HargaSaatDitambahkan,
		}).Error; err != nil {
			return nil, errors.New("gagal menyimpan item pesanan")
		}
	}

	return &order, nil
}

func createCheckoutPayment(tx *gorm.DB, order *models.Order, items []models.CartItem, customer models.User) error {
	orderItems := make([]TripayOrderItem, 0, len(items))
	for _, item := range items {
		orderItems = append(orderItems, TripayOrderItem{
			Sku:      item.ProductID,
			Name:     item.Product.NamaProduk,
			Price:    int64(item.HargaSaatDitambahkan),
			Quantity: int32(item.Jumlah),
		})
	}

	phone := ""
	if customer.NomorTelepon != nil {
		phone = *customer.NomorTelepon
	}

	reference, paymentURL, err := GetPaymentProvider().CreatePayment(
		order.OrderID,
		tripayChannel(order.MetodePembayaran),
		int64(order.TotalPesanan),
		orderItems,
		customer.Nama,
		customer.Email,
		phone,
	)
	if err != nil {
		return err
	}

	payment := models.Payment{
		OrderID:          order.OrderID,
		MetodePembayaran: order.MetodePembayaran,
		Jumlah:           order.TotalPesanan,
		StatusPembayaran: "pending",
		GatewayReference: &reference,
		PaymentURL:       &paymentURL,
	}
	if err := tx.Create(&payment).Error; err != nil {
		return errors.New("gagal menyimpan data pembayaran")
	}
	order.Payment = &payment

	return nil
}

func CheckoutService(customerID string, r requests.CheckoutRequest) ([]models.Order, error) {
	tx := database.DB.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return nil, errors.New("koneksi database bermasalah")
	}

	var cart models.Cart
	if err := tx.
		Preload("Items.Product", func(db *gorm.DB) *gorm.DB {
			return db.Select("product_id", "seller_id", "nama_produk", "stok", "status_publikasi", "harga")
		}).
		Where("customer_id = ?", customerID).
		First(&cart).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("cart tidak ditemukan")
		}
		return nil, errors.New("gagal memuat cart")
	}
	if len(cart.Items) == 0 {
		tx.Rollback()
		return nil, errors.New("cart kosong, tidak ada produk untuk di-checkout")
	}

	var address models.Address
	if err := tx.Where("address_id = ? AND user_id = ?", r.AddressID, customerID).First(&address).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("alamat pengiriman tidak ditemukan")
		}
		return nil, errors.New("gagal memuat alamat pengiriman")
	}

	for _, item := range cart.Items {
		if item.Product.StatusPublikasi != "aktif" {
			tx.Rollback()
			return nil, errors.New("produk tidak aktif: " + item.ProductID)
		}
		if item.Product.Stok < item.Jumlah {
			tx.Rollback()
			return nil, errors.New("stok produk tidak mencukupi: " + item.ProductID)
		}
		res := tx.Model(&models.Product{}).
			Where("product_id = ? AND stok >= ?", item.ProductID, item.Jumlah).
			Update("stok", gorm.Expr("stok - ?", item.Jumlah))
		if res.Error != nil {
			tx.Rollback()
			return nil, errors.New("gagal memperbarui stok produk")
		}
		if res.RowsAffected == 0 {
			tx.Rollback()
			return nil, errors.New("stok produk sudah habis: " + item.ProductID)
		}
	}

	var customer models.User
	if err := tx.Select("nama", "email", "nomor_telepon").Where("user_id = ?", customerID).First(&customer).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal memuat data pelanggan")
	}

	groups := groupCartItemsBySeller(cart.Items)
	orders := make([]models.Order, 0, len(groups))

	for _, items := range groups {
		order, err := createOrderForSeller(tx, customerID, address, address.AddressID, r.MetodePembayaran, items)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := createCheckoutPayment(tx, order, items, customer); err != nil {
			tx.Rollback()
			return nil, err
		}
		orders = append(orders, *order)
	}

	if err := tx.Where("cart_id = ?", cart.CartID).Delete(&models.CartItem{}).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("gagal mengosongkan cart")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, errors.New("data pesanan gagal tersimpan")
	}

	return orders, nil
}
