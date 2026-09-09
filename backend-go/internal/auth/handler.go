package auth

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"nora-photobooth-backend/internal/apperror"
	"nora-photobooth-backend/internal/middleware"
	"nora-photobooth-backend/internal/models"
)

const cookieName = "token"
const cookieMaxAge = 60 * 60 * 24 // 1 day

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, apperror.New(firstBindingError(err, "Email dan password wajib diisi"), 400))
		return
	}

	user, token, err := h.Service.Login(req.Email, req.Password)
	if err != nil {
		middleware.Fail(c, err)
		return
	}

	isProd := os.Getenv("APP_ENV") == "production"
	c.SetSameSite(sameSiteFor(isProd))
	c.SetCookie(cookieName, token, cookieMaxAge, "/", "", isProd, true)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Login berhasil", "data": user})
}

func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie(cookieName, "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Logout berhasil", "data": gin.H{}})
}

func (h *Handler) Me(c *gin.Context) {
	userID := c.GetString("userID")
	user, err := h.Service.CurrentUser(userID)
	if err != nil {
		middleware.Fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "User ditemukan", "data": user})
}

type registerRequest struct {
	FirstName string      `json:"firstName" binding:"required,min=5,max=25"`
	LastName  string      `json:"lastName" binding:"required,min=5,max=25"`
	Email     string      `json:"email" binding:"required,email"`
	Password  string      `json:"password" binding:"required"`
	Role      models.Role `json:"role" binding:"required,oneof=SUPER_ADMIN ADMIN DIGITAL_PHOTOBOOTH SOFTWARE_PHOTOBOOTH"`
	Slug      string      `json:"slug"`
}

func (h *Handler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, apperror.New(firstBindingError(err, "Data registrasi tidak valid"), 400))
		return
	}

	user, err := h.Service.Register(req.FirstName, req.LastName, req.Email, req.Password, req.Role, req.Slug)
	if err != nil {
		middleware.Fail(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Registrasi berhasil", "data": user})
}

// ListCustomers is superadmin-only — the dashboard's list of onboarded
// customer accounts, filtered by ?role= (DIGITAL_PHOTOBOOTH by default).
func (h *Handler) ListCustomers(c *gin.Context) {
	role := models.Role(c.DefaultQuery("role", string(models.RoleDigitalPhotobooth)))
	if !role.IsCustomer() {
		middleware.Fail(c, apperror.New("Role tidak valid", 400))
		return
	}

	users, err := h.Service.ListCustomers(role)
	if err != nil {
		middleware.Fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Customers fetched", "data": users})
}

type updateCustomerRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email" binding:"omitempty,email"`
	Slug      string `json:"slug"`
	Password  string `json:"password"`
}

// UpdateCustomer is superadmin-only — edits a customer account's profile,
// slug, or (optionally) resets their password.
func (h *Handler) UpdateCustomer(c *gin.Context) {
	var req updateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, apperror.New(firstBindingError(err, "Data tidak valid"), 400))
		return
	}

	user, err := h.Service.UpdateCustomer(c.Param("id"), req.FirstName, req.LastName, req.Email, req.Slug, req.Password)
	if err != nil {
		middleware.Fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Pelanggan diperbarui", "data": user})
}

// DeleteCustomer is superadmin-only — removes a customer account.
func (h *Handler) DeleteCustomer(c *gin.Context) {
	if err := h.Service.DeleteCustomer(c.Param("id")); err != nil {
		middleware.Fail(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Pelanggan dihapus", "data": gin.H{}})
}
