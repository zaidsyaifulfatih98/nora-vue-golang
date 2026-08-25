package finance

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"nora-photobooth-backend/internal/apperror"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(c *gin.Context) {
	entries, err := h.service.GetEntries(c.Query("from"), c.Query("to"))
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Finance entries fetched", "data": entries})
}

func (h *Handler) Summary(c *gin.Context) {
	summary, err := h.service.GetSummary(c.Query("from"), c.Query("to"))
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Finance summary fetched", "data": summary})
}

type createRequest struct {
	Type        string  `json:"type" binding:"required,oneof=INCOME EXPENSE"`
	Category    string  `json:"category" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,min=0"`
	Description string  `json:"description"`
	Date        string  `json:"date" binding:"required"`
}

func (h *Handler) Create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.New("Data finance tidak valid", 400))
		c.Abort()
		return
	}

	userID := c.GetString("userID")

	entry, err := h.service.CreateEntry(CreateInput{
		Type:        req.Type,
		Category:    req.Category,
		Amount:      req.Amount,
		Description: req.Description,
		Date:        req.Date,
		CreatedByID: userID,
	})
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Finance entry created", "data": entry})
}

type updateRequest struct {
	Type        *string  `json:"type" binding:"omitempty,oneof=INCOME EXPENSE"`
	Category    *string  `json:"category"`
	Amount      *float64 `json:"amount"`
	Description *string  `json:"description"`
	Date        *string  `json:"date"`
}

func (h *Handler) Update(c *gin.Context) {
	var req updateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.New("Data finance tidak valid", 400))
		c.Abort()
		return
	}

	entry, err := h.service.UpdateEntry(c.Param("id"), UpdateInput{
		Type:        req.Type,
		Category:    req.Category,
		Amount:      req.Amount,
		Description: req.Description,
		Date:        req.Date,
	})
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Finance entry updated", "data": entry})
}

func (h *Handler) Delete(c *gin.Context) {
	if err := h.service.DeleteEntry(c.Param("id")); err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Finance entry deleted", "data": gin.H{}})
}
