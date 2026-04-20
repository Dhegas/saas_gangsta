package http

import (
	"errors"
	"net/http"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/common/response"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant/dto"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// TenantProfileHandler menangani endpoint tenant profile.
type TenantProfileHandler struct {
	usecase domain.TenantProfileUsecase
}

func NewTenantProfileHandler(usecase domain.TenantProfileUsecase) *TenantProfileHandler {
	return &TenantProfileHandler{usecase: usecase}
}

func (h *TenantProfileHandler) Create(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		apperrors.Write(c, apperrors.New("TENANT_NOT_FOUND", "Tenant context is required", http.StatusUnauthorized, nil))
		return
	}

	var req dto.CreateTenantProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var validationErrs validator.ValidationErrors
		details := err.Error()
		if errors.As(err, &validationErrs) {
			details = validationErrs.Error()
		}
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Payload tenant profile tidak valid", http.StatusBadRequest, details))
		return
	}

	res, err := h.usecase.CreateProfile(c.Request.Context(), tenantID, req)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Tenant profile berhasil dibuat", res)
}

func (h *TenantProfileHandler) List(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		apperrors.Write(c, apperrors.New("TENANT_NOT_FOUND", "Tenant context is required", http.StatusUnauthorized, nil))
		return
	}

	res, err := h.usecase.ListProfiles(c.Request.Context(), tenantID)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Daftar tenant profile berhasil diambil", res)
}
