package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	// Pastikan path ini sesuai dengan folder kamu yang menggunakan 'adminTenant'
	"github.com/dhegas/saas_gangsta/internal/modules/adminTenant/domain"
	"github.com/dhegas/saas_gangsta/internal/modules/adminTenant/dto"
)

type TenantHandler struct {
	usecase domain.AdminTenantUsecase
}

// NewTenantHandler adalah constructor untuk handler ini
func NewTenantHandler(usecase domain.AdminTenantUsecase) *TenantHandler {
	return &TenantHandler{usecase: usecase}
}

// GetAllTenants godoc
// @Summary      Get All Tenants
// @Description  Mengambil daftar seluruh toko/merchant (tenant) yang terdaftar di platform
// @Tags         Admin Tenant
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /admin/tenants [get]
func (h *TenantHandler) GetAllTenants(c *gin.Context) {
	ctx := c.Request.Context()

	tenants, err := h.usecase.GetAllTenants(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data tenant",
			"error":   gin.H{"code": "INTERNAL_ERROR", "details": err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data tenant berhasil diambil",
		"data":    tenants,
	})
}

// UpdateTenantStatus godoc
// @Summary      Update Tenant Status
// @Description  Mengubah status operasional tenant (active, inactive, atau suspended)
// @Tags         Admin Tenant
// @Accept       json
// @Produce      json
// @Param        id       path      string                             true  "Tenant ID"
// @Param        request  body      dto.UpdateTenantStatusRequest      true  "Payload Status Baru"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /admin/tenants/{id}/status [patch]
func (h *TenantHandler) UpdateTenantStatus(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := c.Param("id")

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

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Status tenant berhasil diperbarui",
		"data":    nil,
	})
}

// RegisterRoutes digunakan untuk mendaftarkan endpoint ke dalam router Gin
func (h *TenantHandler) RegisterRoutes(router *gin.RouterGroup) {
	adminRoute := router.Group("/admin")
	{
		adminRoute.GET("/tenants", h.GetAllTenants)
		adminRoute.PATCH("/tenants/:id/status", h.UpdateTenantStatus)
	}
}
