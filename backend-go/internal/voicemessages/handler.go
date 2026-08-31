package voicemessages

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"nora-photobooth-backend/internal/apperror"
	"nora-photobooth-backend/internal/middleware"
	"nora-photobooth-backend/internal/models"
	"nora-photobooth-backend/internal/upload"
)

type Handler struct {
	db       *gorm.DB
	uploader *upload.Uploader
}

func NewHandler(db *gorm.DB, uploader *upload.Uploader) *Handler {
	return &Handler{db: db, uploader: uploader}
}

// Create is public — a guest submits a voice greeting from the digital
// photobooth result screen, no auth required.
func (h *Handler) Create(c *gin.Context) {
	file, err := middleware.ExtractAudio(c)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}
	if file == nil {
		_ = c.Error(apperror.New("Audio wajib diunggah", 400))
		c.Abort()
		return
	}
	defer file.Close()

	audioURL, err := h.uploader.UploadAudio(file)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	item := models.VoiceMessage{
		GuestName: c.PostForm("guestName"),
		AudioURL:  audioURL,
		PhotoURL:  c.PostForm("photoUrl"),
	}
	if err := h.db.Create(&item).Error; err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Voice message saved", "data": item})
}

// List is admin-only — the dashboard's collection of guest voice messages.
func (h *Handler) List(c *gin.Context) {
	var items []models.VoiceMessage
	if err := h.db.Order("created_at desc").Find(&items).Error; err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Voice messages fetched", "data": items})
}

// Delete is admin-only.
func (h *Handler) Delete(c *gin.Context) {
	if err := h.db.Delete(&models.VoiceMessage{}, "id = ?", c.Param("id")).Error; err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Voice message deleted", "data": gin.H{}})
}
