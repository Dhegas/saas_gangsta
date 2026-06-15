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
	usecase domain.PartnerOrderUsecase
}

// NewCustomerOrderHandler konstruktor untuk CustomerOrderHandler
func NewCustomerOrderHandler(usecase domain.PartnerOrderUsecase) *CustomerOrderHandler {
	return &CustomerOrderHandler{usecase: usecase}
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
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Tenant context tidak ditemukan", http.StatusBadRequest))
		return
	}
	tenantID, ok := tenantIDVal.(string)
	if !ok || tenantID == "" {
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Tenant context tidak valid", http.StatusBadRequest))
		return
	}

	orderID := c.Param("orderId")
	if orderID == "" || len(orderID) != 36 {
		apperrors.Write(c, apperrors.New("NOT_FOUND", "Pesanan tidak ditemukan atau ID pesanan tidak valid", http.StatusNotFound))
		return
	}

	userIDVal, exists := c.Get("userId")
	if !exists {
		apperrors.Write(c, apperrors.New("UNAUTHORIZED", "Authentication required", http.StatusUnauthorized))
		return
	}
	uID, ok := userIDVal.(string)
	if !ok || uID == "" {
		apperrors.Write(c, apperrors.New("UNAUTHORIZED", "Authentication required", http.StatusUnauthorized))
		return
	}

	res, err := h.usecase.GetPublicOrderStatus(c.Request.Context(), tenantID, orderID)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	// Proteksi anti-IDOR: pastikan order adalah milik user yang sedang login
	if res.UserID == nil || *res.UserID != uID {
		apperrors.Write(c, apperrors.New("FORBIDDEN", "You do not have permission to perform this action", http.StatusForbidden))
		return
	}

	response.Success(c, http.StatusOK, "Order detail fetched successfully", res)
}

// GetPublicOrders godoc
// @Summary      Get Public Orders List
// @Description  Melihat daftar pesanan secara publik (hanya untuk pesanan milik customer yang login)
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
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Tenant context tidak ditemukan", http.StatusBadRequest))
		return
	}
	tenantID, ok := tenantIDVal.(string)
	if !ok || tenantID == "" {
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Tenant context tidak valid", http.StatusBadRequest))
		return
	}

	userIDVal, exists := c.Get("userId")
	if !exists {
		apperrors.Write(c, apperrors.New("UNAUTHORIZED", "Authentication required", http.StatusUnauthorized))
		return
	}
	uID, ok := userIDVal.(string)
	if !ok || uID == "" {
		apperrors.Write(c, apperrors.New("UNAUTHORIZED", "Authentication required", http.StatusUnauthorized))
		return
	}

	var filter dto.PublicOrderFilterParams
	if err := c.ShouldBindQuery(&filter); err != nil {
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Parameter query tidak valid", http.StatusUnprocessableEntity))
		return
	}

	filter.UserID = uID

	res, err := h.usecase.GetPublicOrdersList(c.Request.Context(), tenantID, filter)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Orders list fetched successfully", res)
}

// GetCustomerOrderHistory godoc
// @Summary      Get Customer Order History
// @Description  Melihat seluruh riwayat pesanan customer dari berbagai toko
// @Tags         Public Customer Order
// @Produce      json
// @Success      200        {object}  response.Envelope{data=[]dto.PublicOrderDetailsResponse}
// @Failure      401        {object}  response.Envelope
// @Failure      500        {object}  response.Envelope
// @Router       /customer/orders/history [get]
func (h *CustomerOrderHandler) GetCustomerOrderHistory(c *gin.Context) {
	userIDVal, exists := c.Get("userId")
	if !exists {
		apperrors.Write(c, apperrors.New("UNAUTHORIZED", "Authentication required", http.StatusUnauthorized))
		return
	}
	uID, ok := userIDVal.(string)
	if !ok || uID == "" {
		apperrors.Write(c, apperrors.New("UNAUTHORIZED", "Authentication required", http.StatusUnauthorized))
		return
	}

	res, err := h.usecase.GetCustomerOrderHistory(c.Request.Context(), uID)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Customer order history fetched successfully", res)
}

