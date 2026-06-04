package middleware

import (
	"net/http"
	"time"

	"github.com/dhegas/saas_gangsta/internal/common/cache"
	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant/domain"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var tenantCache = cache.NewLocalCache()

// TenantResolver intercepts public tenant requests, resolves the slug into a valid Tenant ID, and injects context.
func TenantResolver(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		slug := c.Param("tenantSlug")
		if slug == "" {
			apperrors.Abort(c, apperrors.New("VALIDATION_ERROR", "Slug tenant wajib disertakan", http.StatusBadRequest, nil))
			return
		}

		cacheKey := "resolved_tenant:slug:" + slug
		if cached, found := tenantCache.Get(cacheKey); found {
			if tenant, ok := cached.(domain.PublicTenant); ok {
				c.Set("tenantId", tenant.ID)
				c.Set("tenantID", tenant.ID)
				c.Set("tenantSlug", tenant.Slug)
				c.Set("tenantName", tenant.Name)
				c.Next()
				return
			}
		}

		var tenant domain.PublicTenant
		err := db.WithContext(c.Request.Context()).Table("tenants").
			Select("id::text AS id, name, slug, status, is_public").
			Where("slug = ?", slug).
			Where("status = 'active'").
			Where("is_public = true").
			Where("deleted_at IS NULL").
			Scan(&tenant).Error

		if err != nil {
			apperrors.Abort(c, apperrors.New("INTERNAL_ERROR", "Gagal memproses validasi tenant", http.StatusInternalServerError, nil))
			return
		}

		if tenant.ID == "" {
			apperrors.Abort(c, apperrors.New("NOT_FOUND", "Tenant tidak ditemukan atau tidak aktif", http.StatusNotFound, nil))
			return
		}

		// Cache for 10 minutes
		tenantCache.Set(cacheKey, tenant, 10*time.Minute)

		// Inject tenant context into gin.Context
		c.Set("tenantId", tenant.ID)
		c.Set("tenantID", tenant.ID)
		c.Set("tenantSlug", tenant.Slug)
		c.Set("tenantName", tenant.Name)

		c.Next()
	}
}

// TenantResolverMiddleware resolves the tenant slug, validates it, and injects context (alias for TenantResolver).
func TenantResolverMiddleware(db *gorm.DB) gin.HandlerFunc {
	return TenantResolver(db)
}

