package http

import (
	"net/http"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/common/response"
	"github.com/dhegas/saas_gangsta/internal/domains/menu/domain"
	"github.com/gin-gonic/gin"
)

type PublicMenuHandler struct {
	usecase domain.PublicMenuUsecase
}

func NewPublicMenuHandler(usecase domain.PublicMenuUsecase) *PublicMenuHandler {
	return &PublicMenuHandler{usecase: usecase}
}

// GetPublicMenus godoc
// @Summary List active menus for a tenant
// @Description Mengambil daftar menu aktif dalam context tenant dengan filter pencarian dan kategori
// @Tags Public
// @Produce json
// @Param tenantSlug path string true "Tenant Slug"
// @Param categoryId query string false "Filter berdasarkan Category ID"
// @Param search query string false "Pencarian nama atau deskripsi menu"
// @Param isAvailable query bool false "Filter ketersediaan menu"
// @Success 200 {object} response.Envelope{data=[]dto.MenuResponse}
// @Failure 400 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /public/tenant/{tenantSlug}/menus [get]
func (h *PublicMenuHandler) GetPublicMenus(c *gin.Context) {
	tenantID := c.GetString("tenantId")
	if tenantID == "" {
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Tenant context tidak ditemukan", http.StatusBadRequest))
		return
	}

	categoryID := c.Query("categoryId")
	search := c.Query("search")
	isAvailableStr := c.Query("isAvailable")

	var isAvailable *bool
	if isAvailableStr != "" {
		val := isAvailableStr == "true"
		isAvailable = &val
	}

	menus, err := h.usecase.GetPublicMenus(c.Request.Context(), tenantID, categoryID, search, isAvailable)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Menu list fetched successfully", menus)
}
