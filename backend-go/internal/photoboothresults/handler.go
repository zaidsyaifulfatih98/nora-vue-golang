package photoboothresults

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"nora-photobooth-backend/internal/apperror"
	"nora-photobooth-backend/internal/middleware"
	"nora-photobooth-backend/internal/upload"
)

// Handler saves a guest's finished digital photobooth result (uploaded via
// Cloudinary, same as every other image in this app) and hands back a
// public download link, so it can be turned into a QR code for the guest
// to scan and save the photo to their own phone.
//
// Google Drive was tried first, but a bare service account has no Drive
// storage quota of its own — uploads only work into a Google Workspace
// Shared Drive, which isn't available on a personal Gmail account.
type Handler struct {
	uploader *upload.Uploader
}

func NewHandler(uploader *upload.Uploader) *Handler {
	return &Handler{uploader: uploader}
}

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

	// Force a download instead of opening in the browser when the QR code
	// is scanned, via Cloudinary's fl_attachment delivery flag.
	downloadURL := strings.Replace(imageURL, "/upload/", "/upload/fl_attachment/", 1)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Photobooth result saved",
		"data": gin.H{
			"viewUrl":     imageURL,
			"downloadUrl": downloadURL,
		},
	})
}
