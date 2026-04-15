package middleware

import (
	"net/http"
	"strings"

	"github.com/dhegas/saas_gangsta/internal/common/auth"
	"github.com/dhegas/saas_gangsta/internal/common/config"
	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/common/tenant"
	"github.com/gin-gonic/gin"
)

const (
	UserIDKey = "userId"
	RoleKey   = "role"
)

func JWTAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		parts := strings.SplitN(authorization, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			apperrors.Abort(c, apperrors.New("UNAUTHORIZED", "Missing or invalid Authorization header", http.StatusUnauthorized, nil))
			return
		}

		claims, err := auth.ParseAccessToken(strings.TrimSpace(parts[1]), cfg.JWTSecret)
		if err != nil {
			apperrors.Abort(c, apperrors.New("UNAUTHORIZED", "Invalid or expired token", http.StatusUnauthorized, nil))
			return
		}

		c.Set(UserIDKey, claims.Subject)
		c.Set(RoleKey, claims.Role)
		c.Set(tenant.TenantIDKey, claims.TenantID)
		c.Next()
	}
}
