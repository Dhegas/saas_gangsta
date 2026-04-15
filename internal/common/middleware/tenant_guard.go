package middleware

import (
	"net/http"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/common/tenant"
	"github.com/gin-gonic/gin"
)

func TenantGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := tenant.GetTenantID(c); err != nil {
			apperrors.Abort(c, apperrors.New("TENANT_NOT_FOUND", "Tenant context is required", http.StatusUnauthorized, nil))
			return
		}
		c.Next()
	}
}
