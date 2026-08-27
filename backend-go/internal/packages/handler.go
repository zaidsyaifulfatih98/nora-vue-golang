package packages

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"nora-photobooth-backend/internal/apperror"
	"nora-photobooth-backend/internal/crud"
	"nora-photobooth-backend/internal/models"
)

type Handler struct {
	repo *crud.Repository[models.Package]
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{repo: crud.NewRepository[models.Package](db, "Package", "is_active")}
}

func (h *Handler) List() gin.HandlerFunc   { return crud.ListHandler(h.repo) }
func (h *Handler) Delete() gin.HandlerFunc { return crud.DeleteHandler(h.repo) }

type createRequest struct {
	Name        string   `json:"name" binding:"required,max=50"`
	Price       float64  `json:"price" binding:"required,min=0"`
	Duration    string   `json:"duration" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Features    []string `json:"features"`
	IsPopular   bool     `json:"isPopular"`
	IsActive    *bool    `json:"isActive"`
	Order       int      `json:"order"`
}

func (h *Handler) Create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.New("Data paket tidak valid", 400))
		c.Abort()
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	item := models.Package{
		Name:        req.Name,
		Price:       decimal.NewFromFloat(req.Price),
		Duration:    req.Duration,
		Description: req.Description,
		Features:    req.Features,
		IsPopular:   req.IsPopular,
		IsActive:    isActive,
		Order:       req.Order,
	}

	if err := h.repo.Create(&item); err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Package created", "data": item})
}

type updateRequest struct {
	Name        *string   `json:"name"`
	Price       *float64  `json:"price"`
	Duration    *string   `json:"duration"`
	Description *string   `json:"description"`
	Features    *[]string `json:"features"`
	IsPopular   *bool     `json:"isPopular"`
	IsActive    *bool     `json:"isActive"`
	Order       *int      `json:"order"`
}

func (h *Handler) Update(c *gin.Context) {
	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.New("Data paket tidak valid", 400))
		c.Abort()
		return
	}

	updates := map[string]any{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Price != nil {
		updates["price"] = decimal.NewFromFloat(*req.Price)
	}
	if req.Duration != nil {
		updates["duration"] = *req.Duration
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Features != nil {
		updates["features"] = pq.StringArray(*req.Features)
	}
	if req.IsPopular != nil {
		updates["is_popular"] = *req.IsPopular
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if req.Order != nil {
		updates["order"] = *req.Order
	}

	item, err := h.repo.Update(c.Param("id"), updates)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Package updated", "data": item})
}
