package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ehanz12/api-SneakHub/config"
	"github.com/ehanz12/api-SneakHub/database"
	"github.com/ehanz12/api-SneakHub/models"
	"github.com/ehanz12/api-SneakHub/requests"
	"gorm.io/gorm"
)

const biteshipBaseURL = "https://api.biteship.com"

type ShippingOption struct {
	Kurir       string  `json:"kurir"`
	Service     string  `json:"service"`
	ServiceCode string  `json:"service_code,omitempty"`
	Biaya       float64 `json:"biaya"`
	Estimasi    string  `json:"estimasi"`
	IsFallback  bool    `json:"is_fallback"`
}

type SellerShippingRates struct {
	SellerID string           `json:"seller_id"`
	NamaToko string           `json:"nama_toko"`
	Berat    int              `json:"berat"`
	KotaAsal *string          `json:"kota_asal,omitempty"`
	Options  []ShippingOption `json:"options"`
}

type biteshipRateRequest struct {
	OriginPostalCode      string `json:"origin_postal_code"`
	DestinationPostalCode string `json:"destination_postal_code"`
	Couriers              string `json:"couriers"`
	Weight                int    `json:"weight"`
}

type biteshipRateResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Rates   []struct {
		CourierName        string  `json:"courier_name"`
		CourierServiceName string  `json:"courier_service_name"`
		CourierServiceCode string  `json:"courier_service_code"`
		Duration           string  `json:"duration"`
		ShippingAmount     float64 `json:"shipping_amount"`
	} `json:"rates"`
}

func hitungBeratGroup(items []models.CartItem) int {
	total := 0
	for _, item := range items {
		total += item.Product.Berat * item.Jumlah
	}
	return total
}

func fetchBiteshipRates(originPostalCode, destinationPostalCode string, weight int) ([]ShippingOption, error) {
	apiKey := config.AppConfig.BiteshipAPIKey
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("BITESHIP_API_KEY belum dikonfigurasi")
	}
	if weight <= 0 {
		weight = 500
	}

	payload, err := json.Marshal(biteshipRateRequest{
		OriginPostalCode:      originPostalCode,
		DestinationPostalCode: destinationPostalCode,
		Couriers:              "jne,sicepat,anteraja,jnt",
		Weight:                weight,
	})
	if err != nil {
		return nil, errors.New("gagal menyusun request ongkir")
	}

	req, err := http.NewRequest(http.MethodPost, biteshipBaseURL+"/v1/rates/couriers", bytes.NewReader(payload))
	if err != nil {
		return nil, errors.New("gagal menyusun request ongkir")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Biteship "+apiKey)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("gagal terhubung ke layanan ongkir")
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("gagal membaca respons ongkir")
	}

	var result biteshipRateResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, errors.New("respons ongkir tidak valid")
	}
	if !result.Success {
		msg := strings.TrimSpace(result.Message)
		if msg == "" {
			msg = "layanan ongkir gagal"
		}
		return nil, errors.New(msg)
	}

	options := make([]ShippingOption, 0, len(result.Rates))
	for _, rate := range result.Rates {
		options = append(options, ShippingOption{
			Kurir:       strings.ToLower(rate.CourierName),
			Service:     rate.CourierServiceName,
			ServiceCode: strings.ToLower(rate.CourierServiceCode),
			Biaya:       rate.ShippingAmount,
			Estimasi:    rate.Duration,
		})
	}
	if len(options) == 0 {
		return nil, errors.New("tidak ada pilihan kurir tersedia")
	}
	return options, nil
}

func GetShippingRatesService(customerID, addressID string) ([]SellerShippingRates, error) {
	var cart models.Cart
	if err := database.DB.
		Preload("Items.Product", func(db *gorm.DB) *gorm.DB {
			return db.Select("product_id", "seller_id", "nama_produk", "berat", "status_publikasi")
		}).
		Where("customer_id = ?", customerID).
		First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("cart tidak ditemukan")
		}
		return nil, errors.New("gagal memuat cart")
	}
	if len(cart.Items) == 0 {
		return nil, errors.New("cart kosong, tidak ada produk untuk dihitung")
	}

	var address models.Address
	if err := database.DB.Where("address_id = ? AND user_id = ?", addressID, customerID).
		First(&address).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("alamat pengiriman tidak ditemukan")
		}
		return nil, errors.New("gagal memuat alamat pengiriman")
	}

	groups := groupCartItemsBySeller(cart.Items)
	result := make([]SellerShippingRates, 0, len(groups))

	for _, items := range groups {
		var store models.Seller
		if err := database.DB.Select("seller_id", "nama_toko", "kode_pos_asal", "kota_asal").
			Where("seller_id = ?", items[0].Product.SellerID).First(&store).Error; err != nil {
			return nil, errors.New("data toko seller tidak ditemukan")
		}

		rates := SellerShippingRates{
			SellerID: store.SellerID,
			NamaToko: store.NamaToko,
			Berat:    hitungBeratGroup(items),
			KotaAsal: store.KotaAsal,
		}

		if config.AppConfig.BiteshipAPIKey == "" || store.KodePosAsal == nil || strings.TrimSpace(*store.KodePosAsal) == "" {
			rates.Options = flatShippingOption()
		} else {
			options, err := fetchBiteshipRates(*store.KodePosAsal, address.KodePos, rates.Berat)
			if err != nil {
				rates.Options = flatShippingOption()
			} else {
				rates.Options = options
			}
		}
		result = append(result, rates)
	}

	return result, nil
}

func flatShippingOption() []ShippingOption {
	return []ShippingOption{{
		Kurir:      "flat",
		Service:    "ongkir tetap",
		Biaya:      BiayaPengirimanFlat,
		Estimasi:   "1-3 hari",
		IsFallback: true,
	}}
}

func resolveShippingCostForGroup(tx *gorm.DB, items []models.CartItem, destinationKodePos string, requested []requests.CheckoutShippingRequest) (float64, string, string) {
	sellerID := items[0].Product.SellerID

	var store models.Seller
	if err := tx.Select("seller_id", "kode_pos_asal").
		Where("seller_id = ?", sellerID).First(&store).Error; err != nil {
		return BiayaPengirimanFlat, "flat", ""
	}
	if config.AppConfig.BiteshipAPIKey == "" || store.KodePosAsal == nil || strings.TrimSpace(*store.KodePosAsal) == "" {
		return BiayaPengirimanFlat, "flat", ""
	}

	options, err := fetchBiteshipRates(*store.KodePosAsal, destinationKodePos, hitungBeratGroup(items))
	if err != nil {
		return BiayaPengirimanFlat, "flat", ""
	}

	kodeKurir := ""
	for _, req := range requested {
		if req.SellerID == sellerID {
			kodeKurir = strings.ToLower(strings.TrimSpace(req.Kurir))
			break
		}
	}

	if kodeKurir != "" {
		for _, opt := range options {
			if strings.Contains(strings.ToLower(opt.Kurir), kodeKurir) ||
				strings.Contains(kodeKurir, strings.ToLower(opt.Kurir)) {
				return opt.Biaya, opt.Kurir, opt.ServiceCode
			}
		}
	}

	cheapest := options[0]
	for _, opt := range options[1:] {
		if opt.Biaya < cheapest.Biaya {
			cheapest = opt
		}
	}
	return cheapest.Biaya, cheapest.Kurir, cheapest.ServiceCode
}
