package main

import (
	"fmt"

	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"nora-photobooth-backend/internal/config"
	"nora-photobooth-backend/internal/logging"
	"nora-photobooth-backend/internal/models"
)

func main() {
	logging.Init()
	cfg := config.Load()

	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		logging.Log.Fatal().Err(err).Msg("failed to connect to database")
	}

	seedAdmin(db)
	seedPackages(db)
	seedReviews(db)

	fmt.Println("Seed complete.")
}

func seedAdmin(db *gorm.DB) {
	email := "admin@noraphotobooth.id"

	var existing models.User
	if err := db.Where("email = ?", email).First(&existing).Error; err == nil {
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte("admin12345"), 10)
	if err != nil {
		logging.Log.Fatal().Err(err).Msg("failed to hash seed password")
	}

	admin := models.User{
		FirstName: "Super",
		LastName:  "Admin",
		Email:     email,
		Password:  string(hashed),
		Role:      models.RoleSuperAdmin,
	}
	if err := db.Create(&admin).Error; err != nil {
		logging.Log.Fatal().Err(err).Msg("failed to seed admin user")
	}
}

func seedPackages(db *gorm.DB) {
	var count int64
	db.Model(&models.Package{}).Count(&count)
	if count > 0 {
		return
	}

	pkgs := []models.Package{
		{Name: "Silver", Price: decimal.NewFromInt(500000), Duration: "2 jam", Description: "Paket dasar untuk acara kecil", Features: []string{"1 backdrop", "Unlimited foto", "Cetak 50 lembar"}, Order: 1},
		{Name: "Gold", Price: decimal.NewFromInt(1000000), Duration: "3 jam", Description: "Paket populer untuk acara menengah", Features: []string{"2 backdrop", "Unlimited foto", "Cetak 100 lembar", "Props"}, IsPopular: true, Order: 2},
		{Name: "Platinum", Price: decimal.NewFromInt(1500000), Duration: "4 jam", Description: "Paket lengkap untuk acara besar", Features: []string{"3 backdrop", "Unlimited foto", "Cetak unlimited", "Props", "Album cetak"}, Order: 3},
	}
	if err := db.Create(&pkgs).Error; err != nil {
		logging.Log.Fatal().Err(err).Msg("failed to seed packages")
	}
}

func seedReviews(db *gorm.DB) {
	var count int64
	db.Model(&models.Review{}).Count(&count)
	if count > 0 {
		return
	}

	revs := []models.Review{
		{Name: "Andi & Sarah", EventLabel: "Pernikahan", Quote: "Layanan sangat memuaskan, tamu-tamu senang!", Rating: 5, Order: 1},
		{Name: "PT Maju Jaya", EventLabel: "Gathering Kantor", Quote: "Fotobooth-nya seru banget, jadi hiburan tambahan.", Rating: 5, Order: 2},
		{Name: "Rina", EventLabel: "Ulang Tahun", Quote: "Hasil fotonya bagus dan cepat, recommended!", Rating: 4, Order: 3},
	}
	if err := db.Create(&revs).Error; err != nil {
		logging.Log.Fatal().Err(err).Msg("failed to seed reviews")
	}
}
