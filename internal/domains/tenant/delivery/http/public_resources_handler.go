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
// @Router /public/tenant/{tenantSlug}/categories [get]
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
// @Router /public/tenant/{tenantSlug}/menus [get]
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

// GetPublicTables godoc
// @Summary List active dining tables for a tenant
// @Description Mengambil daftar meja aktif dalam context tenant beserta status keterisiannya
// @Tags Public
// @Produce json
// @Param tenantSlug path string true "Tenant Slug"
// @Success 200 {object} response.Envelope
// @Failure 400 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /public/tenant/{tenantSlug}/tables [get]
func (h *PublicResourcesHandler) GetPublicTables(c *gin.Context) {
	tenantID := c.GetString("tenantId")
	if tenantID == "" {
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Tenant context tidak ditemukan", http.StatusBadRequest, nil))
		return
	}

	type PublicTableResponse struct {
		ID        string `json:"id" gorm:"column:id"`
		TenantID  string `json:"tenantId" gorm:"column:tenant_id"`
		TableName string `json:"tableName" gorm:"column:table_name"`
		Status    string `json:"status" gorm:"column:status"`
	}

	var tables []PublicTableResponse
	err := h.db.WithContext(c.Request.Context()).
		Table("dining_tables dt").
		Select(`
			dt.id::text AS id, 
			dt.tenant_id::text AS tenant_id, 
			dt.table_name,
			CASE 
				WHEN EXISTS (
					SELECT 1 FROM orders o 
					WHERE o.dining_tables_id = dt.id 
					  AND o.status NOT IN ('COMPLETED', 'CANCELLED') 
					  AND o.deleted_at IS NULL
				) THEN 'occupied'
				ELSE 'kosong'
			END AS status
		`).
		Where("dt.tenant_id = ?", tenantID).
		Where("dt.deleted_at IS NULL").
		Order("dt.table_name ASC").
		Scan(&tables).Error

	if err != nil {
		apperrors.Write(c, apperrors.New("INTERNAL_ERROR", "Gagal mengambil daftar meja", http.StatusInternalServerError, err.Error()))
		return
	}

	response.Success(c, http.StatusOK, "Table list fetched successfully", tables)
}

