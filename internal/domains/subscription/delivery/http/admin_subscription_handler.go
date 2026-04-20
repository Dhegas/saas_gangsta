package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	// Sesuaikan path import jika perlu
	"github.com/dhegas/saas_gangsta/internal/domains/subscription/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/subscription/dto"
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
// @Router       /api/v1/admin/subscriptions/plans [get]
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

// CreatePlan godoc
// @Summary      Create Subscription Plan
// @Description  Membuat paket langganan baru
// @Tags         Admin Subscription
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateSubscriptionPlanRequest  true  "Payload Data Paket"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /api/v1/admin/subscriptions/plans [post]
func (h *SubscriptionHandler) CreatePlan(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.CreateSubscriptionPlanRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Input tidak valid", "error": err.Error()})
		return
	}

	if err := h.usecase.CreatePlan(ctx, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal membuat paket", "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Paket berhasil dibuat"})
}

// UpdatePlan godoc
// @Summary      Update Subscription Plan
// @Description  Mengubah data paket langganan (termasuk menonaktifkan dengan isActive: false)
// @Tags         Admin Subscription
// @Accept       json
// @Produce      json
// @Param        id       path      string                             true  "Plan ID"
// @Param        request  body      dto.UpdateSubscriptionPlanRequest  true  "Payload Data Update"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /api/v1/admin/subscriptions/plans/{id} [patch]
func (h *SubscriptionHandler) UpdatePlan(c *gin.Context) {
	ctx := c.Request.Context()
	planID := c.Param("id")
	var req dto.UpdateSubscriptionPlanRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Input tidak valid", "error": err.Error()})
		return
	}

	if err := h.usecase.UpdatePlan(ctx, planID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Gagal update paket", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Paket berhasil diupdate"})
}

// RegisterRoutes mendaftarkan endpoint ke dalam group yang sudah diterima.
// Group yang masuk sudah ber-prefix /admin (dari bootstrap/admin_routes.go).
func (h *SubscriptionHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/subscriptions/plans", h.GetAllPlans)
	router.POST("/subscriptions/plans", h.CreatePlan)
	router.PATCH("/subscriptions/plans/:id", h.UpdatePlan)
}
