package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"nora-photobooth-backend/internal/models"
)

type Claims struct {
	ID   string      `json:"id"`
	Role models.Role `json:"role"`
	jwt.RegisteredClaims
}

func Create(id string, role models.Role, secret string) (string, error) {
	claims := Claims{
		ID:   id,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString([]byte(secret))
}

func Verify(tokenString string, secret string) (*Claims, error) {
	claims := &Claims{}
	tok, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !tok.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
