package requests

import (
	"fmt"
	"strings"
)

type OrderAlamatRequest struct {
	NamaPenerima string `json:"nama_penerima"`
	NomorTelepon string `json:"nomor_telepon"`
	Alamat       string `json:"alamat"`
	Kota         string `json:"kota"`
	Provinsi     string `json:"provinsi"`
	KodePos      string `json:"kode_pos"`
}

type OrderItemRequest struct {
	ProductID string `json:"product_id"`
	Jumlah    int    `json:"jumlah"`
}

type CreateOrderRequest struct {
	SellerID         string             `json:"seller_id"`
	CustomerID       string             `json:"customer_id"`
	MetodePembayaran string             `json:"metode_pembayaran"`
	AlamatPengiriman OrderAlamatRequest `json:"alamat_pengiriman"`
	Items            []OrderItemRequest `json:"items"`
}

func (r *CreateOrderRequest) Validate() map[string]string {
	errs := make(map[string]string)

	if strings.TrimSpace(r.SellerID) == "" {
		errs["seller_id"] = "seller_id wajib diisi"
	}
	if strings.TrimSpace(r.MetodePembayaran) == "" {
		errs["metode_pembayaran"] = "metode_pembayaran wajib diisi"
	}

	if strings.TrimSpace(r.AlamatPengiriman.NamaPenerima) == "" {
		errs["alamat_pengiriman.nama_penerima"] = "nama penerima wajib diisi"
	}
	if strings.TrimSpace(r.AlamatPengiriman.Alamat) == "" {
		errs["alamat_pengiriman.alamat"] = "alamat wajib diisi"
	}
	if strings.TrimSpace(r.AlamatPengiriman.Kota) == "" {
		errs["alamat_pengiriman.kota"] = "kota wajib diisi"
	}
	if strings.TrimSpace(r.AlamatPengiriman.Provinsi) == "" {
		errs["alamat_pengiriman.provinsi"] = "provinsi wajib diisi"
	}
	if strings.TrimSpace(r.AlamatPengiriman.KodePos) == "" {
		errs["alamat_pengiriman.kode_pos"] = "kode pos wajib diisi"
	}

	if len(r.Items) == 0 {
		errs["items"] = "minimal satu item wajib diisi"
		return errs
	}

	seen := make(map[string]bool)
	for i, item := range r.Items {
		if strings.TrimSpace(item.ProductID) == "" {
			errs[fmt.Sprintf("items[%d].product_id", i)] = "product_id wajib diisi"
		} else if seen[item.ProductID] {
			errs[fmt.Sprintf("items[%d].product_id", i)] = "product_id duplikat"
		} else {
			seen[item.ProductID] = true
		}
		if item.Jumlah < 1 {
			errs[fmt.Sprintf("items[%d].jumlah", i)] = "jumlah minimal 1"
		}
	}

	return errs
}

type UpdateOrderStatusRequest struct {
	StatusOrder *string `json:"status_order"`
}

func (r *UpdateOrderStatusRequest) Validate() map[string]string {
	errs := make(map[string]string)

	if r.StatusOrder == nil {
		errs["status_order"] = "status_order wajib diisi"
		return errs
	}

	switch strings.ToLower(strings.TrimSpace(*r.StatusOrder)) {
	case "pending", "diproses", "dikirim", "selesai", "dibatalkan":
	default:
		errs["status_order"] = "status_order harus salah satu dari: pending, diproses, dikirim, selesai, dibatalkan"
	}

	return errs
}
