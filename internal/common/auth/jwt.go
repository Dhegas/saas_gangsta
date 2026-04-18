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

func GenerateAccessToken(userID, role, tenantID string, expiry time.Duration, secret string) (string, error) {
	return generateToken(userID, role, tenantID, expiry, secret)
}

func GenerateRefreshToken(userID, role, tenantID string, expiry time.Duration, secret string) (string, error) {
	return generateToken(userID, role, tenantID, expiry, secret)
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

func ParseRefreshToken(tokenString, secret string) (*Claims, error) {
	return ParseAccessToken(tokenString, secret)
}

func generateToken(userID, role, tenantID string, expiry time.Duration, secret string) (string, error) {
	if userID == "" {
		return "", errors.New("user id is required")
	}
	if secret == "" {
		return "", errors.New("jwt secret is required")
	}
	if expiry <= 0 {
		return "", errors.New("token expiry must be greater than zero")
	}

	now := time.Now()
	claims := Claims{
		Role:     role,
		TenantID: tenantID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
