package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv  string
	AppPort string
	AppName string

	DatabaseURL string
	RedisURL    string

	JWTSecret             string
	JWTAccessTokenExpiry  time.Duration
	JWTRefreshTokenExpiry time.Duration

	CORSAllowedOrigins []string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv:                getEnv("APP_ENV", "development"),
		AppPort:               resolveAppPort(),
		AppName:               getEnv("APP_NAME", "saas_gangsta"),
		DatabaseURL:           strings.TrimSpace(getEnv("DATABASE_URL", "")),
		RedisURL:              strings.TrimSpace(getEnv("REDIS_URL", "redis://localhost:6379/0")),
		JWTSecret:             strings.TrimSpace(getEnv("JWT_SECRET", "")),
		CORSAllowedOrigins:    splitCSV(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")),
		JWTAccessTokenExpiry:  15 * time.Minute,
		JWTRefreshTokenExpiry: 7 * 24 * time.Hour,
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	if access, err := time.ParseDuration(getEnv("JWT_ACCESS_TOKEN_EXPIRY", "15m")); err == nil {
		cfg.JWTAccessTokenExpiry = access
	} else {
		return nil, fmt.Errorf("invalid JWT_ACCESS_TOKEN_EXPIRY: %w", err)
	}

	if refresh, err := time.ParseDuration(getEnv("JWT_REFRESH_TOKEN_EXPIRY", "168h")); err == nil {
		cfg.JWTRefreshTokenExpiry = refresh
	} else {
		return nil, fmt.Errorf("invalid JWT_REFRESH_TOKEN_EXPIRY: %w", err)
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return fallback
	}
	return val
}

func splitCSV(raw string) []string {
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	if len(result) == 0 {
		return []string{"http://localhost:3000"}
	}
	return result
}

func resolveAppPort() string {
	appPort := strings.TrimSpace(os.Getenv("APP_PORT"))
	if appPort != "" {
		return appPort
	}

	railwayPort := strings.TrimSpace(os.Getenv("PORT"))
	if railwayPort != "" {
		return railwayPort
	}

	return "8080"
}
