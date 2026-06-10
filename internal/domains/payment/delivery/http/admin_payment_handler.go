package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/common/response"
	paymentdomain "github.com/dhegas/saas_gangsta/internal/domains/payment/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/payment/dto"
)

// AdminPaymentHandler menangani sinkronisasi pembayaran oleh admin
type AdminPaymentHandler struct {
	usecase paymentdomain.PaymentWebhookUsecase
}

func NewAdminPaymentHandler(usecase paymentdomain.PaymentWebhookUsecase) *AdminPaymentHandler {
	return &AdminPaymentHandler{usecase: usecase}
}

// SyncPaymentStatus godoc
// @Summary      Sinkronisasi Status Pembayaran dengan Midtrans (Admin)
// @Description  Sinkronisasi status pembayaran lokal dengan Midtrans. Jika transaksi di Midtrans sudah sukses tetapi lokal belum terupdate, perbaiki data secara aman.
// @Tags         Admin Payment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body  dto.CreateSnapRequest  true  "Order ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      422  {object}  map[string]interface{}
// @Router       /admin/payments/sync [post]
func (h *AdminPaymentHandler) SyncPaymentStatus(c *gin.Context) {
	var req dto.CreateSnapRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "order_id diperlukan", http.StatusBadRequest))
		return
	}

	err := h.usecase.SyncPaymentStatus(c.Request.Context(), req.OrderID)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Status pembayaran berhasil disinkronkan", nil)
}

// RegisterRoutes mendaftarkan endpoint admin payment
func (h *AdminPaymentHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/payments/sync", h.SyncPaymentStatus)
}
