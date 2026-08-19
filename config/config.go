package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DBName      string
	DBUser      string
	DBPassword  string
	DBHost      string
	DBPort      string
	Port        string
	PublicURL   string
	PaymentMode string
	CORSOrigins string

	TripayAPIKey       string
	TripayPrivateKey   string
	TripayMerchantCode string
	TripaySandbox      bool

	MidtransServerKey string
	MidtransClientKey string
	MidtransSandbox   bool

	BiteshipAPIKey string

	ShippingPollMinutes int
	AutoCompleteDays    int
}

var AppConfig *Config

func LoadEnv() {

	if err := godotenv.Load(); err != nil {
		log.Println("Error Not Found file .env !⚠️")
	}

	AppConfig = &Config{
		Port:        os.Getenv("PORT"),
		DBName:      os.Getenv("DB_NAME"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		DBUser:      os.Getenv("DB_USER"),
		DBHost:      os.Getenv("DB_HOST"),
		DBPort:      os.Getenv("DB_PORT"),
		PublicURL:   os.Getenv("PUBLIC_URL"),
		CORSOrigins: os.Getenv("CORS_ORIGINS"),

		PaymentMode: normalizePaymentMode(os.Getenv("PAYMENT_MODE")),

		TripayAPIKey:       os.Getenv("TRIPAY_API_KEY"),
		TripayPrivateKey:   os.Getenv("TRIPAY_PRIVATE_KEY"),
		TripayMerchantCode: os.Getenv("TRIPAY_MERCHANT_CODE"),
		TripaySandbox:      os.Getenv("TRIPAY_IS_SANDBOX") == "true",

		MidtransServerKey: os.Getenv("MIDTRANS_SERVER_KEY"),
		MidtransClientKey: os.Getenv("MIDTRANS_CLIENT_KEY"),
		MidtransSandbox:   os.Getenv("MIDTRANS_IS_SANDBOX") == "true",

		BiteshipAPIKey: os.Getenv("BITESHIP_API_KEY"),

		ShippingPollMinutes: parseIntEnv("SHIPPING_POLL_MINUTES", 30),
		AutoCompleteDays:    parseIntEnv("AUTO_COMPLETE_DAYS", 7),
	}
}

func parseIntEnv(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 {
		return fallback
	}
	return n
}

func normalizePaymentMode(mode string) string {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "mock":
		return "mock"
	case "midtrans":
		return "midtrans"
	default:
		return "tripay"
	}
}
