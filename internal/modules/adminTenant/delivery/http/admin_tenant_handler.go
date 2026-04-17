package http

import (
	"net/http"

	"github.com/dhegas/saas_gangsta/internal/modules/adminTenant/domain"
	"github.com/dhegas/saas_gangsta/internal/modules/adminTenant/dto"
	"github.com/gin-gonic/gin"
)

type TenantHandler struct {
	usecase domain.AdminTenantUsecase
}

// Constructor
func NewTenantHandler(usecase domain.AdminTenantUsecase) *TenantHandler {
	return &TenantHandler{usecase: usecase}
}

// API: GET /api/v1/admin/tenants
func (h *TenantHandler) GetAllTenants(c *gin.Context) {
	ctx := c.Request.Context()

	tenants, err := h.usecase.GetAllTenants(ctx)
	if err != nil {
		// Asumsi kamu akan membuat helper standard error response nanti
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data tenant",
			"error":   gin.H{"code": "INTERNAL_ERROR", "details": err.Error()},
		})
		return
	}

	// Menggunakan Standard Response Envelope sesuai README
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data tenant berhasil diambil",
		"data":    tenants,
	})
}

// API: PATCH /api/v1/admin/tenants/:id/status
func (h *TenantHandler) UpdateTenantStatus(c *gin.Context) {
	ctx := c.Request.Context()

	// Mengambil ID tenant dari parameter URL (misal: /tenants/123e4567-e89b-12d3-a456-426614174000/status)
	tenantID := c.Param("id")

	// 1. Validasi Payload Request
	var req dto.UpdateTenantStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Input tidak valid",
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"details": err.Error(),
			},
		})
		return
	}

	// 2. Panggil Usecase untuk mengeksekusi logika bisnis
	err := h.usecase.UpdateTenantStatus(ctx, tenantID, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memperbarui status tenant",
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"details": err.Error(),
			},
		})
		return
	}

	// 3. Kembalikan Response Sukses
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Status tenant berhasil diperbarui",
		"data":    nil, // Sesuai standar, data bisa null jika tidak ada objek yang perlu dikembalikan
	})
}

// RegisterRoutes digunakan untuk mendaftarkan endpoint ke dalam router Gin
func (h *TenantHandler) RegisterRoutes(router *gin.RouterGroup) {
	// Endpoint ini asumsikan sudah dilewati middleware JWT dan RoleGuard("admin") di level bootstrap
	adminRoute := router.Group("/admin")
	{
		adminRoute.GET("/tenants", h.GetAllTenants)
		adminRoute.PATCH("/tenants/:id/status", h.UpdateTenantStatus)
	}
}
