package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"nora-photobooth-backend/internal/logging"
	"nora-photobooth-backend/internal/models"
)

func ConnectDatabase(cfg *Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
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
		&models.Backdrop{},
		&models.GalleryPhoto{},
		&models.Review{},
		&models.FinanceEntry{},
	)
}
