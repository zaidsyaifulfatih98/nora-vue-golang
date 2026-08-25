package auth

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

// firstBindingError extracts a single readable message from a Gin binding
// error, mirroring the old expressValidation middleware's "first error only"
// behaviour.
func firstBindingError(err error, fallback string) string {
	if verrs, ok := err.(validator.ValidationErrors); ok && len(verrs) > 0 {
		fe := verrs[0]
		return fe.Field() + " " + fe.Tag()
	}
	return fallback
}

func sameSiteFor(isProd bool) http.SameSite {
	if isProd {
		return http.SameSiteNoneMode
	}
	return http.SameSiteLaxMode
}
