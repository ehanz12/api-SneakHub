package services

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ehanz12/api-SneakHub/config"
)

const (
	tripayProdBase = "https://payment.tripay.co.id/api"
	tripaySandbox  = "https://payment.tripay.co.id/api-sandbox"
)

type TripayOrderItem struct {
	Sku      string `json:"sku"`
	Name     string `json:"name"`
	Price    int64  `json:"price"`
	Quantity int32  `json:"quantity"`
}

type PaymentProvider interface {
	CreatePayment(orderID, metode string, amount int64, items []TripayOrderItem, customerName, customerEmail, customerPhone string) (reference, paymentURL string, err error)
	VerifySignature(rawBody []byte, signature string) bool
}

func GetPaymentProvider() PaymentProvider {
	switch config.AppConfig.PaymentMode {
	case "mock":
		return MockProvider{}
	case "midtrans":
		return MidtransProvider{}
	default:
		return TripayProvider{}
	}
}

// TripayProvider menggunakan gateway pembayaran Tripay.
type TripayProvider struct{}

func (TripayProvider) CreatePayment(orderID, metode string, amount int64, items []TripayOrderItem, customerName, customerEmail, customerPhone string) (string, string, error) {
	return createTripayPayment(orderID, metode, amount, items, customerName, customerEmail, customerPhone)
}

func (TripayProvider) VerifySignature(rawBody []byte, signature string) bool {
	return VerifyTripaySignature(rawBody, signature)
}

// MockProvider mensimulasikan gateway pembayaran untuk keperluan testing
// tanpa perlu akun Tripay. Tidak boleh aktif di production.
type MockProvider struct{}

func (MockProvider) CreatePayment(orderID, metode string, amount int64, items []TripayOrderItem, customerName, customerEmail, customerPhone string) (string, string, error) {
	base := strings.TrimSpace(config.AppConfig.PublicURL)
	if base == "" {
		base = "http://localhost:" + config.AppConfig.Port
	}
	reference := "MOCK-" + orderID
	paymentURL := base + "/api/mock-pay/" + orderID + "?amount=" + strconv.FormatInt(amount, 10) + "&name=" + url.QueryEscape(customerName)
	return reference, paymentURL, nil
}

func (MockProvider) VerifySignature(rawBody []byte, signature string) bool {
	return true
}

type tripayPaymentRequest struct {
	Method        string            `json:"method"`
	MerchantRef   string            `json:"merchant_ref"`
	Amount        int64             `json:"amount"`
	CustomerName  string            `json:"customer_name"`
	CustomerEmail string            `json:"customer_email"`
	CustomerPhone string            `json:"customer_phone,omitempty"`
	OrderItems    []TripayOrderItem `json:"order_items"`
	CallbackURL   string            `json:"callback_url,omitempty"`
	ReturnURL     string            `json:"return_url,omitempty"`
	ExpiredTime   int64             `json:"expired_time,omitempty"`
	Signature     string            `json:"signature"`
}

type tripayResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Reference   string `json:"reference"`
		CheckoutURL string `json:"checkout_url"`
	} `json:"data"`
}

func tripayBaseURL() string {
	if config.AppConfig.TripaySandbox {
		return tripaySandbox
	}
	return tripayProdBase
}

func ensureTripayConfigured() error {
	if strings.TrimSpace(config.AppConfig.TripayAPIKey) == "" {
		return errors.New("TRIPAY_API_KEY belum diisi di .env")
	}
	if strings.TrimSpace(config.AppConfig.TripayPrivateKey) == "" {
		return errors.New("TRIPAY_PRIVATE_KEY belum diisi di .env")
	}
	if strings.TrimSpace(config.AppConfig.TripayMerchantCode) == "" {
		return errors.New("TRIPAY_MERCHANT_CODE belum diisi di .env")
	}
	return nil
}

// tripaySignature membuat signature HMAC-SHA256 untuk request create transaksi.
// Signature = HMAC-SHA256(merchant_code + merchant_ref + amount, private_key).
func tripaySignature(merchantRef string, amount int64) string {
	raw := config.AppConfig.TripayMerchantCode + merchantRef + strconv.FormatInt(amount, 10)
	mac := hmac.New(sha256.New, []byte(config.AppConfig.TripayPrivateKey))
	mac.Write([]byte(raw))
	return hex.EncodeToString(mac.Sum(nil))
}

// createTripayPayment membuat transaksi Tripay (closed payment) dan
// mengembalikan reference + checkout URL.
func createTripayPayment(orderID, method string, amount int64, items []TripayOrderItem, customerName, customerEmail, customerPhone string) (string, string, error) {
	if err := ensureTripayConfigured(); err != nil {
		return "", "", err
	}

	req := tripayPaymentRequest{
		Method:        method,
		MerchantRef:   orderID,
		Amount:        amount,
		CustomerName:  customerName,
		CustomerEmail: customerEmail,
		CustomerPhone: customerPhone,
		OrderItems:    items,
		CallbackURL:   config.AppConfig.PublicURL + "/api/payments/notification",
		ExpiredTime:   time.Now().Add(24 * time.Hour).Unix(),
		Signature:     tripaySignature(orderID, amount),
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", "", errors.New("gagal membuat payload transaksi")
	}

	httpReq, err := http.NewRequest(http.MethodPost, tripayBaseURL()+"/transaction/create", bytes.NewReader(body))
	if err != nil {
		return "", "", errors.New("gagal membuat request transaksi")
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+strings.TrimSpace(config.AppConfig.TripayAPIKey))

	client := &http.Client{Timeout: 30 * time.Second}
	res, err := client.Do(httpReq)
	if err != nil {
		return "", "", fmt.Errorf("gagal menghubungi server Tripay: %s", err.Error())
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", "", errors.New("gagal membaca respons Tripay")
	}

	var payload tripayResponse
	if err := json.Unmarshal(resBody, &payload); err != nil {
		return "", "", errors.New("respons Tripay tidak valid")
	}
	if !payload.Success {
		return "", "", fmt.Errorf("Tripay: %s", payload.Message)
	}
	if payload.Data.Reference == "" {
		return "", "", errors.New("Tripay tidak mengembalikan reference transaksi")
	}

	return payload.Data.Reference, payload.Data.CheckoutURL, nil
}

// VerifyTripaySignature memverifikasi X-Callback-Signature dari callback Tripay.
// Signature = HMAC-SHA256(raw JSON body, private_key).
func VerifyTripaySignature(rawBody []byte, signature string) bool {
	if strings.TrimSpace(config.AppConfig.TripayPrivateKey) == "" {
		return false
	}
	mac := hmac.New(sha256.New, []byte(config.AppConfig.TripayPrivateKey))
	mac.Write(rawBody)
	return hmac.Equal(mac.Sum(nil), []byte(strings.TrimSpace(signature)))
}
