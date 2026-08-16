package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// struct untuk config database
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
}

// variable dari struct
var AppConfig *Config

func LoadEnv() {
	// cek env kalo ga ada kasih error
	if err := godotenv.Load(); err != nil {
		log.Println("Error Not Found file .env !⚠️")
	}

	// instalasi untuk config
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
	}
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
