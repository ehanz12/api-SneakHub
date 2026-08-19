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
	"github.com/ehanz12/api-SneakHub/models"
	"gorm.io/datatypes"
)

type alamatPengirimanSnapshot struct {
	NamaPenerima string `json:"nama_penerima"`
	NomorTelepon string `json:"nomor_telepon"`
	Alamat       string `json:"alamat"`
	Kota         string `json:"kota"`
	Provinsi     string `json:"provinsi"`
	KodePos      string `json:"kode_pos"`
}

func parseAlamatSnapshot(raw datatypes.JSON) (alamatPengirimanSnapshot, error) {
	var alamat alamatPengirimanSnapshot
	if err := json.Unmarshal(raw, &alamat); err != nil {
		return alamat, errors.New("data alamat pengiriman tidak valid")
	}
	return alamat, nil
}

type biteshipItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       int64  `json:"value"`
	Quantity    int    `json:"quantity"`
	Weight      int    `json:"weight"`
}

type biteshipCreateOrderRequest struct {
	ShipperContactName      string            `json:"shipper_contact_name"`
	ShipperContactPhone     string            `json:"shipper_contact_phone"`
	ShipperContactEmail     string            `json:"shipper_contact_email"`
	ShipperOrganization     string            `json:"shipper_organization"`
	OriginContactName       string            `json:"origin_contact_name"`
	OriginContactPhone      string            `json:"origin_contact_phone"`
	OriginAddress           string            `json:"origin_address"`
	OriginNote              string            `json:"origin_note"`
	OriginPostalCode        string            `json:"origin_postal_code"`
	DestinationContactName  string            `json:"destination_contact_name"`
	DestinationContactPhone string            `json:"destination_contact_phone"`
	DestinationContactEmail string            `json:"destination_contact_email"`
	DestinationAddress      string            `json:"destination_address"`
	DestinationPostalCode   string            `json:"destination_postal_code"`
	CourierCompany          string            `json:"courier_company"`
	CourierType             string            `json:"courier_type"`
	CourierInsurance        int               `json:"courier_insurance"`
	DeliveryType            string            `json:"delivery_type"`
	OrderNote               string            `json:"order_note"`
	Metadata                map[string]string `json:"metadata"`
	Items                   []biteshipItem    `json:"items"`
}

type biteshipCreateOrderResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		OrderID   string `json:"order_id"`
		WaybillID string `json:"waybill_id"`
		Courier   struct {
			Company    string `json:"company"`
			TrackingID string `json:"tracking_id"`
			WaybillID  string `json:"waybill_id"`
		} `json:"courier"`
		Status string `json:"status"`
	} `json:"data"`
}

type biteshipTrackingResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Status  string `json:"status"`
		Courier struct {
			Company    string `json:"company"`
			TrackingID string `json:"tracking_id"`
		} `json:"courier"`
		History []models.ShipmentTrackingEvent `json:"history"`
	} `json:"data"`
}

func biteshipRequest(method, path string, payload interface{}) ([]byte, error) {
	apiKey := strings.TrimSpace(config.AppConfig.BiteshipAPIKey)
	if apiKey == "" {
		return nil, errors.New("BITESHIP_API_KEY belum dikonfigurasi")
	}

	var body io.Reader
	if payload != nil {
		buf, err := json.Marshal(payload)
		if err != nil {
			return nil, errors.New("gagal menyusun request")
		}
		body = bytes.NewReader(buf)
	}

	req, err := http.NewRequest(method, biteshipBaseURL+path, body)
	if err != nil {
		return nil, errors.New("gagal menyusun request")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Biteship "+apiKey)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("gagal terhubung ke layanan Biteship")
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("gagal membaca respons Biteship")
	}
	return data, nil
}

func CreateBiteshipOrder(order *models.Order, shipment *models.Shipment, items []models.OrderItem, store models.Seller, sellerUser models.User) (resi, waybillID string, err error) {
	alamat, err := parseAlamatSnapshot(order.AlamatPengiriman)
	if err != nil {
		return "", "", err
	}

	if store.AlamatAsal == nil || strings.TrimSpace(*store.AlamatAsal) == "" {
		return "", "", errors.New("alamat_asal toko belum diisi")
	}
	if store.KodePosAsal == nil || strings.TrimSpace(*store.KodePosAsal) == "" {
		return "", "", errors.New("kode_pos_asal toko belum diisi")
	}
	sellerPhone := ""
	if sellerUser.NomorTelepon != nil {
		sellerPhone = strings.TrimSpace(*sellerUser.NomorTelepon)
	}
	if sellerPhone == "" {
		return "", "", errors.New("nomor telepon seller belum diisi")
	}
	if strings.TrimSpace(alamat.NomorTelepon) == "" {
		return "", "", errors.New("nomor telepon penerima belum diisi")
	}

	kodeKurir := strings.ToLower(strings.TrimSpace(shipment.Kurir))
	if kodeKurir == "" {
		return "", "", errors.New("kurir belum dipilih")
	}
	tipeLayanan := ""
	if shipment.Service != nil {
		tipeLayanan = strings.ToLower(strings.TrimSpace(*shipment.Service))
	}
	if tipeLayanan == "" {
		return "", "", errors.New("tipe layanan kurir tidak ditemukan")
	}

	originAddress := strings.TrimSpace(*store.AlamatAsal)
	if store.KotaAsal != nil && strings.TrimSpace(*store.KotaAsal) != "" {
		originAddress += ", " + strings.TrimSpace(*store.KotaAsal)
	}
	destinationAddress := strings.TrimSpace(alamat.Alamat)
	if strings.TrimSpace(alamat.Kota) != "" {
		destinationAddress += ", " + strings.TrimSpace(alamat.Kota)
	}
	if strings.TrimSpace(alamat.Provinsi) != "" {
		destinationAddress += ", " + strings.TrimSpace(alamat.Provinsi)
	}

	biteshipItems := make([]biteshipItem, 0, len(items))
	for _, item := range items {
		berat := item.Product.Berat
		if berat <= 0 {
			berat = 500
		}
		biteshipItems = append(biteshipItems, biteshipItem{
			Name:        item.Product.NamaProduk,
			Description: item.Product.NamaProduk,
			Value:       int64(item.HargaSaatTransaksi),
			Quantity:    item.Jumlah,
			Weight:      berat * item.Jumlah,
		})
	}

	payload := biteshipCreateOrderRequest{
		ShipperContactName:      strings.TrimSpace(sellerUser.Nama),
		ShipperContactPhone:     sellerPhone,
		ShipperContactEmail:     strings.TrimSpace(sellerUser.Email),
		ShipperOrganization:     strings.TrimSpace(store.NamaToko),
		OriginContactName:       strings.TrimSpace(sellerUser.Nama),
		OriginContactPhone:      sellerPhone,
		OriginAddress:           originAddress,
		OriginNote:              "Pesanan dari " + store.NamaToko,
		OriginPostalCode:        strings.TrimSpace(*store.KodePosAsal),
		DestinationContactName:  strings.TrimSpace(alamat.NamaPenerima),
		DestinationContactPhone: strings.TrimSpace(alamat.NomorTelepon),
		DestinationContactEmail: strings.TrimSpace(sellerUser.Email),
		DestinationAddress:      destinationAddress,
		DestinationPostalCode:   strings.TrimSpace(alamat.KodePos),
		CourierCompany:          kodeKurir,
		CourierType:             tipeLayanan,
		CourierInsurance:        0,
		DeliveryType:            "now",
		OrderNote:               "Order " + order.OrderID,
		Metadata: map[string]string{
			"order_id": order.OrderID,
		},
		Items: biteshipItems,
	}

	data, err := biteshipRequest(http.MethodPost, "/v1/orders", payload)
	if err != nil {
		return "", "", err
	}

	var result biteshipCreateOrderResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return "", "", errors.New("respons booking Biteship tidak valid")
	}
	if !result.Success {
		msg := strings.TrimSpace(result.Message)
		if msg == "" {
			msg = "booking kurir Biteship gagal"
		}
		return "", "", errors.New(msg)
	}

	resi = strings.TrimSpace(result.Data.Courier.TrackingID)
	if resi == "" {
		resi = strings.TrimSpace(result.Data.WaybillID)
	}
	if resi == "" {
		return "", "", errors.New("Biteship tidak mengembalikan nomor resi")
	}
	waybillID = strings.TrimSpace(result.Data.Courier.WaybillID)
	if waybillID == "" {
		waybillID = strings.TrimSpace(result.Data.WaybillID)
	}

	return resi, waybillID, nil
}

func TrackBiteshipShipment(trackingID string) (string, []models.ShipmentTrackingEvent, error) {
	trackingID = strings.TrimSpace(trackingID)
	if trackingID == "" {
		return "", nil, errors.New("tracking id kosong")
	}

	data, err := biteshipRequest(http.MethodGet, "/v1/trackings/"+trackingID, nil)
	if err != nil {
		return "", nil, err
	}

	var result biteshipTrackingResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return "", nil, errors.New("respons tracking Biteship tidak valid")
	}
	if !result.Success {
		msg := strings.TrimSpace(result.Message)
		if msg == "" {
			msg = "tracking Biteship gagal"
		}
		return "", nil, errors.New(msg)
	}

	return strings.ToLower(strings.TrimSpace(result.Data.Status)), result.Data.History, nil
}

func mapBiteshipCourierStatus(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "delivered":
		return "sampai"
	case "delivering", "on_delivery", "picked_up", "picked", "picking":
		return "dalam_perjalanan"
	case "cancelled":
		return "menunggu"
	default:
		return "dikirim"
	}
}
