package middleware

import (
	"net/http"
	"strings"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/gin-gonic/gin"
)

func RoleGuard(allowedRole string) gin.HandlerFunc {
	return RoleGuards(allowedRole)
}

func RoleGuards(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(allowedRoles) == 0 {
			apperrors.Abort(c, apperrors.New("FORBIDDEN", "Role guard is not configured", http.StatusForbidden, nil))
			return
		}

		role, ok := c.Get(RoleKey)
		roleStr, okCast := role.(string)

		if !ok || !okCast || roleStr == "" {
			apperrors.Abort(c, apperrors.New("FORBIDDEN", "Role not found", http.StatusForbidden, nil))
			return
		}

		allowed := false
		for _, allowedRole := range allowedRoles {
			if strings.EqualFold(roleStr, allowedRole) {
				allowed = true
				break
			}
		}
		if !allowed {
			apperrors.Abort(c, apperrors.New("FORBIDDEN", "You do not have access to this resource", http.StatusForbidden, nil))
			return
		}
		c.Next()
	}
}
