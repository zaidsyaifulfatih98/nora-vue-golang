package frametemplates

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
	repo     *crud.Repository[models.FrameTemplate]
	uploader *upload.Uploader
}

func NewHandler(db *gorm.DB, uploader *upload.Uploader) *Handler {
	return &Handler{repo: crud.NewRepository[models.FrameTemplate](db, "FrameTemplate", "is_active"), uploader: uploader}
}

func (h *Handler) List() gin.HandlerFunc   { return crud.ListHandler(h.repo) }
func (h *Handler) Delete() gin.HandlerFunc { return crud.DeleteHandler(h.repo) }

func isValidFit(fit string) bool {
	return fit == string(models.FrameFitCover) || fit == string(models.FrameFitContain)
}

func (h *Handler) Create(c *gin.Context) {
	name := c.PostForm("name")
	description := c.PostForm("description")
	fit := c.PostForm("fit")
	if fit == "" {
		fit = string(models.FrameFitCover)
	}
	if name == "" || description == "" {
		_ = c.Error(apperror.New("Name dan description wajib diisi", 400))
		c.Abort()
		return
	}
	if !isValidFit(fit) {
		_ = c.Error(apperror.New("Fit harus COVER atau CONTAIN", 400))
		c.Abort()
		return
	}

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

	item := models.FrameTemplate{
		Name:        name,
		Description: description,
		ImageURL:    imageURL,
		Fit:         models.FrameTemplateFit(fit),
		IsActive:    true,
	}

	if err := h.repo.Create(&item); err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Frame template created", "data": item})
}

type jsonUpdateRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Fit         *string `json:"fit"`
	IsActive    *bool   `json:"isActive"`
	Order       *int    `json:"order"`
}

func (h *Handler) Update(c *gin.Context) {
	updates := map[string]any{}

	if middleware.IsMultipart(c) {
		if v := c.PostForm("name"); v != "" {
			updates["name"] = v
		}
		if v := c.PostForm("description"); v != "" {
			updates["description"] = v
		}
		if v := c.PostForm("fit"); v != "" {
			if !isValidFit(v) {
				_ = c.Error(apperror.New("Fit harus COVER atau CONTAIN", 400))
				c.Abort()
				return
			}
			updates["fit"] = v
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
			updates["image_url"] = imageURL
		}
	} else {
		var req jsonUpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			_ = c.Error(apperror.New("Data frame template tidak valid", 400))
			c.Abort()
			return
		}
		if req.Name != nil {
			updates["name"] = *req.Name
		}
		if req.Description != nil {
			updates["description"] = *req.Description
		}
		if req.Fit != nil {
			if !isValidFit(*req.Fit) {
				_ = c.Error(apperror.New("Fit harus COVER atau CONTAIN", 400))
				c.Abort()
				return
			}
			updates["fit"] = *req.Fit
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

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Frame template updated", "data": item})
}
