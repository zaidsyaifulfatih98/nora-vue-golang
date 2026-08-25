package middleware

import (
	"mime/multipart"
	"strings"

	"github.com/gin-gonic/gin"

	"nora-photobooth-backend/internal/apperror"
)

const maxImageSize = 5 << 20 // 5MB

// IsMultipart reports whether the request body is multipart/form-data, i.e.
// carries a new image, as opposed to a plain JSON PATCH (no image change).
func IsMultipart(c *gin.Context) bool {
	return strings.HasPrefix(c.ContentType(), "multipart/form-data")
}

// ExtractImage pulls the "image" multipart field, enforcing the same
// 5MB size limit and image/* mimetype filter as the old multer config.
// It returns (nil, nil) when no file was provided (image optional on update).
func ExtractImage(c *gin.Context) (multipart.File, error) {
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return nil, nil
	}

	if fileHeader.Size > maxImageSize {
		return nil, apperror.New("Ukuran gambar maksimal 5MB", 400)
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return nil, apperror.New("File harus berupa gambar", 400)
	}

	return fileHeader.Open()
}
