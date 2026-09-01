package reviews

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"nora-photobooth-backend/internal/apperror"
	"nora-photobooth-backend/internal/crud"
	"nora-photobooth-backend/internal/models"
)

type Handler struct {
	repo *crud.Repository[models.Review]
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{repo: crud.NewRepository[models.Review](db, "Review", "is_published")}
}

func (h *Handler) List() gin.HandlerFunc   { return crud.ListHandler(h.repo) }
func (h *Handler) Delete() gin.HandlerFunc { return crud.DeleteHandler(h.repo) }

type createRequest struct {
	Name         string  `json:"name" binding:"required"`
	EventLabel   string  `json:"eventLabel" binding:"required"`
	Quote        string  `json:"quote" binding:"required"`
	EventLabelEn *string `json:"eventLabelEn"`
	QuoteEn      *string `json:"quoteEn"`
	Rating       *int    `json:"rating" binding:"omitempty,min=1,max=5"`
	IsPublished  *bool   `json:"isPublished"`
	Order        int     `json:"order"`
}

func (h *Handler) Create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.New("Data review tidak valid", 400))
		c.Abort()
		return
	}

	rating := 5
	if req.Rating != nil {
		rating = *req.Rating
	}
	published := true
	if req.IsPublished != nil {
		published = *req.IsPublished
	}

	item := models.Review{
		Name:         req.Name,
		EventLabel:   req.EventLabel,
		Quote:        req.Quote,
		EventLabelEn: req.EventLabelEn,
		QuoteEn:      req.QuoteEn,
		Rating:       rating,
		IsPublished:  published,
		Order:        req.Order,
	}

	if err := h.repo.Create(&item); err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Review created", "data": item})
}

type updateRequest struct {
	Name         *string `json:"name"`
	EventLabel   *string `json:"eventLabel"`
	Quote        *string `json:"quote"`
	EventLabelEn *string `json:"eventLabelEn"`
	QuoteEn      *string `json:"quoteEn"`
	Rating       *int    `json:"rating" binding:"omitempty,min=1,max=5"`
	IsPublished  *bool   `json:"isPublished"`
	Order        *int    `json:"order"`
}

func (h *Handler) Update(c *gin.Context) {
	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.New("Data review tidak valid", 400))
		c.Abort()
		return
	}

	updates := map[string]any{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.EventLabel != nil {
		updates["event_label"] = *req.EventLabel
	}
	if req.Quote != nil {
		updates["quote"] = *req.Quote
	}
	if req.EventLabelEn != nil {
		updates["event_label_en"] = *req.EventLabelEn
	}
	if req.QuoteEn != nil {
		updates["quote_en"] = *req.QuoteEn
	}
	if req.Rating != nil {
		updates["rating"] = *req.Rating
	}
	if req.IsPublished != nil {
		updates["is_published"] = *req.IsPublished
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

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Review updated", "data": item})
}
