package tenant

import (
	"errors"

	"github.com/gin-gonic/gin"
)

const TenantIDKey = "tenantId"

func GetTenantID(c *gin.Context) (string, error) {
	// 1. Ambil tenantID yang sah dari Token JWT (diset oleh middleware)
	contextTenantID, exists := c.Get(TenantIDKey)
	if !exists {
		return "", errors.New("tenantId tidak ditemukan di dalam token/konteks")
	}

	authorizedID, ok := contextTenantID.(string)
	if !ok || authorizedID == "" {
		return "", errors.New("tenantId di dalam konteks tidak valid")
	}

	// 2. Cek apakah ada input manual dari Header atau Query
	manualID := c.GetHeader("X-Tenant-ID")
	if manualID == "" {
		manualID = c.Query("tenant_id")
	}

	// 3. Jika ada input manual, validasi apakah sama dengan yang di Token
	if manualID != "" && manualID != authorizedID {
		return "", errors.New("tenantId yang Anda kirimkan tidak sesuai dengan hak akses token Anda")
	}

	return authorizedID, nil
}

func GetID(c *gin.Context) (string, error) {
	return GetTenantID(c)
}
