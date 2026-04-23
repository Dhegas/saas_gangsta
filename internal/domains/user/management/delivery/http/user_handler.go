package http

import (
	"errors"
	"net/http"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/common/response"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant"
	"github.com/dhegas/saas_gangsta/internal/domains/user/management/dto"
	"github.com/dhegas/saas_gangsta/internal/domains/user/management/usecase"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	usecase usecase.UserUsecase
}

func NewUserHandler(usecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{usecase: usecase}
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		apperrors.Write(c, apperrors.New("TENANT_NOT_FOUND", "Tenant context is required", http.StatusUnauthorized, nil))
		return
	}

	res, err := h.usecase.ListUsersByTenant(c.Request.Context(), tenantID)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Daftar user berhasil diambil", res)
}

func (h *UserHandler) GetUserDetail(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		apperrors.Write(c, apperrors.New("TENANT_NOT_FOUND", "Tenant context is required", http.StatusUnauthorized, nil))
		return
	}

	var uri dto.UserIDParam
	if err := c.ShouldBindUri(&uri); err != nil {
		var validationErrs validator.ValidationErrors
		details := err.Error()
		if errors.As(err, &validationErrs) {
			details = validationErrs.Error()
		}
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Parameter user id tidak valid", http.StatusBadRequest, details))
		return
	}

	res, err := h.usecase.GetUserDetailByTenant(c.Request.Context(), tenantID, uri.ID)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Detail user berhasil diambil", res)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		apperrors.Write(c, apperrors.New("TENANT_NOT_FOUND", "Tenant context is required", http.StatusUnauthorized, nil))
		return
	}

	var uri dto.UserIDParam
	if err := c.ShouldBindUri(&uri); err != nil {
		var validationErrs validator.ValidationErrors
		details := err.Error()
		if errors.As(err, &validationErrs) {
			details = validationErrs.Error()
		}
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Parameter user id tidak valid", http.StatusBadRequest, details))
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var validationErrs validator.ValidationErrors
		details := err.Error()
		if errors.As(err, &validationErrs) {
			details = validationErrs.Error()
		}
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Payload update user tidak valid", http.StatusBadRequest, details))
		return
	}

	res, err := h.usecase.UpdateUserByTenant(c.Request.Context(), tenantID, uri.ID, req)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "User berhasil diupdate", res)
}

func (h *UserHandler) SoftDeleteUser(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		apperrors.Write(c, apperrors.New("TENANT_NOT_FOUND", "Tenant context is required", http.StatusUnauthorized, nil))
		return
	}

	var uri dto.UserIDParam
	if err := c.ShouldBindUri(&uri); err != nil {
		var validationErrs validator.ValidationErrors
		details := err.Error()
		if errors.As(err, &validationErrs) {
			details = validationErrs.Error()
		}
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Parameter user id tidak valid", http.StatusBadRequest, details))
		return
	}

	res, err := h.usecase.SoftDeleteUserByTenant(c.Request.Context(), tenantID, uri.ID)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "User berhasil dihapus", res)
}

func (h *UserHandler) ToggleActiveUser(c *gin.Context) {
	tenantID, err := tenant.GetTenantID(c)
	if err != nil {
		apperrors.Write(c, apperrors.New("TENANT_NOT_FOUND", "Tenant context is required", http.StatusUnauthorized, nil))
		return
	}

	var uri dto.UserIDParam
	if err := c.ShouldBindUri(&uri); err != nil {
		var validationErrs validator.ValidationErrors
		details := err.Error()
		if errors.As(err, &validationErrs) {
			details = validationErrs.Error()
		}
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Parameter user id tidak valid", http.StatusBadRequest, details))
		return
	}

	res, err := h.usecase.ToggleUserActiveByTenant(c.Request.Context(), tenantID, uri.ID)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Status user berhasil diubah", res)
}
