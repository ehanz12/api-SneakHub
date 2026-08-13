package services

import (
	"bytes"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ehanz12/api-SneakHub/config"
)

const (
	midtransProdBase   = "https://app.midtrans.com/snap/v1/transactions"
	midtransSandboxURL = "https://app.sandbox.midtrans.com/snap/v1/transactions"
)

// MidtransProvider menggunakan gateway pembayaran Midtrans (Snap API).
type MidtransProvider struct{}

func (MidtransProvider) CreatePayment(orderID, metode string, amount int64, items []TripayOrderItem, customerName, customerEmail, customerPhone string) (string, string, error) {
	return createMidtransPayment(orderID, metode, amount, items, customerName, customerEmail, customerPhone)
}

func (MidtransProvider) VerifySignature(rawBody []byte, signature string) bool {
	return VerifyMidtransSignature(rawBody, signature)
}

type midtransItemDetail struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Price    int64  `json:"price"`
	Quantity int32  `json:"quantity"`
}

type midtransSnapRequest struct {
	TransactionDetails struct {
		OrderID     string `json:"order_id"`
		GrossAmount int64  `json:"gross_amount"`
	} `json:"transaction_details"`
	ItemDetails     []midtransItemDetail `json:"item_details"`
	CustomerDetails struct {
		FirstName string `json:"first_name,omitempty"`
		Email     string `json:"email,omitempty"`
		Phone     string `json:"phone,omitempty"`
	} `json:"customer_details,omitempty"`
	EnabledPayments []string `json:"enabled_payments,omitempty"`
	Expiry          struct {
		StartTime string `json:"start_time"`
		Unit      string `json:"unit"`
		Duration  int    `json:"duration"`
	} `json:"expiry,omitempty"`
	NotificationURL string `json:"notification_url,omitempty"`
}

type midtransSnapResponse struct {
	StatusCode    string `json:"status_code"`
	StatusMessage string `json:"status_message"`
	Token         string `json:"token"`
	RedirectURL   string `json:"redirect_url"`
}

// MidtransNotificationPayload mewakili body callback dari Midtrans.
// GrossAmount sengaja berupa string agar signature dihitung dari nilai
// persis seperti yang dikirim Midtrans.
type MidtransNotificationPayload struct {
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
	OrderID           string `json:"order_id"`
	StatusCode        string `json:"status_code"`
	GrossAmount       string `json:"gross_amount"`
	SignatureKey      string `json:"signature_key"`
	TransactionID     string `json:"transaction_id"`
	PaymentType       string `json:"payment_type"`
}

func midtransBaseURL() string {
	if config.AppConfig.MidtransSandbox {
		return midtransSandboxURL
	}
	return midtransProdBase
}

func ensureMidtransConfigured() error {
	if strings.TrimSpace(config.AppConfig.MidtransServerKey) == "" {
		return errors.New("MIDTRANS_SERVER_KEY belum diisi di .env")
	}
	return nil
}

// midtransEnabledPayments memetakan metode pembayaran yang dipakai checkout
// (kode channel Tripay) ke daftar enabled_payments Midtrans.
func midtransEnabledPayments(metode string) []string {
	switch strings.ToUpper(strings.TrimSpace(metode)) {
	case "QRIS2":
		return []string{"qris"}
	case "BCAVA", "BANK_TRANSFER", "VA":
		return []string{"bank_transfer", "bca_va"}
	default:
		return []string{"qris"}
	}
}

// createMidtransPayment membuat transaksi Snap (open payment) dan
// mengembalikan order_id sebagai reference + redirect_url Snap.
func createMidtransPayment(orderID, metode string, amount int64, items []TripayOrderItem, customerName, customerEmail, customerPhone string) (string, string, error) {
	if err := ensureMidtransConfigured(); err != nil {
		return "", "", err
	}

	itemDetails := make([]midtransItemDetail, 0, len(items)+1)
	itemSum := int64(0)
	for _, item := range items {
		itemSum += item.Price * int64(item.Quantity)
		itemDetails = append(itemDetails, midtransItemDetail{
			ID:       item.Sku,
			Name:     item.Name,
			Price:    item.Price,
			Quantity: item.Quantity,
		})
	}

	// Midtrans memvalidasi gross_amount = total item_details, sehingga
	// ongkos kirim ditambahkan sebagai item terpisah.
	shipping := amount - itemSum
	if shipping > 0 {
		itemDetails = append(itemDetails, midtransItemDetail{
			ID:       "SHIPPING",
			Name:     "Biaya Pengiriman",
			Price:    shipping,
			Quantity: 1,
		})
	}

	req := midtransSnapRequest{}
	req.TransactionDetails.OrderID = orderID
	req.TransactionDetails.GrossAmount = amount
	req.ItemDetails = itemDetails
	req.CustomerDetails.FirstName = customerName
	req.CustomerDetails.Email = customerEmail
	req.CustomerDetails.Phone = customerPhone
	req.EnabledPayments = midtransEnabledPayments(metode)
	req.Expiry.StartTime = time.Now().UTC().Format("2006-01-02 15:04:05 -0700")
	req.Expiry.Unit = "hours"
	req.Expiry.Duration = 24
	req.NotificationURL = config.AppConfig.PublicURL + "/api/payments/midtrans/notification"

	body, err := json.Marshal(req)
	if err != nil {
		return "", "", errors.New("gagal membuat payload transaksi")
	}

	httpReq, err := http.NewRequest(http.MethodPost, midtransBaseURL(), bytes.NewReader(body))
	if err != nil {
		return "", "", errors.New("gagal membuat request transaksi")
	}
	auth := base64.StdEncoding.EncodeToString([]byte(strings.TrimSpace(config.AppConfig.MidtransServerKey) + ":"))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", "Basic "+auth)

	client := &http.Client{Timeout: 30 * time.Second}
	res, err := client.Do(httpReq)
	if err != nil {
		return "", "", fmt.Errorf("gagal menghubungi server Midtrans: %s", err.Error())
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", "", errors.New("gagal membaca respons Midtrans")
	}

	var payload midtransSnapResponse
	if err := json.Unmarshal(resBody, &payload); err != nil {
		return "", "", errors.New("respons Midtrans tidak valid")
	}
	if payload.StatusCode != "201" {
		return "", "", fmt.Errorf("Midtrans: %s", payload.StatusMessage)
	}
	if strings.TrimSpace(payload.RedirectURL) == "" {
		return "", "", errors.New("Midtrans tidak mengembalikan redirect_url")
	}

	return orderID, payload.RedirectURL, nil
}

// VerifyMidtransSignature memverifikasi signature_key dari callback Midtrans.
// Signature = SHA512(order_id + status_code + gross_amount + server_key).
func VerifyMidtransSignature(rawBody []byte, signature string) bool {
	if strings.TrimSpace(config.AppConfig.MidtransServerKey) == "" {
		return false
	}

	var payload MidtransNotificationPayload
	if err := json.Unmarshal(rawBody, &payload); err != nil {
		return false
	}
	if payload.OrderID == "" || payload.StatusCode == "" || payload.GrossAmount == "" {
		return false
	}

	raw := payload.OrderID + payload.StatusCode + payload.GrossAmount + strings.TrimSpace(config.AppConfig.MidtransServerKey)
	sum := sha512.Sum512([]byte(raw))
	return strings.EqualFold(hex.EncodeToString(sum[:]), strings.TrimSpace(signature))
}
