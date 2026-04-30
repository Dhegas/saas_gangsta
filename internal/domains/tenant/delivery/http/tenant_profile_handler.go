package http

import (
	"errors"
	"net/http"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/common/response"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant/dto"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// TenantProfileHandler menangani endpoint tenant profile.
type TenantProfileHandler struct {
	usecase domain.TenantProfileUsecase
}

func NewTenantProfileHandler(usecase domain.TenantProfileUsecase) *TenantProfileHandler {
	return &TenantProfileHandler{usecase: usecase}
}

// Create godoc
// @Summary      Buat tenant profile
// @Tags         tenant-profiles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body dto.CreateTenantProfileRequest true "Payload"
// @Success      201 {object} map[string]interface{}
// @Router       /tenant-profiles [post]
func (h *TenantProfileHandler) Create(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		apperrors.Write(c, apperrors.New("TENANT_NOT_FOUND", "Tenant context is required", http.StatusUnauthorized, nil))
		return
	}

	var req dto.CreateTenantProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Write(c, validationError(err))
		return
	}

	res, err := h.usecase.CreateProfile(c.Request.Context(), tenantID, req)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Tenant profile berhasil dibuat", res)
}

// List godoc
// @Summary      Daftar tenant profile
// @Tags         tenant-profiles
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]interface{}
// @Router       /tenant-profiles [get]
func (h *TenantProfileHandler) List(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		apperrors.Write(c, apperrors.New("TENANT_NOT_FOUND", "Tenant context is required", http.StatusUnauthorized, nil))
		return
	}

	res, err := h.usecase.ListProfiles(c.Request.Context(), tenantID)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Daftar tenant profile berhasil diambil", res)
}

// GetByID godoc
// @Summary      Detail tenant profile
// @Tags         tenant-profiles
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Profile ID"
// @Success      200 {object} map[string]interface{}
// @Router       /tenant-profiles/{id} [get]
func (h *TenantProfileHandler) GetByID(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		apperrors.Write(c, apperrors.New("TENANT_NOT_FOUND", "Tenant context is required", http.StatusUnauthorized, nil))
		return
	}

	profileID := c.Param("id")

	res, err := h.usecase.GetProfileByID(c.Request.Context(), tenantID, profileID)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Detail tenant profile berhasil diambil", res)
}

// Update godoc
// @Summary      Update tenant profile
// @Tags         tenant-profiles
// @Accept       json
// @Security     BearerAuth
// @Produce      json
// @Param        id   path string                          true  "Profile ID"
// @Param        body body dto.UpdateTenantProfileRequest true  "Payload"
// @Success      200 {object} map[string]interface{}
// @Router       /tenant-profiles/{id} [put]
func (h *TenantProfileHandler) Update(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		apperrors.Write(c, apperrors.New("TENANT_NOT_FOUND", "Tenant context is required", http.StatusUnauthorized, nil))
		return
	}

	profileID := c.Param("id")

	var req dto.UpdateTenantProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperrors.Write(c, validationError(err))
		return
	}

	res, err := h.usecase.UpdateProfile(c.Request.Context(), tenantID, profileID, req)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Tenant profile berhasil diupdate", res)
}

// Delete godoc
// @Summary      Hapus tenant profile (soft delete)
// @Tags         tenant-profiles
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Profile ID"
// @Success      200 {object} map[string]interface{}
// @Router       /tenant-profiles/{id} [delete]
func (h *TenantProfileHandler) Delete(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		apperrors.Write(c, apperrors.New("TENANT_NOT_FOUND", "Tenant context is required", http.StatusUnauthorized, nil))
		return
	}

	profileID := c.Param("id")

	if err := h.usecase.DeleteProfile(c.Request.Context(), tenantID, profileID); err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Tenant profile berhasil dihapus", nil)
}

// ToggleActive godoc
// @Summary      Aktif / nonaktifkan tenant profile
// @Tags         tenant-profiles
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Profile ID"
// @Success      200 {object} map[string]interface{}
// @Router       /tenant-profiles/{id}/toggle-active [patch]
func (h *TenantProfileHandler) ToggleActive(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		apperrors.Write(c, apperrors.New("TENANT_NOT_FOUND", "Tenant context is required", http.StatusUnauthorized, nil))
		return
	}

	profileID := c.Param("id")

	res, err := h.usecase.ToggleActive(c.Request.Context(), tenantID, profileID)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Status tenant profile berhasil diubah", res)
}

// --- helpers ---

func validationError(err error) *apperrors.AppError {
	var validationErrs validator.ValidationErrors
	details := err.Error()
	if errors.As(err, &validationErrs) {
		details = validationErrs.Error()
	}
	return apperrors.New("VALIDATION_ERROR", "Payload tenant profile tidak valid", http.StatusBadRequest, details)
}
