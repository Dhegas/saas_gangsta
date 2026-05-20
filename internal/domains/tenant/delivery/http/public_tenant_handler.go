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

type PublicTenantHandler struct {
	usecase domain.PublicTenantUsecase
}

func NewPublicTenantHandler(usecase domain.PublicTenantUsecase) *PublicTenantHandler {
	return &PublicTenantHandler{usecase: usecase}
}

// GetPublicTenantList godoc
// @Summary List public tenants
// @Description Explore public merchants/tenants with search & pagination
// @Tags Public
// @Produce json
// @Param search query string false "Search query for name, description, address"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} response.Envelope{data=[]dto.PublicTenantResponse,meta=dto.PaginationMeta}
// @Failure 400 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /public/tenants [get]
func (h *PublicTenantHandler) GetPublicTenantList(c *gin.Context) {
	var req dto.ListPublicTenantsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		var validationErrs validator.ValidationErrors
		details := err.Error()
		if errors.As(err, &validationErrs) {
			details = validationErrs.Error()
		}
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Query params tidak valid", http.StatusBadRequest, details))
		return
	}

	res, err := h.usecase.ListPublicTenants(c.Request.Context(), req)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.SuccessWithMeta(c, http.StatusOK, "Tenant list fetched successfully", res.Data, res.Meta)
}

// GetPublicTenantDetail godoc
// @Summary Get public tenant detail
// @Description Mengambil detail satu tenant berdasarkan slug untuk branding & context toko
// @Tags Public
// @Produce json
// @Param slug path string true "Tenant Slug"
// @Success 200 {object} response.Envelope{data=dto.PublicTenantDetailResponse}
// @Failure 400 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /public/tenants/{slug} [get]
func (h *PublicTenantHandler) GetPublicTenantDetail(c *gin.Context) {
	slug := c.Param("slug")
	res, err := h.usecase.GetPublicTenantBySlug(c.Request.Context(), slug)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Tenant detail fetched successfully", res)
}
