package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/common/response"
	"github.com/dhegas/saas_gangsta/internal/domains/report/domain"
)

type DashboardHandler struct {
	usecase domain.AdminDashboardUsecase
}

// NewDashboardHandler adalah constructor untuk handler ini
func NewDashboardHandler(usecase domain.AdminDashboardUsecase) *DashboardHandler {
	return &DashboardHandler{usecase: usecase}
}

// GetDashboardStats godoc
// @Summary      Get Admin Dashboard Stats
// @Description  Mengambil overview metrik platform (total tenant, active subscriptions, monthly revenue)
// @Tags         Admin Dashboard
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /admin/dashboard [get]
// @Router       /admin/dashboard [get]
func (h *DashboardHandler) GetDashboardStats(c *gin.Context) {
	ctx := c.Request.Context()

	stats, err := h.usecase.GetStats(ctx)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Data statistik dashboard berhasil diambil", stats)
}

// RegisterRoutes mendaftarkan endpoint ke dalam group yang sudah diterima.
// Group yang masuk sudah ber-prefix /admin (dari bootstrap/admin_routes.go).
func (h *DashboardHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/dashboard", h.GetDashboardStats)
}
