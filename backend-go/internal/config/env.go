package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	DatabaseURL   string
	JWTSecret     string
	CORSWhitelist []string

	CloudinaryCloudName string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string

	GoogleSheetsSpreadsheetID string
	GoogleServiceAccountEmail string
	GooglePrivateKey          string
}

func Load() *Config {
	_ = godotenv.Load()

	whitelist := []string{}
	for _, origin := range strings.Split(os.Getenv("CORS_WHITELIST"), ",") {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			whitelist = append(whitelist, origin)
		}
	}

	return &Config{
		Port:          getEnv("PORT", "8000"),
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		CORSWhitelist: whitelist,

		CloudinaryCloudName: os.Getenv("CLOUDINARY_CLOUD_NAME"),
		CloudinaryAPIKey:    os.Getenv("CLOUDINARY_API_KEY"),
		CloudinaryAPISecret: os.Getenv("CLOUDINARY_API_SECRET"),

		GoogleSheetsSpreadsheetID: os.Getenv("GOOGLE_SHEETS_SPREADSHEET_ID"),
		GoogleServiceAccountEmail: os.Getenv("GOOGLE_SERVICE_ACCOUNT_EMAIL"),
		GooglePrivateKey:          os.Getenv("GOOGLE_PRIVATE_KEY"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
