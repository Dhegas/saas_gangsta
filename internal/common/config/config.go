package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv  string
	AppPort string
	AppName string

	DatabaseURL            string
	RedisURL               string
	SupabaseURL            string
	SupabaseAnonKey        string
	SupabaseServiceRoleKey string

	JWTSecret             string
	JWTAccessTokenExpiry  time.Duration
	JWTRefreshTokenExpiry time.Duration

	CORSAllowedOrigins []string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv:                 getEnv("APP_ENV", "development"),
		AppPort:                resolveAppPort(),
		AppName:                getEnv("APP_NAME", "saas_gangsta"),
		DatabaseURL:            strings.TrimSpace(getEnv("DATABASE_URL", "")),
		RedisURL:               strings.TrimSpace(getEnv("REDIS_URL", "redis://localhost:6379/0")),
		SupabaseURL:            strings.TrimSpace(getEnv("SUPABASE_URL", "")),
		SupabaseAnonKey:        strings.TrimSpace(getEnv("SUPABASE_ANON_KEY", "")),
		SupabaseServiceRoleKey: strings.TrimSpace(getEnv("SUPABASE_SERVICE_ROLE_KEY", "")),
		JWTSecret:              strings.TrimSpace(getEnv("JWT_SECRET", "")),
		CORSAllowedOrigins:     splitCSV(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")),
		JWTAccessTokenExpiry:   15 * time.Minute,
		JWTRefreshTokenExpiry:  7 * 24 * time.Hour,
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	if err := validateDatabaseURL(cfg.DatabaseURL); err != nil {
		return nil, err
	}

	// Supabase direct integration values should be configured together.
	if (cfg.SupabaseURL == "") != (cfg.SupabaseServiceRoleKey == "") {
		return nil, fmt.Errorf("SUPABASE_URL and SUPABASE_SERVICE_ROLE_KEY must be set together")
	}

	if strings.EqualFold(cfg.AppEnv, "production") {
		if cfg.SupabaseURL == "" {
			return nil, fmt.Errorf("SUPABASE_URL is required in production")
		}
		if cfg.SupabaseServiceRoleKey == "" {
			return nil, fmt.Errorf("SUPABASE_SERVICE_ROLE_KEY is required in production")
		}
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

func validateDatabaseURL(raw string) error {
	parsed, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("invalid DATABASE_URL: %w. If password contains special characters (for example @, :, /, #, ?), URL-encode it first", err)
	}

	if parsed.Scheme != "postgres" && parsed.Scheme != "postgresql" {
		return fmt.Errorf("invalid DATABASE_URL scheme %q: use postgresql:// or postgres://", parsed.Scheme)
	}

	if parsed.Hostname() == "" {
		return fmt.Errorf("invalid DATABASE_URL: hostname is required")
	}

	return nil
}
