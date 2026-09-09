package auth

import (
	"errors"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"nora-photobooth-backend/internal/apperror"
	"nora-photobooth-backend/internal/models"
	"nora-photobooth-backend/internal/token"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9-]+$`)

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

func (s *Service) Register(firstName, lastName, email, password string, role models.Role, slug string) (*models.User, error) {
	var existing models.User
	err := s.DB.Where("email = ?", email).First(&existing).Error
	if err == nil {
		return nil, apperror.New("Email sudah terdaftar", 400)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var slugPtr *string
	if role.IsCustomer() {
		normalized := strings.ToLower(strings.TrimSpace(slug))
		if normalized == "" || !slugPattern.MatchString(normalized) {
			return nil, apperror.New("Slug wajib diisi (huruf kecil, angka, dan strip saja)", 400)
		}

		var existingSlug models.User
		err := s.DB.Where("slug = ?", normalized).First(&existingSlug).Error
		if err == nil {
			return nil, apperror.New("Slug sudah digunakan", 409)
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		slugPtr = &normalized
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
		Slug:      slugPtr,
	}
	if err := s.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	user.Password = ""
	return &user, nil
}

// ListCustomers returns every account of the given customer role
// (DIGITAL_PHOTOBOOTH or SOFTWARE_PHOTOBOOTH), for the superadmin onboarding
// dashboards.
func (s *Service) ListCustomers(role models.Role) ([]models.User, error) {
	var users []models.User
	err := s.DB.Where("role = ?", role).Order("created_at desc").Find(&users).Error
	return users, err
}

// UpdateCustomer edits a customer account (DIGITAL_PHOTOBOOTH or
// SOFTWARE_PHOTOBOOTH only — editing an admin account isn't exposed through
// this path). Empty fields are left unchanged; password is only rehashed
// when non-empty.
func (s *Service) UpdateCustomer(id, firstName, lastName, email, slug, password string) (*models.User, error) {
	var user models.User
	err := s.DB.Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) || (err == nil && !user.Role.IsCustomer()) {
		return nil, apperror.New("Pelanggan tidak ditemukan", 404)
	}
	if err != nil {
		return nil, err
	}

	if firstName != "" {
		user.FirstName = firstName
	}
	if lastName != "" {
		user.LastName = lastName
	}

	if email != "" && email != user.Email {
		var existing models.User
		err := s.DB.Where("email = ? AND id <> ?", email, id).First(&existing).Error
		if err == nil {
			return nil, apperror.New("Email sudah terdaftar", 400)
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		user.Email = email
	}

	if slug != "" {
		normalized := strings.ToLower(strings.TrimSpace(slug))
		if !slugPattern.MatchString(normalized) {
			return nil, apperror.New("Slug tidak valid (huruf kecil, angka, dan strip saja)", 400)
		}
		if user.Slug == nil || normalized != *user.Slug {
			var existing models.User
			err := s.DB.Where("slug = ? AND id <> ?", normalized, id).First(&existing).Error
			if err == nil {
				return nil, apperror.New("Slug sudah digunakan", 409)
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
			user.Slug = &normalized
		}
	}

	if password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), 10)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashed)
	}

	if err := s.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	user.Password = ""
	return &user, nil
}

// DeleteCustomer removes a customer account (DIGITAL_PHOTOBOOTH or
// SOFTWARE_PHOTOBOOTH only). Their frames/results/voice messages are left
// in place (still owner-scoped to the now-deleted account id) rather than
// cascade-deleted, consistent with every other soft-delete in this app.
func (s *Service) DeleteCustomer(id string) error {
	var user models.User
	err := s.DB.Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) || (err == nil && !user.Role.IsCustomer()) {
		return apperror.New("Pelanggan tidak ditemukan", 404)
	}
	if err != nil {
		return err
	}
	return s.DB.Delete(&user).Error
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
