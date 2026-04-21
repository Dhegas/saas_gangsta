package http

import (
	"net/http"

	"github.com/dhegas/saas_gangsta/internal/common/response"
	"github.com/dhegas/saas_gangsta/internal/domains/category/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/category/dto"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant"
	"github.com/gin-gonic/gin"
)

type MerchantCategoryHandler struct {
	usecase domain.MerchantCategoryUsecase
}

func NewMerchantCategoryHandler(usecase domain.MerchantCategoryUsecase) *MerchantCategoryHandler {
	return &MerchantCategoryHandler{usecase: usecase}
}

// GetAllCategories godoc
// @Summary      List Categories
// @Description  Mengambil daftar category per tenant
// @Tags         Category Management
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /categories [get]
func (h *MerchantCategoryHandler) GetAllCategories(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		return
	}

	categories, err := h.usecase.GetAllCategories(c.Request.Context(), tenantID)
	if err != nil {
		return // Usecase returns formatted apperrors
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil daftar category", categories)
}

// GetCategoryByID godoc
// @Summary      Detail Category
// @Description  Mengambil detail satu category
// @Tags         Category Management
// @Produce      json
// @Security     BearerAuth
// @Param        id  path      string  true  "Category ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /categories/{id} [get]
func (h *MerchantCategoryHandler) GetCategoryByID(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		return
	}
	categoryID := c.Param("id")

	category, err := h.usecase.GetCategoryByID(c.Request.Context(), tenantID, categoryID)
	if err != nil {
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil detail category", category)
}

// CreateCategory godoc
// @Summary      Create Category
// @Description  Buat category baru
// @Tags         Category Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.CreateCategoryRequest  true  "Payload Create"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /categories [post]
func (h *MerchantCategoryHandler) CreateCategory(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		return
	}

	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Data tidak valid", gin.H{"code": "VALIDATION_ERROR", "details": err.Error()})
		return
	}

	category, err := h.usecase.CreateCategory(c.Request.Context(), tenantID, req)
	if err != nil {
		return
	}

	response.Success(c, http.StatusCreated, "Berhasil membuat category", category)
}

// UpdateCategory godoc
// @Summary      Update Category
// @Description  Perbarui nama atau deskripsi category
// @Tags         Category Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                     true  "Category ID"
// @Param        request  body      dto.UpdateCategoryRequest  true  "Payload Update"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      404      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /categories/{id} [put]
func (h *MerchantCategoryHandler) UpdateCategory(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		return
	}
	categoryID := c.Param("id")

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Data tidak valid", gin.H{"code": "VALIDATION_ERROR", "details": err.Error()})
		return
	}

	category, err := h.usecase.UpdateCategory(c.Request.Context(), tenantID, categoryID, req)
	if err != nil {
		return
	}

	response.Success(c, http.StatusOK, "Berhasil memperbarui category", category)
}

// SoftDeleteCategory godoc
// @Summary      Soft Delete Category
// @Description  Hapus category (soft delete)
// @Tags         Category Management
// @Produce      json
// @Security     BearerAuth
// @Param        id  path      string  true  "Category ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /categories/{id} [delete]
func (h *MerchantCategoryHandler) SoftDeleteCategory(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		return
	}
	categoryID := c.Param("id")

	if err := h.usecase.SoftDeleteCategory(c.Request.Context(), tenantID, categoryID); err != nil {
		return
	}

	response.Success(c, http.StatusOK, "Berhasil menghapus category", nil)
}

// ToggleCategoryActive godoc
// @Summary      Aktif / Nonaktifkan Category
// @Description  Ubah status is_active category
// @Tags         Category Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                           true  "Category ID"
// @Param        request  body      dto.ToggleCategoryActiveRequest  true  "Payload Toggle"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      404      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /categories/{id}/toggle-active [patch]
func (h *MerchantCategoryHandler) ToggleCategoryActive(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		return
	}
	categoryID := c.Param("id")

	var req dto.ToggleCategoryActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Data tidak valid", gin.H{"code": "VALIDATION_ERROR", "details": err.Error()})
		return
	}

	if err := h.usecase.ToggleCategoryActive(c.Request.Context(), tenantID, categoryID, *req.IsActive); err != nil {
		return
	}

	response.Success(c, http.StatusOK, "Berhasil memperbarui status category", nil)
}

// ReorderCategories godoc
// @Summary      Update sort_order
// @Description  Ubah urutan (sort_order) untuk beberapa category sekaligus
// @Tags         Category Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.ReorderCategoryRequest  true  "Payload Reorder"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /categories/reorder [patch]
func (h *MerchantCategoryHandler) ReorderCategories(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		return
	}

	var req dto.ReorderCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Data tidak valid", gin.H{"code": "VALIDATION_ERROR", "details": err.Error()})
		return
	}

	if err := h.usecase.ReorderCategories(c.Request.Context(), tenantID, req); err != nil {
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengurutkan category", nil)
}
