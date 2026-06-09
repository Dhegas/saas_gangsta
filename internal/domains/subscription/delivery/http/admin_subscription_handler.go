package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dhegas/saas_gangsta/internal/common/response"
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
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /admin/subscriptions/plans [get]
func (h *SubscriptionHandler) GetAllPlans(c *gin.Context) {
	ctx := c.Request.Context()

	plans, err := h.usecase.GetAllPlans(ctx)
	if err != nil {
		slog.Error("GetAllPlans failed",
			slog.String("error", err.Error()),
		)
		response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", gin.H{
			"code": "INTERNAL_ERROR",
		})
		return
	}

	response.Success(c, http.StatusOK, "Data paket langganan berhasil diambil", plans)
}

// CreatePlan godoc
// @Summary      Create Subscription Plan
// @Description  Membuat paket langganan baru
// @Tags         Admin Subscription
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.CreateSubscriptionPlanRequest  true  "Payload Data Paket"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /admin/subscriptions/plans [post]
func (h *SubscriptionHandler) CreatePlan(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.CreateSubscriptionPlanRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "Validation failed", gin.H{
			"code": "VALIDATION_ERROR",
		})
		return
	}

	if err := h.usecase.CreatePlan(ctx, req); err != nil {
		slog.Error("CreatePlan failed",
			slog.String("error", err.Error()),
		)
		response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", gin.H{
			"code": "INTERNAL_ERROR",
		})
		return
	}

	response.Success(c, http.StatusCreated, "Paket berhasil dibuat", nil)
}

// UpdatePlan godoc
// @Summary      Update Subscription Plan
// @Description  Mengubah data paket langganan (termasuk menonaktifkan dengan isActive: false)
// @Tags         Admin Subscription
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                             true  "Plan ID"
// @Param        request  body      dto.UpdateSubscriptionPlanRequest  true  "Payload Data Update"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /admin/subscriptions/plans/{id} [patch]
func (h *SubscriptionHandler) UpdatePlan(c *gin.Context) {
	ctx := c.Request.Context()
	planID := c.Param("id")
	var req dto.UpdateSubscriptionPlanRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusUnprocessableEntity, "Validation failed", gin.H{
			"code": "VALIDATION_ERROR",
		})
		return
	}

	if err := h.usecase.UpdatePlan(ctx, planID, req); err != nil {
		slog.Error("UpdatePlan failed",
			slog.String("plan_id", planID),
			slog.String("error", err.Error()),
		)
		response.Error(c, http.StatusInternalServerError, "An unexpected error occurred", gin.H{
			"code": "INTERNAL_ERROR",
		})
		return
	}

	response.Success(c, http.StatusOK, "Paket berhasil diupdate", nil)
}

// RegisterRoutes mendaftarkan endpoint ke dalam group yang sudah diterima.
// Group yang masuk sudah ber-prefix /admin (dari bootstrap/admin_routes.go).
func (h *SubscriptionHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/subscriptions/plans", h.GetAllPlans)
	router.POST("/subscriptions/plans", h.CreatePlan)
	router.PATCH("/subscriptions/plans/:id", h.UpdatePlan)
}
