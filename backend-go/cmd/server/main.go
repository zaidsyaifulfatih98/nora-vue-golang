package main

import (
	"google.golang.org/api/sheets/v4"

	"nora-photobooth-backend/internal/config"
	"nora-photobooth-backend/internal/logging"
	"nora-photobooth-backend/internal/router"
)

func main() {
	logging.Init()

	cfg := config.Load()

	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		logging.Log.Fatal().Err(err).Msg("failed to connect to database")
	}

	var sheetsSvc *sheets.Service
	if cfg.GoogleServiceAccountEmail != "" && cfg.GooglePrivateKey != "" {
		sheetsSvc, err = config.ConnectSheets(cfg)
		if err != nil {
			logging.Log.Warn().Err(err).Msg("failed to connect to Google Sheets, finance module disabled")
			sheetsSvc = nil
		}
	} else {
		logging.Log.Warn().Msg("Google Sheets credentials not set, finance module disabled")
	}

	r, err := router.New(cfg, db, sheetsSvc)
	if err != nil {
		logging.Log.Fatal().Err(err).Msg("failed to build router")
	}

	logging.Log.Info().Str("port", cfg.Port).Msg("server starting")
	if err := r.Run(":" + cfg.Port); err != nil {
		logging.Log.Fatal().Err(err).Msg("server stopped")
	}
}
