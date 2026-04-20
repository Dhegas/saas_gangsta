package tenant

import (
	"errors"

	"github.com/gin-gonic/gin"
)

const TenantIDKey = "tenantId"

func GetTenantID(c *gin.Context) (string, error) {
	tenantID, exists := c.Get(TenantIDKey)
	if !exists {
		return "", errors.New("tenantId not found in context")
	}

	tenantIDStr, ok := tenantID.(string)
	if !ok || tenantIDStr == "" {
		return "", errors.New("tenantId in context is invalid")
	}

	return tenantIDStr, nil
}

func GetID(c *gin.Context) (string, error) {
	return GetTenantID(c)
}
