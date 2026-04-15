package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Role     string `json:"role"`
	TenantID string `json:"tenantId"`
	jwt.RegisteredClaims
}

func ParseAccessToken(tokenString, secret string) (*Claims, error) {
	if tokenString == "" {
		return nil, errors.New("token is required")
	}
	if secret == "" {
		return nil, errors.New("jwt secret is required")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
