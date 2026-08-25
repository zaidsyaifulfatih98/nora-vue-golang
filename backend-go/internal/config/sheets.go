package config

import (
	"context"
	"strings"

	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const FinanceSheetName = "FinanceEntries"

// FinanceRange is the read/write range covering every data column (A:id
// through J:deletedAt), starting at row 2 (row 1 is the header).
const FinanceRange = FinanceSheetName + "!A2:J"

func ConnectSheets(cfg *Config) (*sheets.Service, error) {
	privateKey := strings.ReplaceAll(cfg.GooglePrivateKey, "\\n", "\n")

	jwtConfig := &jwt.Config{
		Email:      cfg.GoogleServiceAccountEmail,
		PrivateKey: []byte(privateKey),
		TokenURL:   google.JWTTokenURL,
		Scopes:     []string{sheets.SpreadsheetsScope},
	}

	client := jwtConfig.Client(context.Background())
	return sheets.NewService(context.Background(), option.WithHTTPClient(client))
}
