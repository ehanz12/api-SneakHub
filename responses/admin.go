package responses

import "time"

type AdminUserListItemResponse struct {
	UserID     string `json:"user_id"`
	Nama       string `json:"nama"`
	Email      string `json:"email"`
	Peran      string `json:"peran"`
	StatusAkun string `json:"status_akun"`
}

type AdminUserListDataResponse struct {
	Items      []AdminUserListItemResponse `json:"items"`
	Pagination PaginationResponse          `json:"pagination"`
}

type AdminUserStatusResponse struct {
	UserID     string `json:"user_id"`
	StatusAkun string `json:"status_akun"`
}

type AdminSellerListItemResponse struct {
	SellerID         string    `json:"seller_id"`
	UserID           string    `json:"user_id"`
	NamaToko         string    `json:"nama_toko"`
	StatusVerifikasi string    `json:"status_verifikasi"`
	NamaUser         string    `json:"nama_user"`
	Email            string    `json:"email"`
	CreatedAt        time.Time `json:"created_at"`
}

type AdminSellerListDataResponse struct {
	Items      []AdminSellerListItemResponse `json:"items"`
	Pagination PaginationResponse            `json:"pagination"`
}

type AdminSellerVerificationResponse struct {
	SellerID         string `json:"seller_id"`
	StatusVerifikasi string `json:"status_verifikasi"`
}

type AdminUserRoleResponse struct {
	UserID string `json:"user_id"`
	Peran  string `json:"peran"`
}

type AdminProductListItemResponse struct {
	ProductID       string  `json:"product_id"`
	NamaProduk      string  `json:"nama_produk"`
	SellerID        string  `json:"seller_id"`
	Harga           float64 `json:"harga"`
	ImageURL        string  `json:"image_url"`
	StatusPublikasi string  `json:"status_publikasi"`
}

type AdminProductListDataResponse struct {
	Items      []AdminProductListItemResponse `json:"items"`
	Pagination PaginationResponse             `json:"pagination"`
}

type AdminProductStatusResponse struct {
	ProductID       string `json:"product_id"`
	StatusPublikasi string `json:"status_publikasi"`
}

type AdminOrderListItemResponse struct {
	OrderID         string    `json:"order_id"`
	CustomerID      string    `json:"customer_id"`
	SellerID        string    `json:"seller_id"`
	StatusOrder     string    `json:"status_order"`
	TotalPembayaran float64   `json:"total_pembayaran"`
	CreatedAt       time.Time `json:"created_at"`
}

type AdminOrderListDataResponse struct {
	Items      []AdminOrderListItemResponse `json:"items"`
	Pagination PaginationResponse           `json:"pagination"`
}

type AdminReportResponse struct {
	Period        string  `json:"period"`
	TotalUsers    int64   `json:"total_users"`
	TotalSellers  int64   `json:"total_sellers"`
	TotalProducts int64   `json:"total_products"`
	TotalOrders   int64   `json:"total_orders"`
	TotalRevenue  float64 `json:"total_revenue"`
}
