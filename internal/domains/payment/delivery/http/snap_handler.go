package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/common/response"
	paymentdomain "github.com/dhegas/saas_gangsta/internal/domains/payment/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/payment/dto"
)

// SnapHandler menangani inisiasi pembayaran via Midtrans Snap
type SnapHandler struct {
	usecase paymentdomain.PaymentSnapUsecase
}

func NewSnapHandler(usecase paymentdomain.PaymentSnapUsecase) *SnapHandler {
	return &SnapHandler{usecase: usecase}
}

// CreateSnapTransaction godoc
// @Summary      Inisiasi Pembayaran Midtrans Snap
// @Description  Membuat Snap token Midtrans untuk order yang belum dibayar. Customer menggunakan snap_token atau redirect_url untuk menyelesaikan pembayaran di halaman Midtrans.
// @Tags         Payment
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body  dto.CreateSnapRequest  true  "Order ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      422  {object}  map[string]interface{}
// @Router       /customer/payments/initiate [post]
func (h *SnapHandler) CreateSnapTransaction(c *gin.Context) {
	var req dto.CreateSnapRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "order_id diperlukan", http.StatusBadRequest))
		return
	}

	userID := c.GetString("userId")
	if userID == "" {
		apperrors.Write(c, apperrors.New("UNAUTHORIZED", "User tidak teridentifikasi", http.StatusUnauthorized))
		return
	}

	result, err := h.usecase.CreateSnapTransaction(c.Request.Context(), userID, req.OrderID)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Snap token berhasil dibuat", result)
}
