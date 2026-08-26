// Run once against the target database: go run ./cmd/migrate
// Uses DATABASE_URL from the environment / .env, same as cmd/server.
package main

import (
	"nora-photobooth-backend/internal/config"
	"nora-photobooth-backend/internal/logging"
)

func main() {
	logging.Init()

	cfg := config.Load()

	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		logging.Log.Fatal().Err(err).Msg("failed to connect to database")
	}

	logging.Log.Info().Msg("running migrations")
	if err := config.Migrate(db); err != nil {
		logging.Log.Fatal().Err(err).Msg("migration failed")
	}
	logging.Log.Info().Msg("migrations complete")
}
