package features

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
	repo     *crud.Repository[models.Feature]
	uploader *upload.Uploader
}

func NewHandler(db *gorm.DB, uploader *upload.Uploader) *Handler {
	return &Handler{repo: crud.NewRepository[models.Feature](db, "Feature", "is_active"), uploader: uploader}
}

func (h *Handler) List() gin.HandlerFunc   { return crud.ListHandler(h.repo) }
func (h *Handler) Delete() gin.HandlerFunc { return crud.DeleteHandler(h.repo) }

func (h *Handler) Create(c *gin.Context) {
	title := c.PostForm("title")
	description := c.PostForm("description")
	titleEn := c.PostForm("titleEn")
	descriptionEn := c.PostForm("descriptionEn")
	if title == "" || description == "" {
		_ = c.Error(apperror.New("Title dan description wajib diisi", 400))
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

	item := models.Feature{
		Title:         title,
		Description:   description,
		TitleEn:       middleware.NilIfEmpty(titleEn),
		DescriptionEn: middleware.NilIfEmpty(descriptionEn),
		ImageURL:      imageURL,
		IsActive:      true,
	}

	if err := h.repo.Create(&item); err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Feature created", "data": item})
}

type jsonUpdateRequest struct {
	Title         *string `json:"title"`
	Description   *string `json:"description"`
	TitleEn       *string `json:"titleEn"`
	DescriptionEn *string `json:"descriptionEn"`
	IsActive      *bool   `json:"isActive"`
	Order         *int    `json:"order"`
}

func (h *Handler) Update(c *gin.Context) {
	updates := map[string]any{}

	if middleware.IsMultipart(c) {
		if v := c.PostForm("title"); v != "" {
			updates["title"] = v
		}
		if v := c.PostForm("description"); v != "" {
			updates["description"] = v
		}
		if v := c.PostForm("titleEn"); v != "" {
			updates["title_en"] = v
		}
		if v := c.PostForm("descriptionEn"); v != "" {
			updates["description_en"] = v
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
			_ = c.Error(apperror.New("Data feature tidak valid", 400))
			c.Abort()
			return
		}
		if req.Title != nil {
			updates["title"] = *req.Title
		}
		if req.Description != nil {
			updates["description"] = *req.Description
		}
		if req.TitleEn != nil {
			updates["title_en"] = *req.TitleEn
		}
		if req.DescriptionEn != nil {
			updates["description_en"] = *req.DescriptionEn
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

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Feature updated", "data": item})
}
