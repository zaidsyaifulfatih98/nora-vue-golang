package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"nora-photobooth-backend/internal/logging"
	"nora-photobooth-backend/internal/models"
)

func ConnectDatabase(cfg *Config) (*gorm.DB, error) {
	// PreferSimpleProtocol disables pgx's server-side prepared statement
	// caching. DATABASE_URL points at Neon's pooled (PgBouncer
	// transaction-mode) endpoint, which multiplexes client sessions across
	// physical connections — a cached plan from one session can be reused
	// against a connection that has since seen a schema change, surfacing as
	// "cached plan must not change result type" (SQLSTATE 0A000) the next
	// time a table's columns change.
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  cfg.DatabaseURL,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	logging.Log.Info().Msg("boot: gorm.Open succeeded, pinging database")

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}
	logging.Log.Info().Msg("boot: ping succeeded")

	return db, nil
}

// Migrate runs schema migrations. Run it once via `go run ./cmd/migrate`
// against the target database instead of on every server boot: DDL against
// Neon's pooled (PgBouncer transaction-mode) endpoint can hang indefinitely,
// and running it on every serverless cold start risks lock contention
// between concurrent invocations.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Package{},
		&models.Feature{},
		&models.FrameTemplate{},
		&models.PhotoboothFrame{},
		&models.VoiceMessage{},
		&models.Backdrop{},
		&models.GalleryPhoto{},
		&models.Review{},
		&models.FinanceEntry{},
	)
}
