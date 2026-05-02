package http

import (
	"net/http"

	"github.com/dhegas/saas_gangsta/internal/common/response"
	"github.com/dhegas/saas_gangsta/internal/domains/order/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/order/dto"
	"github.com/gin-gonic/gin"
)

type PartnerOrderHandler struct {
	usecase domain.PartnerOrderUsecase
}

func NewPartnerOrderHandler(usecase domain.PartnerOrderUsecase) *PartnerOrderHandler {
	return &PartnerOrderHandler{usecase: usecase}
}

// extractTenantID is a helper to get tenant_id for both CUSTOMER and PARTNER
func (h *PartnerOrderHandler) extractTenantID(c *gin.Context) (string, error) {
	// First check query param (mostly for CUSTOMER role)
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		// Then check header
		tenantID = c.GetHeader("X-Tenant-ID")
	}
	if tenantID == "" {
		// Finally check gin context (for PARTNER who passed TenantGuard)
		if val, exists := c.Get("tenantId"); exists {
			if tID, ok := val.(string); ok && tID != "" {
				tenantID = tID
			}
		}
	}

	if tenantID == "" {
		response.Error(c, http.StatusBadRequest, "Tenant ID diperlukan", gin.H{"code": "TENANT_NOT_FOUND", "details": "Konteks tenant atau parameter tenant_id tidak ditemukan"})
		return "", http.ErrNoCookie
	}
	return tenantID, nil
}

// GetAllOrders godoc
// @Summary      List Orders
// @Description  Mengambil daftar pesanan per tenant (Partner only)
// @Tags         Order Management
// @Produce      json
// @Security     BearerAuth
// @Param        tenant_id query     string  false  "Tenant ID (Wajib untuk CUSTOMER)"
// @Param        status    query     string  false  "Filter by status (PENDING, PROCESSING, dll)"
// @Param        table_id  query     string  false  "Filter by table id"
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /orders [get]
func (h *PartnerOrderHandler) GetAllOrders(c *gin.Context) {
	tenantID, err := h.extractTenantID(c)
	if err != nil {
		return
	}

	var filter dto.OrderFilterParams
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.Error(c, http.StatusBadRequest, "Parameter query tidak valid", gin.H{"code": "VALIDATION_ERROR", "details": err.Error()})
		return
	}

	orders, err := h.usecase.GetAllOrders(c.Request.Context(), tenantID, filter)
	if err != nil {
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil daftar pesanan", orders)
}

// GetOrderByID godoc
// @Summary      Detail Order
// @Description  Mengambil detail pesanan beserta itemnya
// @Tags         Order Management
// @Produce      json
// @Security     BearerAuth
// @Param        id        path      string  true   "Order ID"
// @Param        tenant_id query     string  false  "Tenant ID (Wajib untuk CUSTOMER)"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /orders/{id} [get]
func (h *PartnerOrderHandler) GetOrderByID(c *gin.Context) {
	tenantID, err := h.extractTenantID(c)
	if err != nil {
		return
	}
	orderID := c.Param("id")

	order, err := h.usecase.GetOrderByID(c.Request.Context(), tenantID, orderID)
	if err != nil {
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil detail pesanan", order)
}

// CreateOrder godoc
// @Summary      Create Order
// @Description  Membuat pesanan baru (Customer & Partner)
// @Tags         Order Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        tenant_id query     string                  false "Tenant ID (Wajib untuk Customer, opsional untuk Partner)"
// @Param        request   body      dto.CreateOrderRequest  true  "Payload Create Order"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /orders [post]
func (h *PartnerOrderHandler) CreateOrder(c *gin.Context) {
	tenantID, err := h.extractTenantID(c)
	if err != nil {
		return
	}

	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Data tidak valid", gin.H{"code": "VALIDATION_ERROR", "details": err.Error()})
		return
	}

	order, err := h.usecase.CreateOrder(c.Request.Context(), tenantID, req)
	if err != nil {
		return
	}

	response.Success(c, http.StatusCreated, "Berhasil membuat pesanan", order)
}

// UpdateOrderStatus godoc
// @Summary      Update Order Status
// @Description  Memperbarui status pesanan (Partner Only)
// @Tags         Order Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                        true  "Order ID"
// @Param        request  body      dto.UpdateOrderStatusRequest  true  "Payload Update Status"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      404      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /orders/{id}/status [patch]
func (h *PartnerOrderHandler) UpdateOrderStatus(c *gin.Context) {
	tenantID, err := h.extractTenantID(c)
	if err != nil {
		return
	}
	orderID := c.Param("id")

	var req dto.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Data tidak valid", gin.H{"code": "VALIDATION_ERROR", "details": err.Error()})
		return
	}

	order, err := h.usecase.UpdateOrderStatus(c.Request.Context(), tenantID, orderID, req)
	if err != nil {
		return
	}

	response.Success(c, http.StatusOK, "Berhasil memperbarui status pesanan", order)
}

// SoftDeleteOrder godoc
// @Summary      Soft Delete / Cancel Order
// @Description  Membatalkan/menghapus pesanan secara logic (Partner Only)
// @Tags         Order Management
// @Produce      json
// @Security     BearerAuth
// @Param        id  path      string  true  "Order ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /orders/{id} [delete]
func (h *PartnerOrderHandler) SoftDeleteOrder(c *gin.Context) {
	tenantID, err := h.extractTenantID(c)
	if err != nil {
		return
	}
	orderID := c.Param("id")

	if err := h.usecase.SoftDeleteOrder(c.Request.Context(), tenantID, orderID); err != nil {
		return
	}

	response.Success(c, http.StatusOK, "Berhasil menghapus pesanan", nil)
}
