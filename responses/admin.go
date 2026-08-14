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

type AdminProductListItemResponse struct {
	ProductID       string  `json:"product_id"`
	NamaProduk      string  `json:"nama_produk"`
	SellerID        string  `json:"seller_id"`
	Harga           float64 `json:"harga"`
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
