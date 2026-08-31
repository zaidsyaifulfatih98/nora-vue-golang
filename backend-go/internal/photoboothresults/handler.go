package photoboothresults

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"nora-photobooth-backend/internal/apperror"
	"nora-photobooth-backend/internal/middleware"
	"nora-photobooth-backend/internal/models"
	"nora-photobooth-backend/internal/upload"
)

// Handler saves a guest's finished digital photobooth result (uploaded via
// Cloudinary, same as every other image in this app), records it so the
// dashboard can browse the collection, and hands back a public download
// link, so it can be turned into a QR code for the guest to scan and save
// the photo to their own phone.
//
// Google Drive was tried first, but a bare service account has no Drive
// storage quota of its own — uploads only work into a Google Workspace
// Shared Drive, which isn't available on a personal Gmail account.
type Handler struct {
	db       *gorm.DB
	uploader *upload.Uploader
}

func NewHandler(db *gorm.DB, uploader *upload.Uploader) *Handler {
	return &Handler{db: db, uploader: uploader}
}

// downloadURL forces a download instead of opening in the browser when the
// QR code is scanned, via Cloudinary's fl_attachment delivery flag.
func downloadURL(imageURL string) string {
	return strings.Replace(imageURL, "/upload/", "/upload/fl_attachment/", 1)
}

// Create is public — a guest saves their finished result from the digital
// photobooth result screen, no auth required.
func (h *Handler) Create(c *gin.Context) {
	file, err := middleware.ExtractImage(c)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}
	if file == nil {
		_ = c.Error(apperror.New("Gambar wajib diunggah", 400))
		c.Abort()
		return
	}
	defer file.Close()

	imageURL, err := h.uploader.UploadImage(file)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	item := models.PhotoboothResult{ImageURL: imageURL}
	if err := h.db.Create(&item).Error; err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Photobooth result saved",
		"data": gin.H{
			"id":          item.ID,
			"createdAt":   item.CreatedAt,
			"viewUrl":     imageURL,
			"downloadUrl": downloadURL(imageURL),
		},
	})
}

// List is admin-only — the dashboard's collection of saved digital
// photobooth results.
func (h *Handler) List(c *gin.Context) {
	var items []models.PhotoboothResult
	if err := h.db.Order("created_at desc").Find(&items).Error; err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	data := make([]gin.H, len(items))
	for i, item := range items {
		data[i] = gin.H{
			"id":          item.ID,
			"createdAt":   item.CreatedAt,
			"viewUrl":     item.ImageURL,
			"downloadUrl": downloadURL(item.ImageURL),
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Photobooth results fetched", "data": data})
}

// Delete is admin-only.
func (h *Handler) Delete(c *gin.Context) {
	if err := h.db.Delete(&models.PhotoboothResult{}, "id = ?", c.Param("id")).Error; err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Photobooth result deleted", "data": gin.H{}})
}
