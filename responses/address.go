package responses

import "time"

type AddressResponse struct {
	AddressID    string    `json:"address_id"`
	NamaPenerima string    `json:"nama_penerima"`
	NomorTelepon string    `json:"nomor_telepon"`
	Alamat       string    `json:"alamat"`
	Kota         string    `json:"kota"`
	Provinsi     string    `json:"provinsi"`
	KodePos      string    `json:"kode_pos"`
	IsDefault    bool      `json:"is_default"`
	CreatedAt    time.Time `json:"created_at"`
}
