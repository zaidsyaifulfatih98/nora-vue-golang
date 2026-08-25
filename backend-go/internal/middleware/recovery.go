package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"nora-photobooth-backend/internal/logging"
)

// Recovery catches panics and responds with the same envelope as ErrorHandler,
// instead of gin's default plain-text 500.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logging.Log.Error().Interface("panic", r).Str("path", c.Request.URL.Path).Msg("panic recovered")
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "something went wrong",
					"data":    gin.H{},
				})
			}
		}()
		c.Next()
	}
}
