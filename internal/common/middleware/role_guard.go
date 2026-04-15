package middleware

import (
	"net/http"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/gin-gonic/gin"
)

func RoleGuard(allowedRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get(RoleKey)
		roleStr, okCast := role.(string)
		if !ok || !okCast || roleStr == "" {
			apperrors.Abort(c, apperrors.New("FORBIDDEN", "Role not found", http.StatusForbidden, nil))
			return
		}
		if roleStr != allowedRole {
			apperrors.Abort(c, apperrors.New("FORBIDDEN", "You do not have access to this resource", http.StatusForbidden, nil))
			return
		}
		c.Next()
	}
}
