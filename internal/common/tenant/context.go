package tenant

import (
	"errors"

	"github.com/gin-gonic/gin"
)

const TenantIDKey = "tenantId"

func GetTenantID(c *gin.Context) (string, error) {
	raw, ok := c.Get(TenantIDKey)
	if !ok {
		return "", errors.New("tenantId not found in context")
	}
	tenantID, ok := raw.(string)
	if !ok || tenantID == "" {
		return "", errors.New("invalid tenantId value")
	}
	return tenantID, nil
}
