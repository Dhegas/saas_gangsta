package http

import (
	"net/http"
	"strings"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/common/response"
	"github.com/dhegas/saas_gangsta/internal/domains/order/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/order/dto"
	"github.com/gin-gonic/gin"
)

type CustomerOrderHandler struct {
	usecase domain.PartnerOrderUsecase
}

// NewCustomerOrderHandler konstruktor untuk CustomerOrderHandler
func NewCustomerOrderHandler(usecase domain.PartnerOrderUsecase) *CustomerOrderHandler {
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

	// ADAPTER: Map CreateCustomerOrderRequest (guest mode) to unified CreateOrderRequest
	items := make([]dto.CreateOrderItemRequest, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, dto.CreateOrderItemRequest{
			MenuID:    item.MenuID,
			Quantity:  item.Quantity,
			Notes:     item.Notes,
		})
	}

	adaptedReq := dto.CreateOrderRequest{
		DiningTablesID: req.DiningTableID,
		Items:          items,
		Customer: &dto.CreateCustomerDetailsRequest{
			FullName:    req.Customer.FullName,
			PhoneNumber: req.Customer.PhoneNumber,
		},
	}

	// Panggil usecase terpadu (Unified CreateOrder)
	res, err := h.usecase.CreateOrder(c.Request.Context(), tenantID, adaptedReq)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	// ADAPTER: Map unified OrderResponse ke CreateCustomerOrderResponse (camelCase)
	response.Success(c, http.StatusCreated, "Order created successfully", dto.CreateCustomerOrderResponse{
		OrderID:    res.ID,
		Status:     strings.ToLower(res.Status),
		TotalPrice: res.TotalPrice,
	})
}

