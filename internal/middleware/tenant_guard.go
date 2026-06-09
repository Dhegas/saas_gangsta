package middleware

import (
	"net/http"
	"strings"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TenantGuard(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, err := tenant.GetTenantID(c)
		if err != nil {
			apperrors.Abort(c, apperrors.New("TENANT_NOT_FOUND", "Tenant context is required", http.StatusUnauthorized))
			return
		}

		// Ambil role dan userID dari konteks
		roleVal, _ := c.Get(RoleKey)
		roleStr, _ := roleVal.(string)

		// Jika user adalah PARTNER, validasi apakah ia memiliki hak akses (ownership) ke tenant tersebut
		if strings.EqualFold(roleStr, "PARTNER") && db != nil {
			userIDVal, _ := c.Get(UserIDKey)
			userIDStr, _ := userIDVal.(string)

			var count int64
			err := db.WithContext(c.Request.Context()).Table("tenants").
				Where("id = NULLIF(?, '')::uuid AND user_id = NULLIF(?, '')::uuid AND deleted_at IS NULL", tenantID, userIDStr).
				Count(&count).Error

			if err != nil || count == 0 {
				apperrors.Abort(c, apperrors.New("FORBIDDEN", "Anda tidak memiliki akses ke tenant ini", http.StatusForbidden))
				return
			}
		}

		// Set tenantId ke konteks agar handler berikutnya bisa mengambilnya
		c.Set(tenant.TenantIDKey, tenantID)
		c.Set("tenantID", tenantID)

		c.Next()
	}
}
