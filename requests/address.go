package requests

import "strings"

type AddressCreateRequest struct {
	NamaPenerima string `json:"nama_penerima"`
	NomorTelepon string `json:"nomor_telepon"`
	Alamat       string `json:"alamat"`
	Kota         string `json:"kota"`
	Provinsi     string `json:"provinsi"`
	KodePos      string `json:"kode_pos"`
	IsDefault    bool   `json:"is_default"`
}

func (r *AddressCreateRequest) Validate() map[string]string {
	errs := make(map[string]string)

	if strings.TrimSpace(r.NamaPenerima) == "" {
		errs["nama_penerima"] = "nama penerima harus diisi"
	}
	if strings.TrimSpace(r.Alamat) == "" {
		errs["alamat"] = "alamat harus diisi"
	}
	if strings.TrimSpace(r.Kota) == "" {
		errs["kota"] = "kota harus diisi"
	}

	if strings.TrimSpace(r.Provinsi) == "" {
		errs["provinsi"] = "provinsi wajib diisi"
	}

	if strings.TrimSpace(r.KodePos) == "" {
		errs["kode_pos"] = "kode pos wajib diisi"
	}
	return errs
}

type AddressUpdateRequest struct {
	NamaPenerima *string `json:"nama_penerima"`
	NomorTelepon *string `json:"nomor_telepon"`
	Alamat       *string `json:"alamat"`
	Kota         *string `json:"kota"`
	Provinsi     *string `json:"provinsi"`
	KodePos      *string `json:"kode_pos"`
	IsDefault    *bool   `json:"is_default"`
}

func (r *AddressUpdateRequest) Validate() map[string]string {
	errs := make(map[string]string)

	if r.NamaPenerima == nil && r.NomorTelepon == nil && r.Alamat == nil &&
		r.Kota == nil && r.Provinsi == nil && r.KodePos == nil && r.IsDefault == nil {
		errs["request"] = "minimal satu field wajib diisi"
		return errs
	}

	if r.NamaPenerima != nil && strings.TrimSpace(*r.NamaPenerima) == "" {
		errs["nama_penerima"] = "nama penerima harus diisi"
	}
	if r.NomorTelepon != nil && strings.TrimSpace(*r.NomorTelepon) == "" {
		errs["nomor_telepon"] = "nomor telepon harus diisi"
	}
	if r.Alamat != nil && strings.TrimSpace(*r.Alamat) == "" {
		errs["alamat"] = "alamat harus diisi"
	}
	if r.Kota != nil && strings.TrimSpace(*r.Kota) == "" {
		errs["kota"] = "kota harus diisi"
	}
	if r.Provinsi != nil && strings.TrimSpace(*r.Provinsi) == "" {
		errs["provinsi"] = "provinsi wajib diisi"
	}
	if r.KodePos != nil && strings.TrimSpace(*r.KodePos) == "" {
		errs["kode_pos"] = "kode pos wajib diisi"
	}
	return errs
}
