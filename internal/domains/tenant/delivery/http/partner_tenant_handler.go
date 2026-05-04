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

type PartnerTenantHandler struct {
	usecase domain.PartnerTenantUsecase
}

func NewPartnerTenantHandler(usecase domain.PartnerTenantUsecase) *PartnerTenantHandler {
	return &PartnerTenantHandler{usecase: usecase}
}

// CreatePartnerTenant godoc
// @Summary Create partner tenant
// @Description Partner membuat tenant baru miliknya sesuai limit paket subscription
// @Tags Partner
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreatePartnerTenantRequest true "Create partner tenant payload"
// @Success 201 {object} response.Envelope
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /partner/tenants [post]
func (h *PartnerTenantHandler) CreatePartnerTenant(c *gin.Context) {
	userID, _ := c.Get("userId")
	userIDStr, _ := userID.(string)

	var req dto.CreatePartnerTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var validationErrs validator.ValidationErrors
		details := err.Error()
		if errors.As(err, &validationErrs) {
			details = validationErrs.Error()
		}
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Payload create tenant tidak valid", http.StatusBadRequest, details))
		return
	}

	res, err := h.usecase.CreatePartnerTenant(c.Request.Context(), userIDStr, req)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Tenant partner berhasil dibuat", res)
}

// ListPartnerTenants godoc
// @Summary List partner tenants
// @Description Ambil daftar tenant milik partner login
// @Tags Partner
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /partner/tenants [get]
func (h *PartnerTenantHandler) ListPartnerTenants(c *gin.Context) {
	userID, _ := c.Get("userId")
	userIDStr, _ := userID.(string)

	res, err := h.usecase.ListPartnerTenants(c.Request.Context(), userIDStr)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Daftar tenant partner berhasil diambil", res)
}
