package photoboothframes

import (
	"encoding/json"
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

// parseSlots validates that raw is a JSON array (of slot objects), returning
// "[]" when raw is empty so the column is never left invalid.
func parseSlots(raw string) (models.JSONText, error) {
	if raw == "" {
		return models.JSONText("[]"), nil
	}
	var v []map[string]float64
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return nil, apperror.New("Data slot frame tidak valid", 400)
	}
	return models.JSONText(raw), nil
}

type Handler struct {
	repo     *crud.Repository[models.PhotoboothFrame]
	uploader *upload.Uploader
}

func NewHandler(db *gorm.DB, uploader *upload.Uploader) *Handler {
	return &Handler{repo: crud.NewRepository[models.PhotoboothFrame](db, "PhotoboothFrame", "is_active"), uploader: uploader}
}

func (h *Handler) List() gin.HandlerFunc   { return crud.ListHandler(h.repo) }
func (h *Handler) Delete() gin.HandlerFunc { return crud.DeleteHandler(h.repo) }

func (h *Handler) Create(c *gin.Context) {
	name := c.PostForm("name")
	if name == "" {
		_ = c.Error(apperror.New("Nama frame wajib diisi", 400))
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

	slots, err := parseSlots(c.PostForm("slots"))
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	item := models.PhotoboothFrame{Name: name, ImageURL: imageURL, Slots: slots, IsActive: true}
	if err := h.repo.Create(&item); err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Photobooth frame created", "data": item})
}

type jsonUpdateRequest struct {
	Name     *string          `json:"name"`
	Slots    *models.JSONText `json:"slots"`
	IsActive *bool            `json:"isActive"`
	Order    *int             `json:"order"`
}

func (h *Handler) Update(c *gin.Context) {
	updates := map[string]any{}

	if middleware.IsMultipart(c) {
		if v := c.PostForm("name"); v != "" {
			updates["name"] = v
		}
		if v := c.PostForm("isActive"); v != "" {
			updates["is_active"] = v == "true"
		}
		if v := c.PostForm("order"); v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				updates["order"] = n
			}
		}
		if v, ok := c.GetPostForm("slots"); ok {
			slots, err := parseSlots(v)
			if err != nil {
				_ = c.Error(err)
				c.Abort()
				return
			}
			updates["slots"] = slots
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
			_ = c.Error(apperror.New("Data frame tidak valid", 400))
			c.Abort()
			return
		}
		if req.Name != nil {
			updates["name"] = *req.Name
		}
		if req.Slots != nil {
			slots, err := parseSlots(string(*req.Slots))
			if err != nil {
				_ = c.Error(err)
				c.Abort()
				return
			}
			updates["slots"] = slots
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

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Photobooth frame updated", "data": item})
}
