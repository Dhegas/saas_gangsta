package http

import (
	"errors"
	"net/http"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/common/response"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant/dto"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
		var validationErrs validator.ValidationErrors
		details := err.Error()
		if errors.As(err, &validationErrs) {
			details = validationErrs.Error()
		}
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Payload create tenant tidak valid", http.StatusBadRequest, details))
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
		var validationErrs validator.ValidationErrors
		details := err.Error()
		if errors.As(err, &validationErrs) {
			details = validationErrs.Error()
		}
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Query parameter tidak valid", http.StatusBadRequest, details))
		return
	}

	res, err := h.usecase.ListAllTenants(c.Request.Context(), req)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Semua tenant berhasil diambil oleh admin", res)
}
