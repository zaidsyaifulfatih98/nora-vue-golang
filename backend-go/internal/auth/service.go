package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"nora-photobooth-backend/internal/apperror"
	"nora-photobooth-backend/internal/models"
	"nora-photobooth-backend/internal/token"
)

type Service struct {
	DB        *gorm.DB
	JWTSecret string
}

func NewService(db *gorm.DB, jwtSecret string) *Service {
	return &Service{DB: db, JWTSecret: jwtSecret}
}

// Login validates credentials and returns a sanitized user + signed token.
func (s *Service) Login(email, password string) (*models.User, string, error) {
	var user models.User
	err := s.DB.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, "", apperror.New("Email atau password salah", 401)
	}
	if err != nil {
		return nil, "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", apperror.New("Email atau password salah", 401)
	}

	tok, err := token.Create(user.ID, user.Role, s.JWTSecret)
	if err != nil {
		return nil, "", err
	}

	user.Password = ""
	return &user, tok, nil
}

func (s *Service) Register(firstName, lastName, email, password string, role models.Role) (*models.User, error) {
	var existing models.User
	err := s.DB.Where("email = ?", email).First(&existing).Error
	if err == nil {
		return nil, apperror.New("Email sudah terdaftar", 400)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}

	user := models.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  string(hashed),
		Role:      role,
	}
	if err := s.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	user.Password = ""
	return &user, nil
}

func (s *Service) CurrentUser(id string) (*models.User, error) {
	var user models.User
	err := s.DB.Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.New("User tidak ditemukan", 404)
	}
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return &user, nil
}
