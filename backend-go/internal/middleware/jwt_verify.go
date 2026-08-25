package middleware

import (
	"github.com/gin-gonic/gin"

	"nora-photobooth-backend/internal/apperror"
	"nora-photobooth-backend/internal/token"
)

// JWTVerify reads the "token" cookie, verifies it, and stores userID/role in
// the Gin context. Any failure yields a generic 401, matching the old
// jwt-verify.middlewere.ts (which never leaks the underlying jwt error).
func JWTVerify(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookieValue, err := c.Cookie("token")
		if err != nil || cookieValue == "" {
			Fail(c, apperror.New("Unauthorized, please login", 401))
			return
		}

		claims, err := token.Verify(cookieValue, jwtSecret)
		if err != nil {
			Fail(c, apperror.New("Unauthorized, please login", 401))
			return
		}

		c.Set("userID", claims.ID)
		c.Set("userRole", string(claims.Role))
		c.Next()
	}
}
