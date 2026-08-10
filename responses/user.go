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
