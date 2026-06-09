package http

import (
	"net/http"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/common/response"
	"github.com/dhegas/saas_gangsta/internal/domains/report/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/report/dto"
	"github.com/gin-gonic/gin"
)

type PartnerReportHandler struct {
	usecase domain.PartnerReportUsecase
}

func NewPartnerReportHandler(usecase domain.PartnerReportUsecase) *PartnerReportHandler {
	return &PartnerReportHandler{usecase: usecase}
}

// extractReportTenantID mendapatkan tenant_id dari gin context (sudah diset oleh TenantGuard)
func extractReportTenantID(c *gin.Context) (string, bool) {
	val, exists := c.Get("tenantId")
	if !exists {
		apperrors.Write(c, apperrors.New("TENANT_NOT_FOUND", "Tenant ID diperlukan", http.StatusBadRequest))
		return "", false
	}
	tenantID, ok := val.(string)
	if !ok || tenantID == "" {
		apperrors.Write(c, apperrors.New("TENANT_NOT_FOUND", "Tenant ID tidak valid", http.StatusBadRequest))
		return "", false
	}
	return tenantID, true
}

// GetRevenue godoc
// @Summary      Total Pendapatan (Revenue)
// @Description  Menghitung total pendapatan dari order COMPLETED dalam rentang tanggal yang ditentukan.
// @Tags         Reporting & Analytics
// @Produce      json
// @Security     BearerAuth
// @Param        from  query     string  true  "Tanggal awal (YYYY-MM-DD)"
// @Param        to    query     string  true  "Tanggal akhir (YYYY-MM-DD)"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /reports/revenue [get]
func (h *PartnerReportHandler) GetRevenue(c *gin.Context) {
	tenantID, ok := extractReportTenantID(c)
	if !ok {
		return
	}

	var params dto.RevenueFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Parameter tidak valid", http.StatusUnprocessableEntity))
		return
	}

	result, err := h.usecase.GetRevenue(c.Request.Context(), tenantID, params)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil data revenue", result)
}

// GetTopMenus godoc
// @Summary      Menu Terlaris
// @Description  Mengambil daftar menu terlaris berdasarkan total qty terjual dari order COMPLETED.
// @Tags         Reporting & Analytics
// @Produce      json
// @Security     BearerAuth
// @Param        from   query     string  false  "Tanggal awal (YYYY-MM-DD)"
// @Param        to     query     string  false  "Tanggal akhir (YYYY-MM-DD)"
// @Param        limit  query     int     false  "Jumlah item (default: 10, max: 100)"
// @Success      200    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /reports/top-menus [get]
func (h *PartnerReportHandler) GetTopMenus(c *gin.Context) {
	tenantID, ok := extractReportTenantID(c)
	if !ok {
		return
	}

	var params dto.TopMenusFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Parameter tidak valid", http.StatusUnprocessableEntity))
		return
	}

	result, err := h.usecase.GetTopMenus(c.Request.Context(), tenantID, params)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil data menu terlaris", result)
}

// GetOrdersByTable godoc
// @Summary      Order Terbanyak per Meja
// @Description  Mengambil daftar meja dengan jumlah order terbanyak dari order COMPLETED.
// @Tags         Reporting & Analytics
// @Produce      json
// @Security     BearerAuth
// @Param        from   query     string  false  "Tanggal awal (YYYY-MM-DD)"
// @Param        to     query     string  false  "Tanggal akhir (YYYY-MM-DD)"
// @Param        limit  query     int     false  "Jumlah item (default: 10, max: 100)"
// @Success      200    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /reports/orders-by-table [get]
func (h *PartnerReportHandler) GetOrdersByTable(c *gin.Context) {
	tenantID, ok := extractReportTenantID(c)
	if !ok {
		return
	}

	var params dto.OrdersByTableFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Parameter tidak valid", http.StatusUnprocessableEntity))
		return
	}

	result, err := h.usecase.GetOrdersByTable(c.Request.Context(), tenantID, params)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil data order per meja", result)
}

// GetDailySummary godoc
// @Summary      Ringkasan Harian
// @Description  Mengambil ringkasan order dan revenue per hari. Default 7 hari ke belakang jika from/to tidak diisi.
// @Tags         Reporting & Analytics
// @Produce      json
// @Security     BearerAuth
// @Param        from  query     string  false  "Tanggal awal (YYYY-MM-DD, default: 7 hari lalu)"
// @Param        to    query     string  false  "Tanggal akhir (YYYY-MM-DD, default: hari ini)"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /reports/daily-summary [get]
func (h *PartnerReportHandler) GetDailySummary(c *gin.Context) {
	tenantID, ok := extractReportTenantID(c)
	if !ok {
		return
	}

	var params dto.DailySummaryFilterParams
	if err := c.ShouldBindQuery(&params); err != nil {
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Parameter tidak valid", http.StatusUnprocessableEntity))
		return
	}

	result, err := h.usecase.GetDailySummary(c.Request.Context(), tenantID, params)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil ringkasan harian", result)
}
