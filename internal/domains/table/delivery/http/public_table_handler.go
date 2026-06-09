package http

import (
	"net/http"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/common/response"
	"github.com/dhegas/saas_gangsta/internal/domains/table/domain"
	"github.com/gin-gonic/gin"
)

type PublicTableHandler struct {
	usecase domain.PublicTableUsecase
}

func NewPublicTableHandler(usecase domain.PublicTableUsecase) *PublicTableHandler {
	return &PublicTableHandler{usecase: usecase}
}

// GetPublicTables godoc
// @Summary List active dining tables for a tenant
// @Description Mengambil daftar meja aktif dalam context tenant beserta status keterisiannya
// @Tags Public
// @Produce json
// @Param tenantSlug path string true "Tenant Slug"
// @Success 200 {object} response.Envelope{data=[]dto.PublicTableResponse}
// @Failure 400 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /public/tenant/{tenantSlug}/tables [get]
func (h *PublicTableHandler) GetPublicTables(c *gin.Context) {
	tenantID := c.GetString("tenantId")
	if tenantID == "" {
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Tenant context tidak ditemukan", http.StatusBadRequest))
		return
	}

	tables, err := h.usecase.GetPublicTables(c.Request.Context(), tenantID)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Table list fetched successfully", tables)
}
