package http

import (
	"net/http"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/common/response"
	"github.com/dhegas/saas_gangsta/internal/domains/order/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/order/dto"
	"github.com/gin-gonic/gin"
)

type CustomerOrderHandler struct {
	usecase domain.CustomerOrderUsecase
}

// NewCustomerOrderHandler konstruktor untuk CustomerOrderHandler
func NewCustomerOrderHandler(usecase domain.CustomerOrderUsecase) *CustomerOrderHandler {
	return &CustomerOrderHandler{usecase: usecase}
}

// CreateOrder godoc
// @Summary      Create Public Customer Order
// @Description  Membuat pesanan baru oleh pelanggan lewat halaman menu publik menggunakan scan QR code (self-order)
// @Tags         Public Customer Order
// @Accept       json
// @Produce      json
// @Param        tenantSlug path      string                         true  "Tenant Slug"
// @Param        request    body      dto.CreateCustomerOrderRequest true  "Payload Create Customer Order"
// @Success      201        {object}  response.Envelope{data=dto.CreateCustomerOrderResponse}
// @Failure      400        {object}  response.Envelope
// @Failure      500        {object}  response.Envelope
// @Router       /public/tenant/{tenantSlug}/orders [post]
func (h *CustomerOrderHandler) CreateOrder(c *gin.Context) {
	tenantIDVal, exists := c.Get("tenantId")
	if !exists {
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Tenant context tidak ditemukan", http.StatusBadRequest, nil))
		return
	}
	tenantID, ok := tenantIDVal.(string)
	if !ok || tenantID == "" {
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Tenant context tidak valid", http.StatusBadRequest, nil))
		return
	}

	var req dto.CreateCustomerOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Data payload order tidak valid", gin.H{
			"code":    "VALIDATION_ERROR",
			"details": err.Error(),
		})
		return
	}

	res, err := h.usecase.CreateCustomerOrder(c.Request.Context(), tenantID, req)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Order created successfully", res)
}
