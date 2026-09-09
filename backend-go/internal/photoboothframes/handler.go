package photoboothframes

import (
	"encoding/json"
	"errors"
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

// List is the public, unauthenticated main-site listing — unchanged from
// before this package supported customer-owned frames: it only ever returns
// frames with no owner (owner_id IS NULL).
func (h *Handler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		all := c.Query("all") == "true"
		var items []models.PhotoboothFrame
		q := h.repo.DB.Where("owner_id IS NULL").Order("\"order\" asc")
		if !all {
			q = q.Where("is_active = ?", true)
		}
		if err := q.Find(&items).Error; err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "PhotoboothFrame fetched", "data": items})
	}
}

// Mine returns the caller's own frames — a DIGITAL_PHOTOBOOTH customer's
// dashboard listing (or the full admin listing when an admin/superadmin
// happens to call it, though the dashboard doesn't use this route for them).
func (h *Handler) Mine(c *gin.Context) {
	userID := c.GetString("userID")
	var items []models.PhotoboothFrame
	if err := h.repo.DB.Where("owner_id = ?", userID).Order("\"order\" asc").Find(&items).Error; err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "PhotoboothFrame fetched", "data": items})
}

// ByOwnerSlug is public — it resolves a customer's slug (DIGITAL_PHOTOBOOTH
// or SOFTWARE_PHOTOBOOTH) to their active frames, powering that customer's
// guest try-it / kiosk sub-page.
func (h *Handler) ByOwnerSlug(c *gin.Context) {
	slug := c.Param("slug")

	var owner models.User
	err := h.repo.DB.Where("slug = ? AND role IN ?", slug, []models.Role{models.RoleDigitalPhotobooth, models.RoleSoftwarePhotobooth}).First(&owner).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		_ = c.Error(apperror.New("Halaman tidak ditemukan", 404))
		c.Abort()
		return
	}
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	var frames []models.PhotoboothFrame
	if err := h.repo.DB.Where("owner_id = ? AND is_active = ?", owner.ID, true).Order("\"order\" asc").Find(&frames).Error; err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "PhotoboothFrame fetched",
		"data": gin.H{
			"owner":  gin.H{"firstName": owner.FirstName, "lastName": owner.LastName},
			"frames": frames,
		},
	})
}

// ownsFrame reports whether the caller may modify the given frame: an
// admin/superadmin may modify any frame, a customer (DIGITAL_PHOTOBOOTH or
// SOFTWARE_PHOTOBOOTH) only their own (owner_id = their id). Writes a 403
// and returns false if not.
func (h *Handler) ownsFrame(c *gin.Context, item *models.PhotoboothFrame) bool {
	if !models.Role(c.GetString("userRole")).IsCustomer() {
		return true
	}
	if item.OwnerID != nil && *item.OwnerID == c.GetString("userID") {
		return true
	}
	_ = c.Error(apperror.New("Forbidden, insufficient role", 403))
	c.Abort()
	return false
}

func (h *Handler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		item, err := h.repo.FindByID(c.Param("id"))
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}
		if !h.ownsFrame(c, item) {
			return
		}
		_ = h.uploader.DeleteImage(item.ImageURL)
		if err := h.repo.SoftDelete(item.ID); err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "PhotoboothFrame deleted", "data": gin.H{}})
	}
}

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

	var ownerID *string
	if models.Role(c.GetString("userRole")).IsCustomer() {
		// A customer's own upload is always scoped to themselves — never
		// trust a client-supplied ownerId here.
		userID := c.GetString("userID")
		ownerID = &userID
	} else if raw := c.PostForm("ownerId"); raw != "" {
		var owner models.User
		err := h.repo.DB.Where("id = ? AND role IN ?", raw, []models.Role{models.RoleDigitalPhotobooth, models.RoleSoftwarePhotobooth}).First(&owner).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			_ = c.Error(apperror.New("Pelanggan tidak ditemukan", 400))
			c.Abort()
			return
		}
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}
		ownerID = &owner.ID
	}

	item := models.PhotoboothFrame{Name: name, ImageURL: imageURL, Slots: slots, IsActive: true, OwnerID: ownerID}
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
	existing, err := h.repo.FindByID(c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}
	if !h.ownsFrame(c, existing) {
		return
	}

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

	var oldImageURL string
	if newImageURL, replacingImage := updates["image_url"]; replacingImage {
		if existing, err := h.repo.FindByID(c.Param("id")); err == nil && existing.ImageURL != newImageURL {
			oldImageURL = existing.ImageURL
		}
	}

	item, err := h.repo.Update(c.Param("id"), updates)
	if err != nil {
		_ = c.Error(err)
		c.Abort()
		return
	}

	if oldImageURL != "" {
		_ = h.uploader.DeleteImage(oldImageURL)
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Photobooth frame updated", "data": item})
}
