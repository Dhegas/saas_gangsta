package http

import (
	"net/http"
	"strings"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/common/response"
	categorydto "github.com/dhegas/saas_gangsta/internal/domains/category/dto"
	menudto "github.com/dhegas/saas_gangsta/internal/domains/menu/dto"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PublicResourcesHandler struct {
	db *gorm.DB
}

func NewPublicResourcesHandler(db *gorm.DB) *PublicResourcesHandler {
	return &PublicResourcesHandler{db: db}
}

// GetPublicCategories godoc
// @Summary List active menu categories for a tenant
// @Description Mengambil daftar kategori menu aktif dalam context tenant
// @Tags Public
// @Produce json
// @Param tenantSlug path string true "Tenant Slug"
// @Success 200 {object} response.Envelope{data=[]dto.CategoryResponse}
// @Failure 400 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /public/t/{tenantSlug}/categories [get]
func (h *PublicResourcesHandler) GetPublicCategories(c *gin.Context) {
	tenantID := c.GetString("tenantId")
	if tenantID == "" {
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Tenant context tidak ditemukan", http.StatusBadRequest, nil))
		return
	}

	var categories []categorydto.CategoryResponse
	err := h.db.WithContext(c.Request.Context()).Table("categories").
		Where("tenant_id = ?", tenantID).
		Where("is_active = true").
		Where("deleted_at IS NULL").
		Order("sort_order ASC, name ASC").
		Select("id::text AS id, tenant_id::text AS tenant_id, name, description, sort_order, is_active, created_at, updated_at").
		Scan(&categories).Error

	if err != nil {
		apperrors.Write(c, apperrors.New("INTERNAL_ERROR", "Gagal mengambil daftar kategori", http.StatusInternalServerError, err.Error()))
		return
	}

	response.Success(c, http.StatusOK, "Category list fetched successfully", categories)
}

// GetPublicMenus godoc
// @Summary List active menus for a tenant
// @Description Mengambil daftar menu aktif dalam context tenant dengan filter pencarian dan kategori
// @Tags Public
// @Produce json
// @Param tenantSlug path string true "Tenant Slug"
// @Param categoryId query string false "Filter berdasarkan Category ID"
// @Param search query string false "Pencarian nama atau deskripsi menu"
// @Param isAvailable query bool false "Filter ketersediaan menu"
// @Success 200 {object} response.Envelope{data=[]dto.MenuResponse}
// @Failure 400 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /public/t/{tenantSlug}/menus [get]
func (h *PublicResourcesHandler) GetPublicMenus(c *gin.Context) {
	tenantID := c.GetString("tenantId")
	if tenantID == "" {
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Tenant context tidak ditemukan", http.StatusBadRequest, nil))
		return
	}

	categoryID := c.Query("categoryId")
	search := c.Query("search")
	isAvailableStr := c.Query("isAvailable")

	query := h.db.WithContext(c.Request.Context()).Table("menus").
		Where("tenant_id = ?", tenantID).
		Where("deleted_at IS NULL")

	// Public customers should see available menus by default
	if isAvailableStr != "" {
		if isAvailableStr == "true" {
			query = query.Where("is_available = true")
		} else if isAvailableStr == "false" {
			query = query.Where("is_available = false")
		}
	} else {
		query = query.Where("is_available = true")
	}

	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchPattern, searchPattern)
	}

	var menus []menudto.MenuResponse
	err := query.Select("id::text AS id, tenant_id::text AS tenant_id, category_id::text AS category_id, name, description, price, image_url, is_available, created_at, updated_at").
		Order("name ASC").
		Scan(&menus).Error

	if err != nil {
		apperrors.Write(c, apperrors.New("INTERNAL_ERROR", "Gagal mengambil daftar menu", http.StatusInternalServerError, err.Error()))
		return
	}

	response.Success(c, http.StatusOK, "Menu list fetched successfully", menus)
}
