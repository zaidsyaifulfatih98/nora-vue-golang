package config

import (
	"time"

	"github.com/gin-contrib/cors"
)

// BuildCORS mirrors the old dynamic-origin-whitelist behaviour: requests with
// no Origin header (curl / server-to-server) are always allowed, otherwise
// the origin must be in the whitelist.
func BuildCORS(cfg *Config) cors.Config {
	whitelist := make(map[string]bool, len(cfg.CORSWhitelist))
	for _, o := range cfg.CORSWhitelist {
		whitelist[o] = true
	}

	return cors.Config{
		AllowOriginFunc: func(origin string) bool {
			if origin == "" {
				return true
			}
			return whitelist[origin]
		},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
}
