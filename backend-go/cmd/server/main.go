package main

import (
	"fmt"

	"google.golang.org/api/sheets/v4"

	"nora-photobooth-backend/internal/config"
	"nora-photobooth-backend/internal/logging"
	"nora-photobooth-backend/internal/router"
)

func main() {
	fmt.Println("boot: process started")

	logging.Init()
	logging.Log.Info().Msg("boot: logger ready, loading config")

	cfg := config.Load()
	logging.Log.Info().Str("port", cfg.Port).Bool("has_database_url", cfg.DatabaseURL != "").Msg("boot: config loaded, connecting to database")

	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		logging.Log.Fatal().Err(err).Msg("failed to connect to database")
	}
	logging.Log.Info().Msg("boot: database connected")

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
