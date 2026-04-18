package tenant

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

const TenantIDKey = "tenantId"

const (
	TenantHeaderKey    = "X-Tenant-Id"
	TenantHeaderAltKey = "X-Tenant-ID"
)

func GetTenantID(c *gin.Context) (string, error) {
	raw, ok := c.Get(TenantIDKey)
	if ok {
		tenantID, ok := raw.(string)
		if ok && tenantID != "" {
			return tenantID, nil
		}
	}

	headerTenantID := strings.TrimSpace(c.GetHeader(TenantHeaderKey))
	if headerTenantID == "" {
		headerTenantID = strings.TrimSpace(c.GetHeader(TenantHeaderAltKey))
	}
	if headerTenantID != "" {
		return headerTenantID, nil
	}

	return "", errors.New("tenantId not found in context or header")
}
