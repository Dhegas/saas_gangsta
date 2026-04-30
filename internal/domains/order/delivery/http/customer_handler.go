package http

import (
	"net/http"

	"github.com/dhegas/saas_gangsta/internal/common/response"
	"github.com/dhegas/saas_gangsta/internal/domains/order/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/order/dto"
	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	usecase domain.CustomerUsecase
}

func NewCustomerHandler(usecase domain.CustomerUsecase) *CustomerHandler {
	return &CustomerHandler{usecase: usecase}
}

// CreateCustomer godoc
// @Summary      Tambah Customer ke Order
// @Description  Menambahkan data customer yang terkait dengan sebuah order. Customer hanya dapat dibuat sekali per order.
// @Tags         Customer Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                        true  "Order ID"
// @Param        request  body      dto.CreateCustomerRequest     true  "Payload Create Customer"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      404      {object}  map[string]interface{}
// @Failure      409      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /orders/{id}/customer [post]
func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	tenantID, err := extractTenantIDFromCtx(c)
	if err != nil {
		return
	}
	orderID := c.Param("id")

	var req dto.CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Data tidak valid", gin.H{"code": "VALIDATION_ERROR", "details": err.Error()})
		return
	}

	customer, err := h.usecase.CreateCustomer(c.Request.Context(), tenantID, orderID, req)
	if err != nil {
		return
	}

	response.Success(c, http.StatusCreated, "Berhasil menambahkan data customer", customer)
}

// GetCustomer godoc
// @Summary      Detail Customer dari Order
// @Description  Mengambil data customer yang terkait dengan sebuah order.
// @Tags         Customer Management
// @Produce      json
// @Security     BearerAuth
// @Param        id  path      string  true  "Order ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /orders/{id}/customer [get]
func (h *CustomerHandler) GetCustomer(c *gin.Context) {
	tenantID, err := extractTenantIDFromCtx(c)
	if err != nil {
		return
	}
	orderID := c.Param("id")

	customer, err := h.usecase.GetCustomerByOrderID(c.Request.Context(), tenantID, orderID)
	if err != nil {
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil data customer", customer)
}

// UpdateCustomer godoc
// @Summary      Update Data Customer dari Order
// @Description  Memperbarui data customer yang terkait dengan sebuah order.
// @Tags         Customer Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                        true  "Order ID"
// @Param        request  body      dto.UpdateCustomerRequest     true  "Payload Update Customer"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      404      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /orders/{id}/customer [put]
func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	tenantID, err := extractTenantIDFromCtx(c)
	if err != nil {
		return
	}
	orderID := c.Param("id")

	var req dto.UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Data tidak valid", gin.H{"code": "VALIDATION_ERROR", "details": err.Error()})
		return
	}

	customer, err := h.usecase.UpdateCustomer(c.Request.Context(), tenantID, orderID, req)
	if err != nil {
		return
	}

	response.Success(c, http.StatusOK, "Berhasil memperbarui data customer", customer)
}

// extractTenantIDFromCtx adalah helper lokal untuk mendapatkan tenant_id dari context/query/header
func extractTenantIDFromCtx(c *gin.Context) (string, error) {
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		tenantID = c.GetHeader("X-Tenant-ID")
	}
	if tenantID == "" {
		if val, exists := c.Get("tenantId"); exists {
			if tID, ok := val.(string); ok && tID != "" {
				tenantID = tID
			}
		}
	}

	if tenantID == "" {
		response.Error(c, http.StatusBadRequest, "Tenant ID diperlukan", gin.H{
			"code":    "TENANT_NOT_FOUND",
			"details": "Konteks tenant atau parameter tenant_id tidak ditemukan",
		})
		return "", http.ErrNoCookie
	}
	return tenantID, nil
}
