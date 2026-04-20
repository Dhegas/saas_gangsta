package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

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
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /admin/dashboard [get]
func (h *DashboardHandler) GetDashboardStats(c *gin.Context) {
	ctx := c.Request.Context()

	stats, err := h.usecase.GetStats(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data statistik dashboard",
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data statistik dashboard berhasil diambil",
		"data":    stats,
	})
}

// RegisterRoutes mendaftarkan endpoint ke dalam group yang sudah diterima.
// Group yang masuk sudah ber-prefix /admin (dari bootstrap/admin_routes.go).
func (h *DashboardHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/dashboard", h.GetDashboardStats)
}
