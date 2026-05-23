package tenant

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

const TenantIDKey = "tenantId"

func GetTenantID(c *gin.Context) (string, error) {
	// 1. Ambil role dari konteks
	roleVal, _ := c.Get("role")
	roleStr, _ := roleVal.(string)

	// 2. Cek apakah ada input manual dari Header atau Query
	manualID := c.GetHeader("X-Tenant-ID")
	if manualID == "" {
		manualID = c.Query("tenant_id")
	}
	if manualID == "" {
		manualID = c.Query("tenantId")
	}

	// 3. Jika role adalah PARTNER atau ADMIN, mereka boleh mengakses tenant mana pun
	// (validasi kepemilikan tenant untuk partner dilakukan di TenantGuard middleware)
	if strings.EqualFold(roleStr, "PARTNER") || strings.EqualFold(roleStr, "ADMIN") {
		if manualID != "" {
			return manualID, nil
		}
		// Fallback ke tenantId dari token jika ada
		contextTenantID, exists := c.Get(TenantIDKey)
		if exists {
			if authorizedID, ok := contextTenantID.(string); ok && authorizedID != "" {
				return authorizedID, nil
			}
		}
		return "", errors.New("tenantId tidak ditemukan untuk partner/admin")
	}

	// 4. Default logic untuk CUSTOMER / BASIC / dll. (harus sama dengan token)
	contextTenantID, exists := c.Get(TenantIDKey)
	if !exists {
		return "", errors.New("tenantId tidak ditemukan di dalam token/konteks")
	}

	authorizedID, ok := contextTenantID.(string)
	if !ok || authorizedID == "" {
		return "", errors.New("tenantId di dalam konteks tidak valid")
	}

	if manualID != "" && manualID != authorizedID {
		return "", errors.New("tenantId yang Anda kirimkan tidak sesuai dengan hak akses token Anda")
	}

	return authorizedID, nil
}

func GetID(c *gin.Context) (string, error) {
	return GetTenantID(c)
}
