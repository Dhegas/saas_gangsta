package http

import (
	"net/http"

	"github.com/dhegas/saas_gangsta/internal/common/response"
	"github.com/dhegas/saas_gangsta/internal/domains/table/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/table/dto"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant"
	"github.com/gin-gonic/gin"
)

type PartnerTableHandler struct {
	usecase domain.PartnerTableUsecase
}

func NewPartnerTableHandler(usecase domain.PartnerTableUsecase) *PartnerTableHandler {
	return &PartnerTableHandler{usecase: usecase}
}

// GetAllTables godoc
// @Summary      List Tables
// @Description  Mengambil daftar meja per tenant
// @Tags         Dining Table Management
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /dining-tables [get]
func (h *PartnerTableHandler) GetAllTables(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		return
	}

	tables, err := h.usecase.GetAllTables(c.Request.Context(), tenantID)
	if err != nil {
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil daftar meja", tables)
}

// GetTableByID godoc
// @Summary      Detail Table
// @Description  Mengambil detail satu meja
// @Tags         Dining Table Management
// @Produce      json
// @Security     BearerAuth
// @Param        id  path      string  true  "Table ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /dining-tables/{id} [get]
func (h *PartnerTableHandler) GetTableByID(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		return
	}
	tableID := c.Param("id")

	table, err := h.usecase.GetTableByID(c.Request.Context(), tenantID, tableID)
	if err != nil {
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil detail meja", table)
}

// GetTableStatus godoc
// @Summary      Cek Status Meja
// @Description  Mengecek apakah meja sedang kosong atau occupied (berdasarkan transaksi pesanan)
// @Tags         Dining Table Management
// @Produce      json
// @Security     BearerAuth
// @Param        id  path      string  true  "Table ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /dining-tables/{id}/status [get]
func (h *PartnerTableHandler) GetTableStatus(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		return
	}
	tableID := c.Param("id")

	statusRes, err := h.usecase.GetTableStatus(c.Request.Context(), tenantID, tableID)
	if err != nil {
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengecek status meja", statusRes)
}

// CreateTable godoc
// @Summary      Create Table
// @Description  Mendaftarkan meja baru
// @Tags         Dining Table Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.CreateTableRequest  true  "Payload Create Table"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /dining-tables [post]
func (h *PartnerTableHandler) CreateTable(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		return
	}

	var req dto.CreateTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Data tidak valid", gin.H{"code": "VALIDATION_ERROR", "details": err.Error()})
		return
	}

	table, err := h.usecase.CreateTable(c.Request.Context(), tenantID, req)
	if err != nil {
		return
	}

	response.Success(c, http.StatusCreated, "Berhasil membuat meja", table)
}

// UpdateTable godoc
// @Summary      Update Table
// @Description  Memperbarui nama meja
// @Tags         Dining Table Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                  true  "Table ID"
// @Param        request  body      dto.UpdateTableRequest  true  "Payload Update Table"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      404      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /dining-tables/{id} [put]
func (h *PartnerTableHandler) UpdateTable(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		return
	}
	tableID := c.Param("id")

	var req dto.UpdateTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Data tidak valid", gin.H{"code": "VALIDATION_ERROR", "details": err.Error()})
		return
	}

	table, err := h.usecase.UpdateTable(c.Request.Context(), tenantID, tableID, req)
	if err != nil {
		return
	}

	response.Success(c, http.StatusOK, "Berhasil memperbarui meja", table)
}

// SoftDeleteTable godoc
// @Summary      Soft Delete Table
// @Description  Menghapus meja (soft delete)
// @Tags         Dining Table Management
// @Produce      json
// @Security     BearerAuth
// @Param        id  path      string  true  "Table ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /dining-tables/{id} [delete]
func (h *PartnerTableHandler) SoftDeleteTable(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		return
	}
	tableID := c.Param("id")

	if err := h.usecase.SoftDeleteTable(c.Request.Context(), tenantID, tableID); err != nil {
		return
	}

	response.Success(c, http.StatusOK, "Berhasil menghapus meja", nil)
}
