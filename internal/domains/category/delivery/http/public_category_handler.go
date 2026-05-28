package http

import (
	"net/http"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/common/response"
	"github.com/dhegas/saas_gangsta/internal/domains/category/domain"
	"github.com/gin-gonic/gin"
)

type PublicCategoryHandler struct {
	usecase domain.PublicCategoryUsecase
}

func NewPublicCategoryHandler(usecase domain.PublicCategoryUsecase) *PublicCategoryHandler {
	return &PublicCategoryHandler{usecase: usecase}
}

// GetPublicCategories godoc
// @Summary List active menu categories for a tenant
// @Description Mengambil daftar kategori menu aktif dalam context tenant
// @Tags Public
// @Produce json
// @Param tenantSlug path string true "Tenant Slug"
// @Success 200 {object} response.Envelope{data=[]dto.CategoryResponse}
// @Failure 400 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /public/tenant/{tenantSlug}/categories [get]
func (h *PublicCategoryHandler) GetPublicCategories(c *gin.Context) {
	tenantID := c.GetString("tenantId")
	if tenantID == "" {
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Tenant context tidak ditemukan", http.StatusBadRequest, nil))
		return
	}

	categories, err := h.usecase.GetPublicCategories(c.Request.Context(), tenantID)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Category list fetched successfully", categories)
}
