package http

import (
	"net/http"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/common/response"
	"github.com/dhegas/saas_gangsta/internal/domains/menu/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/menu/dto"
	"github.com/gin-gonic/gin"
)

// AdminMenuHandler mengelola endpoint admin untuk manajemen menu tenant partner.
// Tenant ID diambil langsung dari header X-Tenant-ID karena JWT admin tidak
// mengandung tenantId — validasi keamanan dijamin oleh RoleGuard("ADMIN").
type AdminMenuHandler struct {
	usecase domain.PartnerMenuUsecase
}

func NewAdminMenuHandler(usecase domain.PartnerMenuUsecase) *AdminMenuHandler {
	return &AdminMenuHandler{usecase: usecase}
}

// extractAdminTenantID mengambil tenant ID dari header X-Tenant-ID atau
// query param tenant_id. Return error jika tidak ada.
func (h *AdminMenuHandler) extractAdminTenantID(c *gin.Context) (string, error) {
	tenantID := c.GetHeader("X-Tenant-ID")
	if tenantID == "" {
		tenantID = c.Query("tenant_id")
	}
	if tenantID == "" {
		return "", apperrors.New("TENANT_NOT_FOUND", "Header X-Tenant-ID wajib disertakan untuk endpoint admin menu", http.StatusBadRequest, nil)
	}
	return tenantID, nil
}

// GetAllMenus godoc
// @Summary      [Admin] List Menus by Tenant
// @Description  Admin mengambil daftar menu milik tenant partner tertentu
// @Tags         Admin
// @Produce      json
// @Security     BearerAuth
// @Param        X-Tenant-ID   header    string  true   "Tenant ID yang ingin dikelola"
// @Param        category_id   query     string  false  "Filter berdasarkan Category ID"
// @Param        is_available  query     bool    false  "Filter berdasarkan ketersediaan"
// @Success      200  {object}  response.Envelope
// @Failure      400  {object}  response.Envelope
// @Failure      401  {object}  response.Envelope
// @Failure      403  {object}  response.Envelope
// @Failure      500  {object}  response.Envelope
// @Router       /admin/menus [get]
func (h *AdminMenuHandler) GetAllMenus(c *gin.Context) {
	tenantID, err := h.extractAdminTenantID(c)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	var filter dto.MenuFilterParams
	if err := c.ShouldBindQuery(&filter); err != nil {
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Parameter query tidak valid", http.StatusBadRequest, err.Error()))
		return
	}

	menus, err := h.usecase.GetAllMenus(c.Request.Context(), tenantID, filter)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Daftar menu berhasil diambil oleh admin", menus)
}

// GetMenuByID godoc
// @Summary      [Admin] Detail Menu
// @Description  Admin mengambil detail satu menu dari tenant partner tertentu
// @Tags         Admin
// @Produce      json
// @Security     BearerAuth
// @Param        X-Tenant-ID  header  string  true  "Tenant ID yang ingin dikelola"
// @Param        id           path    string  true  "Menu ID (UUID)"
// @Success      200  {object}  response.Envelope
// @Failure      400  {object}  response.Envelope
// @Failure      401  {object}  response.Envelope
// @Failure      403  {object}  response.Envelope
// @Failure      404  {object}  response.Envelope
// @Failure      500  {object}  response.Envelope
// @Router       /admin/menus/{id} [get]
func (h *AdminMenuHandler) GetMenuByID(c *gin.Context) {
	tenantID, err := h.extractAdminTenantID(c)
	if err != nil {
		apperrors.Write(c, err)
		return
	}
	menuID := c.Param("id")

	menu, err := h.usecase.GetMenuByID(c.Request.Context(), tenantID, menuID)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Detail menu berhasil diambil oleh admin", menu)
}

// CreateMenu godoc
// @Summary      [Admin] Create Menu
// @Description  Admin membuat menu baru pada tenant partner tertentu
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        X-Tenant-ID  header  string                 true  "Tenant ID yang ingin dikelola"
// @Param        request      body    dto.CreateMenuRequest  true  "Payload Create Menu"
// @Success      201  {object}  response.Envelope
// @Failure      400  {object}  response.Envelope
// @Failure      401  {object}  response.Envelope
// @Failure      403  {object}  response.Envelope
// @Failure      500  {object}  response.Envelope
// @Router       /admin/menus [post]
func (h *AdminMenuHandler) CreateMenu(c *gin.Context) {
	tenantID, err := h.extractAdminTenantID(c)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	var req dto.CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Data menu tidak valid", http.StatusBadRequest, err.Error()))
		return
	}

	menu, err := h.usecase.CreateMenu(c.Request.Context(), tenantID, req)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Menu berhasil dibuat oleh admin", menu)
}

// UpdateMenu godoc
// @Summary      [Admin] Update Menu
// @Description  Admin memperbarui informasi menu pada tenant partner tertentu
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        X-Tenant-ID  header  string                 true  "Tenant ID yang ingin dikelola"
// @Param        id           path    string                 true  "Menu ID (UUID)"
// @Param        request      body    dto.UpdateMenuRequest  true  "Payload Update Menu"
// @Success      200  {object}  response.Envelope
// @Failure      400  {object}  response.Envelope
// @Failure      401  {object}  response.Envelope
// @Failure      403  {object}  response.Envelope
// @Failure      404  {object}  response.Envelope
// @Failure      500  {object}  response.Envelope
// @Router       /admin/menus/{id} [put]
func (h *AdminMenuHandler) UpdateMenu(c *gin.Context) {
	tenantID, err := h.extractAdminTenantID(c)
	if err != nil {
		apperrors.Write(c, err)
		return
	}
	menuID := c.Param("id")

	var req dto.UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Data update menu tidak valid", http.StatusBadRequest, err.Error()))
		return
	}

	menu, err := h.usecase.UpdateMenu(c.Request.Context(), tenantID, menuID, req)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Menu berhasil diperbarui oleh admin", menu)
}

// SoftDeleteMenu godoc
// @Summary      [Admin] Delete Menu
// @Description  Admin menghapus (soft delete) menu pada tenant partner tertentu
// @Tags         Admin
// @Produce      json
// @Security     BearerAuth
// @Param        X-Tenant-ID  header  string  true  "Tenant ID yang ingin dikelola"
// @Param        id           path    string  true  "Menu ID (UUID)"
// @Success      200  {object}  response.Envelope
// @Failure      400  {object}  response.Envelope
// @Failure      401  {object}  response.Envelope
// @Failure      403  {object}  response.Envelope
// @Failure      404  {object}  response.Envelope
// @Failure      500  {object}  response.Envelope
// @Router       /admin/menus/{id} [delete]
func (h *AdminMenuHandler) SoftDeleteMenu(c *gin.Context) {
	tenantID, err := h.extractAdminTenantID(c)
	if err != nil {
		apperrors.Write(c, err)
		return
	}
	menuID := c.Param("id")

	if err := h.usecase.SoftDeleteMenu(c.Request.Context(), tenantID, menuID); err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Menu berhasil dihapus oleh admin", nil)
}

// ToggleMenuAvailable godoc
// @Summary      [Admin] Toggle Menu Availability
// @Description  Admin mengubah status ketersediaan menu pada tenant partner tertentu
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        X-Tenant-ID  header  string                          true  "Tenant ID yang ingin dikelola"
// @Param        id           path    string                          true  "Menu ID (UUID)"
// @Param        request      body    dto.ToggleMenuAvailableRequest  true  "Payload Toggle"
// @Success      200  {object}  response.Envelope
// @Failure      400  {object}  response.Envelope
// @Failure      401  {object}  response.Envelope
// @Failure      403  {object}  response.Envelope
// @Failure      404  {object}  response.Envelope
// @Failure      500  {object}  response.Envelope
// @Router       /admin/menus/{id}/toggle-available [patch]
func (h *AdminMenuHandler) ToggleMenuAvailable(c *gin.Context) {
	tenantID, err := h.extractAdminTenantID(c)
	if err != nil {
		apperrors.Write(c, err)
		return
	}
	menuID := c.Param("id")

	var req dto.ToggleMenuAvailableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Data toggle menu tidak valid", http.StatusBadRequest, err.Error()))
		return
	}

	if err := h.usecase.ToggleMenuAvailable(c.Request.Context(), tenantID, menuID, *req.IsAvailable); err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Status ketersediaan menu berhasil diubah oleh admin", nil)
}
