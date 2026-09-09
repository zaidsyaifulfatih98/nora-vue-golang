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
// photobooth result screen, no auth required. An optional ownerSlug (present
// when the guest used a customer's sub-URL) tags the message to that
// customer; it's always resolved server-side through the slug, never taken
// as a raw id, so a guest can't tag a message to an arbitrary account.
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

	var ownerID *string
	if slug := c.PostForm("ownerSlug"); slug != "" {
		var owner models.User
		if err := h.db.Where("slug = ? AND role IN ?", slug, []models.Role{models.RoleDigitalPhotobooth, models.RoleSoftwarePhotobooth}).First(&owner).Error; err == nil {
			ownerID = &owner.ID
		}
	}

	item := models.VoiceMessage{
		GuestName: c.PostForm("guestName"),
		AudioURL:  audioURL,
		PhotoURL:  c.PostForm("photoUrl"),
		OwnerID:   ownerID,
	}
	if err := h.db.Create(&item).Error; err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Voice message saved", "data": item})
}

// List is admin/customer — the dashboard's collection of guest voice
// messages. A customer (DIGITAL_PHOTOBOOTH or SOFTWARE_PHOTOBOOTH) only
// sees their own (owner_id = their id); admin/superadmin only see the main
// site's (owner_id IS NULL), unchanged from before customer accounts
// existed.
func (h *Handler) List(c *gin.Context) {
	q := h.db.Order("created_at desc")
	if models.Role(c.GetString("userRole")).IsCustomer() {
		q = q.Where("owner_id = ?", c.GetString("userID"))
	} else {
		q = q.Where("owner_id IS NULL")
	}

	var items []models.VoiceMessage
	if err := q.Find(&items).Error; err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Voice messages fetched", "data": items})
}

// Delete is admin-only. Cloudinary cleanup is best-effort — an unreachable
// Cloudinary shouldn't block the delete the admin actually asked for. Only
// the recording itself is removed; PhotoURL (if set) points at a separate
// PhotoboothResult that manages its own Cloudinary asset lifecycle.
func (h *Handler) Delete(c *gin.Context) {
	var item models.VoiceMessage
	if err := h.db.First(&item, "id = ?", c.Param("id")).Error; err == nil {
		_ = h.uploader.DeleteAudio(item.AudioURL)
	}

	if err := h.db.Delete(&models.VoiceMessage{}, "id = ?", c.Param("id")).Error; err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Voice message deleted", "data": gin.H{}})
}
