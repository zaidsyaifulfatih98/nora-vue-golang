package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"nora-photobooth-backend/internal/models"
)

func ConnectDatabase(cfg *Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Package{},
		&models.Feature{},
		&models.FrameTemplate{},
		&models.Backdrop{},
		&models.GalleryPhoto{},
		&models.Review{},
		&models.FinanceEntry{},
	); err != nil {
		return nil, err
	}

	return db, nil
}
