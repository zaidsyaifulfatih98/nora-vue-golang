package middleware

import (
	"github.com/gin-gonic/gin"

	"nora-photobooth-backend/internal/apperror"
)

// RoleVerify must run after JWTVerify. It checks that the authenticated
// user's role is one of allowedRoles.
func RoleVerify(allowedRoles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(allowedRoles))
	for _, r := range allowedRoles {
		allowed[r] = true
	}

	return func(c *gin.Context) {
		role := c.GetString("userRole")
		if role == "" || !allowed[role] {
			Fail(c, apperror.New("Forbidden, insufficient role", 403))
			return
		}
		c.Next()
	}
}
