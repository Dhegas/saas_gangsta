package http

import (
	"net/http"

	"github.com/dhegas/saas_gangsta/internal/common/response"
	"github.com/dhegas/saas_gangsta/internal/domains/menu/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/menu/dto"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant"
	"github.com/gin-gonic/gin"
)

type PartnerMenuHandler struct {
	usecase domain.PartnerMenuUsecase
}

func NewPartnerMenuHandler(usecase domain.PartnerMenuUsecase) *PartnerMenuHandler {
	return &PartnerMenuHandler{usecase: usecase}
}

// GetAllMenus godoc
// @Summary      List Menus
// @Description  Mengambil daftar menu per tenant (mendukung filter category_id dan is_available)
// @Tags         Menu Management
// @Produce      json
// @Security     BearerAuth
// @Param        category_id   query     string  false  "Filter berdasarkan Category ID"
// @Param        is_available  query     bool    false  "Filter berdasarkan ketersediaan (true/false)"
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /menus [get]
func (h *PartnerMenuHandler) GetAllMenus(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		tenantID = c.Query("tenant_id")
		if tenantID == "" {
			tenantID = c.GetHeader("X-Tenant-ID")
		}
		if tenantID == "" {
			response.Error(c, http.StatusBadRequest, "Tenant ID diperlukan", gin.H{"code": "TENANT_NOT_FOUND", "details": "Konteks tenant atau parameter tenant_id tidak ditemukan"})
			return
		}
	}

	var filter dto.MenuFilterParams
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.Error(c, http.StatusBadRequest, "Parameter query tidak valid", gin.H{"code": "VALIDATION_ERROR", "details": err.Error()})
		return
	}

	menus, err := h.usecase.GetAllMenus(c.Request.Context(), tenantID, filter)
	if err != nil {
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil daftar menu", menus)
}

// GetMenuByID godoc
// @Summary      Detail Menu
// @Description  Mengambil detail satu menu
// @Tags         Menu Management
// @Produce      json
// @Security     BearerAuth
// @Param        id  path      string  true  "Menu ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /menus/{id} [get]
func (h *PartnerMenuHandler) GetMenuByID(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		tenantID = c.Query("tenant_id")
		if tenantID == "" {
			tenantID = c.GetHeader("X-Tenant-ID")
		}
		if tenantID == "" {
			response.Error(c, http.StatusBadRequest, "Tenant ID diperlukan", gin.H{"code": "TENANT_NOT_FOUND", "details": "Konteks tenant atau parameter tenant_id tidak ditemukan"})
			return
		}
	}
	menuID := c.Param("id")

	menu, err := h.usecase.GetMenuByID(c.Request.Context(), tenantID, menuID)
	if err != nil {
		return
	}

	response.Success(c, http.StatusOK, "Berhasil mengambil detail menu", menu)
}

// CreateMenu godoc
// @Summary      Create Menu
// @Description  Mendaftarkan menu baru
// @Tags         Menu Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.CreateMenuRequest  true  "Payload Create Menu"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /menus [post]
func (h *PartnerMenuHandler) CreateMenu(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		return
	}

	var req dto.CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Data tidak valid", gin.H{"code": "VALIDATION_ERROR", "details": err.Error()})
		return
	}

	menu, err := h.usecase.CreateMenu(c.Request.Context(), tenantID, req)
	if err != nil {
		return
	}

	response.Success(c, http.StatusCreated, "Berhasil membuat menu", menu)
}

// UpdateMenu godoc
// @Summary      Update Menu
// @Description  Memperbarui informasi menu
// @Tags         Menu Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                 true  "Menu ID"
// @Param        request  body      dto.UpdateMenuRequest  true  "Payload Update Menu"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      404      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /menus/{id} [put]
func (h *PartnerMenuHandler) UpdateMenu(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		return
	}
	menuID := c.Param("id")

	var req dto.UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Data tidak valid", gin.H{"code": "VALIDATION_ERROR", "details": err.Error()})
		return
	}

	menu, err := h.usecase.UpdateMenu(c.Request.Context(), tenantID, menuID, req)
	if err != nil {
		return
	}

	response.Success(c, http.StatusOK, "Berhasil memperbarui menu", menu)
}

// SoftDeleteMenu godoc
// @Summary      Soft Delete Menu
// @Description  Menghapus menu (soft delete)
// @Tags         Menu Management
// @Produce      json
// @Security     BearerAuth
// @Param        id  path      string  true  "Menu ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /menus/{id} [delete]
func (h *PartnerMenuHandler) SoftDeleteMenu(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		return
	}
	menuID := c.Param("id")

	if err := h.usecase.SoftDeleteMenu(c.Request.Context(), tenantID, menuID); err != nil {
		return
	}

	response.Success(c, http.StatusOK, "Berhasil menghapus menu", nil)
}

// ToggleMenuAvailable godoc
// @Summary      Toggle menu availability
// @Description  Mengubah status is_available pada menu
// @Tags         Menu Management
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                          true  "Menu ID"
// @Param        request  body      dto.ToggleMenuAvailableRequest  true  "Payload Toggle"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      404      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /menus/{id}/toggle-available [patch]
func (h *PartnerMenuHandler) ToggleMenuAvailable(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		return
	}
	menuID := c.Param("id")

	var req dto.ToggleMenuAvailableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Data tidak valid", gin.H{"code": "VALIDATION_ERROR", "details": err.Error()})
		return
	}

	if err := h.usecase.ToggleMenuAvailable(c.Request.Context(), tenantID, menuID, *req.IsAvailable); err != nil {
		return
	}

	response.Success(c, http.StatusOK, "Berhasil memperbarui status ketersediaan menu", nil)
}
