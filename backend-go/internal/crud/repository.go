package crud

import (
	"errors"

	"gorm.io/gorm"

	"nora-photobooth-backend/internal/apperror"
)

// Repository is a generic soft-delete CRUD repository shared by the six
// similar resources (Package, Feature, FrameTemplate, Backdrop, GalleryPhoto,
// Review). VisibilityColumn lets Review use "is_published" instead of the
// "is_active" column every other resource uses.
type Repository[T any] struct {
	DB               *gorm.DB
	ResourceName     string
	VisibilityColumn string
}

func NewRepository[T any](db *gorm.DB, resourceName, visibilityColumn string) *Repository[T] {
	if visibilityColumn == "" {
		visibilityColumn = "is_active"
	}
	return &Repository[T]{DB: db, ResourceName: resourceName, VisibilityColumn: visibilityColumn}
}

func (r *Repository[T]) List(includeInactive bool) ([]T, error) {
	var items []T
	q := r.DB.Order("\"order\" asc")
	if !includeInactive {
		q = q.Where(r.VisibilityColumn+" = ?", true)
	}
	if err := q.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository[T]) FindByID(id string) (*T, error) {
	var item T
	err := r.DB.Where("id = ?", id).First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.New(r.ResourceName+" tidak ditemukan", 404)
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *Repository[T]) Create(item *T) error {
	return r.DB.Create(item).Error
}

func (r *Repository[T]) Update(id string, updates map[string]any) (*T, error) {
	item, err := r.FindByID(id)
	if err != nil {
		return nil, err
	}
	if err := r.DB.Model(item).Updates(updates).Error; err != nil {
		return nil, err
	}
	return r.FindByID(id)
}

func (r *Repository[T]) SoftDelete(id string) error {
	item, err := r.FindByID(id)
	if err != nil {
		return err
	}
	return r.DB.Delete(item).Error
}
