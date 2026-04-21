package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/dhegas/saas_gangsta/internal/domains/tenant/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant/dto"
)

// TenantHandler memegang referensi ke usecase tenant
type TenantHandler struct {
	usecase domain.AdminTenantUsecase
}

// NewTenantHandler adalah constructor untuk dependency injection
func NewTenantHandler(usecase domain.AdminTenantUsecase) *TenantHandler {
	return &TenantHandler{usecase: usecase}
}

// ─── Helper ──────────────────────────────────────────────────────────────────

func errorResponse(c *gin.Context, status int, code, message string, detail interface{}) {
	c.JSON(status, gin.H{
		"success": false,
		"message": message,
		"error": gin.H{
			"code":    code,
			"details": detail,
		},
	})
}

// ─── Handlers ─────────────────────────────────────────────────────────────────

// GetAllTenants godoc
// @Summary      List All Tenants
// @Description  Mengambil daftar seluruh tenant yang terdaftar di platform
// @Tags         Admin Tenant
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /admin/tenants [get]
func (h *TenantHandler) GetAllTenants(c *gin.Context) {
	ctx := c.Request.Context()

	tenants, err := h.usecase.GetAllTenants(ctx)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal mengambil data tenant", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data tenant berhasil diambil",
		"data":    tenants,
	})
}

// GetTenantByID godoc
// @Summary      Get Tenant Detail
// @Description  Mengambil detail satu tenant berdasarkan ID
// @Tags         Admin Tenant
// @Produce      json
// @Security     BearerAuth
// @Param        id  path      string  true  "Tenant ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /admin/tenants/{id} [get]
func (h *TenantHandler) GetTenantByID(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := c.Param("id")

	tenant, err := h.usecase.GetTenantByID(ctx, tenantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse(c, http.StatusNotFound, "TENANT_NOT_FOUND", "Tenant tidak ditemukan", nil)
			return
		}
		errorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal mengambil detail tenant", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Detail tenant berhasil diambil",
		"data":    tenant,
	})
}

// CreateTenant godoc
// @Summary      Register New Tenant
// @Description  Mendaftarkan tenant (merchant/toko) baru ke dalam platform
// @Tags         Admin Tenant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.CreateTenantRequest  true  "Payload Tenant Baru"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /admin/tenants [post]
func (h *TenantHandler) CreateTenant(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Input tidak valid", err.Error())
		return
	}

	tenant, err := h.usecase.CreateTenant(ctx, req)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal membuat tenant baru", err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Tenant berhasil didaftarkan",
		"data":    tenant,
	})
}

// UpdateTenant godoc
// @Summary      Update Tenant
// @Description  Memperbarui data tenant (name, slug, atau status)
// @Tags         Admin Tenant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                   true  "Tenant ID"
// @Param        request  body      dto.UpdateTenantRequest  true  "Payload Update Tenant"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      404      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /admin/tenants/{id} [put]
func (h *TenantHandler) UpdateTenant(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := c.Param("id")

	var req dto.UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Input tidak valid", err.Error())
		return
	}

	tenant, err := h.usecase.UpdateTenant(ctx, tenantID, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse(c, http.StatusNotFound, "TENANT_NOT_FOUND", "Tenant tidak ditemukan", nil)
			return
		}
		errorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal memperbarui tenant", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data tenant berhasil diperbarui",
		"data":    tenant,
	})
}

// SoftDeleteTenant godoc
// @Summary      Soft Delete Tenant
// @Description  Menghapus tenant secara soft delete (mengisi deleted_at, data tetap ada di DB)
// @Tags         Admin Tenant
// @Produce      json
// @Security     BearerAuth
// @Param        id  path      string  true  "Tenant ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /admin/tenants/{id} [delete]
func (h *TenantHandler) SoftDeleteTenant(c *gin.Context) {
	ctx := c.Request.Context()
	tenantID := c.Param("id")

	if err := h.usecase.SoftDeleteTenant(ctx, tenantID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse(c, http.StatusNotFound, "TENANT_NOT_FOUND", "Tenant tidak ditemukan", nil)
			return
		}
		errorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal menghapus tenant", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tenant berhasil dihapus",
		"data":    nil,
	})
}

// UpdateTenantStatus godoc
// @Summary      Update Tenant Status
// @Description  Mengubah status operasional tenant (active, inactive, atau suspended)
// @Tags         Admin Tenant
// @Accept       json
// @Produce      json
// @Security     BearerAuth
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
		errorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Input tidak valid", err.Error())
		return
	}

	if err := h.usecase.UpdateTenantStatus(ctx, tenantID, req.Status); err != nil {
		errorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Gagal memperbarui status tenant", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Status tenant berhasil diperbarui",
		"data":    nil,
	})
}

// Catatan: Route registration untuk modul ini dilakukan di
// internal/bootstrap/admin_routes.go — bukan di sini.
// Lihat RegisterAdminRoutes() untuk daftar lengkap endpoint.
