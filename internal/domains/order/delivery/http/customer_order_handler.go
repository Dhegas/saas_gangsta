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
		DiningTablesID: &req.DiningTableID,
		Items:          items,
		Customer: &dto.CreateCustomerDetailsRequest{
			FullName:    req.Customer.FullName,
			Email:       req.Customer.Email,
			Password:    req.Customer.Password,
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
		OrderID:     res.ID,
		Status:      strings.ToLower(res.Status),
		TotalPrice:  res.TotalPrice,
		AccessToken: res.AccessToken,
	})
}

// GetOrderStatus godoc
// @Summary      Get Public Order Status
// @Description  Melihat detail dan status pesanan secara publik (untuk tracking status customer)
// @Tags         Public Customer Order
// @Produce      json
// @Param        tenantSlug path      string  true  "Tenant Slug"
// @Param        orderId    path      string  true  "Order ID"
// @Success      200        {object}  response.Envelope{data=dto.PublicOrderDetailsResponse}
// @Failure      400        {object}  response.Envelope
// @Failure      404        {object}  response.Envelope
// @Failure      500        {object}  response.Envelope
// @Router       /public/tenant/{tenantSlug}/orders/{orderId} [get]
func (h *CustomerOrderHandler) GetOrderStatus(c *gin.Context) {
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

	orderID := c.Param("orderId")
	if orderID == "" {
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Order ID wajib disertakan", http.StatusBadRequest, nil))
		return
	}

	res, err := h.usecase.GetPublicOrderStatus(c.Request.Context(), tenantID, orderID)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Order detail fetched successfully", res)
}

// GetPublicOrders godoc
// @Summary      Get Public Orders List
// @Description  Melihat daftar pesanan secara publik (misal untuk antrean / monitoring status customer)
// @Tags         Public Customer Order
// @Produce      json
// @Param        tenantSlug path      string  true  "Tenant Slug"
// @Param        status     query     string  false "Filter by status"
// @Param        table_id   query     string  false "Filter by Table ID"
// @Success      200        {object}  response.Envelope{data=[]dto.PublicOrderDetailsResponse}
// @Failure      400        {object}  response.Envelope
// @Failure      500        {object}  response.Envelope
// @Router       /public/tenant/{tenantSlug}/orders [get]
func (h *CustomerOrderHandler) GetPublicOrders(c *gin.Context) {
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

	var filter dto.PublicOrderFilterParams
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.Error(c, http.StatusBadRequest, "Parameter query tidak valid", gin.H{"code": "VALIDATION_ERROR", "details": err.Error()})
		return
	}

	res, err := h.usecase.GetPublicOrdersList(c.Request.Context(), tenantID, filter)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Orders list fetched successfully", res)
}
