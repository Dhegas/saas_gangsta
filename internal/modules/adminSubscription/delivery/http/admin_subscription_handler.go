package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	// Sesuaikan path import jika perlu
	"github.com/dhegas/saas_gangsta/internal/modules/adminSubscription/domain"
)

type SubscriptionHandler struct {
	usecase domain.AdminSubscriptionUsecase
}

// Constructor
func NewSubscriptionHandler(usecase domain.AdminSubscriptionUsecase) *SubscriptionHandler {
	return &SubscriptionHandler{usecase: usecase}
}

// GetAllPlans godoc
// @Summary      Get All Subscription Plans
// @Description  Mengambil daftar semua paket langganan (misal: Basic, Pro, Enterprise)
// @Tags         Admin Subscription
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /admin/subscriptions/plans [get]
func (h *SubscriptionHandler) GetAllPlans(c *gin.Context) {
	ctx := c.Request.Context()

	plans, err := h.usecase.GetAllPlans(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data paket langganan",
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"details": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data paket langganan berhasil diambil",
		"data":    plans,
	})
}

// RegisterRoutes mendaftarkan endpoint ke router Gin
func (h *SubscriptionHandler) RegisterRoutes(router *gin.RouterGroup) {
	adminRoute := router.Group("/admin")
	{
		// Endpoint: GET /api/v1/admin/subscriptions/plans
		adminRoute.GET("/subscriptions/plans", h.GetAllPlans)
	}
}
