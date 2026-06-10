package http

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/common/response"
	paymentdomain "github.com/dhegas/saas_gangsta/internal/domains/payment/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/payment/dto"
)

// WebhookHandler menangani notifikasi dari Midtrans
type WebhookHandler struct {
	usecase paymentdomain.PaymentWebhookUsecase
}

func NewWebhookHandler(usecase paymentdomain.PaymentWebhookUsecase) *WebhookHandler {
	return &WebhookHandler{usecase: usecase}
}

// HandleMidtrans godoc
// @Summary      Midtrans Payment Webhook
// @Description  Endpoint untuk menerima notifikasi pembayaran dari Midtrans. Tidak membutuhkan autentikasi JWT — divalidasi via signature Midtrans.
// @Tags         Webhook
// @Accept       json
// @Produce      json
// @Param        payload  body  dto.MidtransWebhookPayload  true  "Midtrans Notification Payload"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Router       /webhook/midtrans [post]
func (h *WebhookHandler) HandleMidtrans(c *gin.Context) {
	var payload dto.MidtransWebhookPayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		slog.Warn("midtrans webhook: invalid JSON payload",
			slog.String("error", err.Error()),
			slog.String("remote_addr", c.ClientIP()),
		)
		// Tetap return 200 agar Midtrans tidak retry terus-menerus
		// untuk payload yang memang tidak valid (bukan transient error)
		response.Success(c, http.StatusOK, "Webhook diterima", nil)
		return
	}

	slog.Info("midtrans webhook received",
		slog.String("order_id", payload.OrderID),
		slog.String("transaction_status", payload.TransactionStatus),
		slog.String("fraud_status", payload.FraudStatus),
		slog.String("gross_amount", payload.GrossAmount),
	)

	if err := h.usecase.HandleMidtransWebhook(c.Request.Context(), payload); err != nil {
		var appErr *apperrors.AppError
		if e, ok := err.(*apperrors.AppError); ok {
			appErr = e
		}

		// Untuk signature invalid → return 403 (tidak retry Midtrans untuk ini)
		if appErr != nil && appErr.Code == "INVALID_SIGNATURE" {
			apperrors.Write(c, err)
			return
		}

		// Untuk error internal lainnya → return 500 agar Midtrans retry
		slog.Error("midtrans webhook: processing failed",
			slog.String("order_id", payload.OrderID),
			slog.String("error", err.Error()),
		)
		response.Error(c, http.StatusInternalServerError, "Gagal memproses pembayaran", gin.H{
			"code": "PROCESS_FAILED",
		})
		return
	}

	response.Success(c, http.StatusOK, "Pembayaran berhasil diproses", nil)
}
