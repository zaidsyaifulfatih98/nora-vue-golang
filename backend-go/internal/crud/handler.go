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
	return func(c *gin.Context) {
		id := c.Param("id")
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
