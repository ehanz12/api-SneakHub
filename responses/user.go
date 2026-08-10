package responses

type UserResponse struct {
	UserID       string  `json:"user_id"`
	Nama         string  `json:"nama"`
	Email        string  `json:"email"`
	NomorTelepon *string `json:"nomor_telepon"`
	Peran        string  `json:"peran"`
	StatusAkun   string  `json:"status_akun"`
}

type LoginRes struct {
	UserID string `json:"user_id"`
	Nama   string `json:"nama"`
	Email  string `json:"email"`
	Peran  string `json:"peran"`
}

type UserBigResponse struct {
	UserID           string   `json:"user_id"`
	Nama             string   `json:"nama"`
	Email            string   `json:"email"`
	NomorTelepon     *string  `json:"nomor_telepon"`
	Peran            string   `json:"peran"`
	StatusAkun       string   `json:"status_akun"`
	PreferensiUkuran []string `json:"preferensi_ukuran"`
	BrandFavorit     []string `json:"brand_favorit"`
}

type UpdateUserResponse struct {
	UserID           string   `json:"user_id"`
	Nama             string   `json:"nama"`
	NomorTelepon     *string  `json:"nomor_telepon"`
	PreferensiUkuran []string `json:"preferensi_ukuran"`
	BrandFavorit     []string `json:"brand_favorit"`
}
