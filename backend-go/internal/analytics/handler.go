package analytics

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"nora-photobooth-backend/internal/apperror"
	"nora-photobooth-backend/internal/models"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

type trackRequest struct {
	EventType string `json:"eventType" binding:"required,oneof=page_view cta_click"`
	Path      string `json:"path" binding:"required,max=255"`
	CtaLabel  string `json:"ctaLabel" binding:"max=100"`
}

// Track is public — called anonymously from the landing page, digital
// photobooth page, etc. No auth, so it only ever writes, never reads back
// anything a caller sent.
func (h *Handler) Track(c *gin.Context) {
	var req trackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperror.New("Data analytics tidak valid", 400))
		c.Abort()
		return
	}

	event := models.AnalyticsEvent{
		EventType: req.EventType,
		Path:      req.Path,
		CtaLabel:  req.CtaLabel,
	}
	if err := h.db.Create(&event).Error; err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Event tracked", "data": gin.H{}})
}

type countRow struct {
	Label string `json:"label"`
	Count int64  `json:"count"`
}

type summary struct {
	TotalPageViews int64      `json:"totalPageViews"`
	TotalCtaClicks int64      `json:"totalCtaClicks"`
	TodayPageViews int64      `json:"todayPageViews"`
	TodayCtaClicks int64      `json:"todayCtaClicks"`
	CtaBreakdown   []countRow `json:"ctaBreakdown"`
	TopPages       []countRow `json:"topPages"`
}

// Summary is dashboard-only (auth required) — aggregate counts for the
// analytics page.
func (h *Handler) Summary(c *gin.Context) {
	var out summary
	startOfDay := time.Now().Truncate(24 * time.Hour)

	h.db.Model(&models.AnalyticsEvent{}).Where("event_type = ?", "page_view").Count(&out.TotalPageViews)
	h.db.Model(&models.AnalyticsEvent{}).Where("event_type = ?", "cta_click").Count(&out.TotalCtaClicks)
	h.db.Model(&models.AnalyticsEvent{}).Where("event_type = ? AND created_at >= ?", "page_view", startOfDay).Count(&out.TodayPageViews)
	h.db.Model(&models.AnalyticsEvent{}).Where("event_type = ? AND created_at >= ?", "cta_click", startOfDay).Count(&out.TodayCtaClicks)

	h.db.Model(&models.AnalyticsEvent{}).
		Select("cta_label as label, count(*) as count").
		Where("event_type = ? AND cta_label <> ''", "cta_click").
		Group("cta_label").
		Order("count DESC").
		Scan(&out.CtaBreakdown)

	h.db.Model(&models.AnalyticsEvent{}).
		Select("path as label, count(*) as count").
		Where("event_type = ?", "page_view").
		Group("path").
		Order("count DESC").
		Limit(10).
		Scan(&out.TopPages)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Analytics summary fetched", "data": out})
}
