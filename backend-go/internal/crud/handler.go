package crud

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListHandler returns a Gin handler for GET list endpoints, honoring the
// shared ?all=true convention used across every resource.
func ListHandler[T any](repo *Repository[T]) gin.HandlerFunc {
	return func(c *gin.Context) {
		all := c.Query("all") == "true"
		items, err := repo.List(all)
		if err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": repo.ResourceName + " fetched",
			"data":    items,
		})
	}
}

// DeleteHandler returns a Gin handler for DELETE endpoints (soft delete).
func DeleteHandler[T any](repo *Repository[T]) gin.HandlerFunc {
	return DeleteHandlerWithCleanup(repo, nil)
}

// DeleteHandlerWithCleanup is DeleteHandler plus a hook run against the item
// right before it's soft-deleted — e.g. removing its uploaded image from
// Cloudinary, which the database row alone has no way to do on its own.
// cleanup runs best-effort: its error is ignored so an unreachable Cloudinary
// never blocks the delete the user actually asked for.
func DeleteHandlerWithCleanup[T any](repo *Repository[T], cleanup func(item *T) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if cleanup != nil {
			if item, err := repo.FindByID(id); err == nil {
				_ = cleanup(item)
			}
		}

		if err := repo.SoftDelete(id); err != nil {
			_ = c.Error(err)
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": repo.ResourceName + " deleted",
			"data":    gin.H{},
		})
	}
}
