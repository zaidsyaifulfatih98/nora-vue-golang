package gallery

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"nora-photobooth-backend/internal/apperror"
	"nora-photobooth-backend/internal/crud"
	"nora-photobooth-backend/internal/middleware"
	"nora-photobooth-backend/internal/models"
	"nora-photobooth-backend/internal/upload"
)

type Handler struct {
	repo     *crud.Repository[models.GalleryPhoto]
	uploader *upload.Uploader
}

func NewHandler(db *gorm.DB, uploader *upload.Uploader) *Handler {
	return &Handler{repo: crud.NewRepository[models.GalleryPhoto](db, "GalleryPhoto", "is_active"), uploader: uploader}
}

func (h *Handler) List() gin.HandlerFunc   { return crud.ListHandler(h.repo) }
func (h *Handler) Delete() gin.HandlerFunc { return crud.DeleteHandler(h.repo) }

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

	item := models.GalleryPhoto{URL: imageURL, Caption: c.PostForm("caption"), IsActive: true}
	if err := h.repo.Create(&item); err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Gallery photo created", "data": item})
}

type jsonUpdateRequest struct {
	Caption  *string `json:"caption"`
	IsActive *bool   `json:"isActive"`
	Order    *int    `json:"order"`
}

func (h *Handler) Update(c *gin.Context) {
	updates := map[string]any{}

	if middleware.IsMultipart(c) {
		if v := c.PostForm("caption"); v != "" {
			updates["caption"] = v
		}
		if v := c.PostForm("isActive"); v != "" {
			updates["is_active"] = v == "true"
		}
		if v := c.PostForm("order"); v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				updates["order"] = n
			}
		}

		file, err := middleware.ExtractImage(c)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}
		if file != nil {
			defer file.Close()
			imageURL, err := h.uploader.UploadImage(file)
			if err != nil {
				_ = c.Error(err)
				c.Abort()
				return
			}
			updates["url"] = imageURL
		}
	} else {
		var req jsonUpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			_ = c.Error(apperror.New("Data galeri tidak valid", 400))
			c.Abort()
			return
		}
		if req.Caption != nil {
			updates["caption"] = *req.Caption
		}
		if req.IsActive != nil {
			updates["is_active"] = *req.IsActive
		}
		if req.Order != nil {
			updates["order"] = *req.Order
		}
	}

	item, err := h.repo.Update(c.Param("id"), updates)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Gallery photo updated", "data": item})
}
