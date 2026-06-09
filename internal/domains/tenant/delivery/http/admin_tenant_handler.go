package http

import (
	"net/http"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/common/response"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant/dto"
	"github.com/gin-gonic/gin"
)

type AdminTenantHandler struct {
	usecase domain.AdminTenantUsecase
}

func NewAdminTenantHandler(usecase domain.AdminTenantUsecase) *AdminTenantHandler {
	return &AdminTenantHandler{usecase: usecase}
}

// CreateAdminTenant godoc
// @Summary Create admin tenant
// @Description Admin membuat tenant baru untuk partner tertentu
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateAdminTenantRequest true "Create admin tenant payload"
// @Success 201 {object} response.Envelope
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /admin/tenants [post]
func (h *AdminTenantHandler) CreateAdminTenant(c *gin.Context) {
	var req dto.CreateAdminTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Validation failed", http.StatusUnprocessableEntity))
		return
	}

	res, err := h.usecase.CreateAdminTenant(c.Request.Context(), req)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Tenant berhasil dibuat oleh admin", res)
}

// ListAllTenants godoc
// @Summary List all tenants
// @Description Admin mengambil daftar semua tenant yang ada di database
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /admin/tenants [get]
func (h *AdminTenantHandler) ListAllTenants(c *gin.Context) {
	var req dto.ListAllTenantsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Query parameter tidak valid", http.StatusUnprocessableEntity))
		return
	}

	res, err := h.usecase.ListAllTenants(c.Request.Context(), req)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Semua tenant berhasil diambil oleh admin", res)
}

// SoftDeleteTenant godoc
// @Summary Soft delete tenant
// @Description Admin men-soft delete tenant berdasarkan ID
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param id path string true "Tenant ID"
// @Success 200 {object} response.Envelope
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /admin/tenants/{id} [delete]
func (h *AdminTenantHandler) SoftDeleteTenant(c *gin.Context) {
	tenantID := c.Param("id")
	err := h.usecase.SoftDeleteTenant(c.Request.Context(), tenantID)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Tenant berhasil di delete oleh admin", nil)
}

// GetTenantsByUserID godoc
// @Summary Get tenants by user ID
// @Description Admin mengambil daftar tenant berdasarkan User ID partner
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User ID Partner"
// @Success 200 {object} response.Envelope
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /admin/tenants/users/{userId} [get]
func (h *AdminTenantHandler) GetTenantsByUserID(c *gin.Context) {
	userID := c.Param("userId")
	res, err := h.usecase.GetTenantsByUserID(c.Request.Context(), userID)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Tenant berdasarkan user ID berhasil diambil oleh admin", res)
}

// GetTenantByID godoc
// @Summary Get tenant detail by ID
// @Description Admin mengambil detail satu tenant berdasarkan ID
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param id path string true "Tenant ID"
// @Success 200 {object} response.Envelope{data=dto.AdminTenantResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /admin/tenants/{id} [get]
func (h *AdminTenantHandler) GetTenantByID(c *gin.Context) {
	tenantID := c.Param("id")
	res, err := h.usecase.GetTenantByID(c.Request.Context(), tenantID)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Detail tenant berhasil diambil oleh admin", res)
}
