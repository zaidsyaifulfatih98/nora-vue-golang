package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"nora-photobooth-backend/internal/apperror"
	"nora-photobooth-backend/internal/logging"
)

// ErrorHandler centralizes error responses to match the old Express global
// handler's {success:false, message, data:{}} envelope.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		if appErr, ok := err.(*apperror.AppError); ok && appErr.Expose {
			logging.Log.Warn().
				Str("path", c.Request.URL.Path).
				Int("status", appErr.StatusCode).
				Msg(appErr.Message)
			c.JSON(appErr.StatusCode, gin.H{
				"success": false,
				"message": appErr.Message,
				"data":    gin.H{},
			})
			return
		}

		logging.Log.Error().
			Str("path", c.Request.URL.Path).
			Err(err).
			Msg("unhandled error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "something went wrong",
			"data":    gin.H{},
		})
	}
}

// Fail is a small helper for handlers: attach an AppError and abort.
func Fail(c *gin.Context, err error) {
	_ = c.Error(err)
	c.Abort()
}
