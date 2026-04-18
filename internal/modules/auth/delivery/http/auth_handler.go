package http

import (
	"errors"
	"net/http"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/common/response"
	"github.com/dhegas/saas_gangsta/internal/modules/auth/dto"
	"github.com/dhegas/saas_gangsta/internal/modules/auth/usecase"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	usecase usecase.AuthUsecase
}

func NewAuthHandler(usecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{usecase: usecase}
}

// Register godoc
// @Summary Register user account
// @Description Membuat akun baru untuk login dengan role default customer
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register payload"
// @Success 201 {object} response.Envelope
// @Failure 400 {object} response.Envelope
// @Failure 409 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var validationErrs validator.ValidationErrors
		details := err.Error()
		if errors.As(err, &validationErrs) {
			details = validationErrs.Error()
		}
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Payload registrasi tidak valid", http.StatusBadRequest, details))
		return
	}

	res, err := h.usecase.Register(c.Request.Context(), req)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Akun berhasil dibuat", res)
}

// Login godoc
// @Summary Login user
// @Description Login untuk semua role (customer, merchant, admin)
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login payload"
// @Success 200 {object} response.Envelope
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var validationErrs validator.ValidationErrors
		details := err.Error()
		if errors.As(err, &validationErrs) {
			details = validationErrs.Error()
		}
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Payload login tidak valid", http.StatusBadRequest, details))
		return
	}

	res, err := h.usecase.Login(c.Request.Context(), req)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Login berhasil", res)
}

// Subscribe godoc
// @Summary Subscribe customer to merchant plan
// @Description Customer subscribe paket untuk upgrade akun menjadi merchant; tenant dibuat kemudian di dashboard merchant
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.SubscribeRequest true "Subscribe payload"
// @Success 200 {object} response.Envelope
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 409 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /auth/subscribe [post]
func (h *AuthHandler) Subscribe(c *gin.Context) {
	userID, _ := c.Get("userId")
	userIDStr, _ := userID.(string)

	var req dto.SubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var validationErrs validator.ValidationErrors
		details := err.Error()
		if errors.As(err, &validationErrs) {
			details = validationErrs.Error()
		}
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Payload subscribe tidak valid", http.StatusBadRequest, details))
		return
	}

	res, err := h.usecase.Subscribe(c.Request.Context(), userIDStr, req)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Subscribe berhasil, akun di-upgrade menjadi merchant", res)
}

// CreateMerchantTenant godoc
// @Summary Create merchant tenant
// @Description Merchant membuat tenant baru miliknya sesuai limit paket subscription
// @Tags Merchant
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateMerchantTenantRequest true "Create merchant tenant payload"
// @Success 201 {object} response.Envelope
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /merchant/tenants [post]
func (h *AuthHandler) CreateMerchantTenant(c *gin.Context) {
	userID, _ := c.Get("userId")
	userIDStr, _ := userID.(string)

	var req dto.CreateMerchantTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var validationErrs validator.ValidationErrors
		details := err.Error()
		if errors.As(err, &validationErrs) {
			details = validationErrs.Error()
		}
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Payload create tenant tidak valid", http.StatusBadRequest, details))
		return
	}

	res, err := h.usecase.CreateMerchantTenant(c.Request.Context(), userIDStr, req)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Tenant merchant berhasil dibuat", res)
}

// ListMerchantTenants godoc
// @Summary List merchant tenants
// @Description Ambil daftar tenant milik merchant login
// @Tags Merchant
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /merchant/tenants [get]
func (h *AuthHandler) ListMerchantTenants(c *gin.Context) {
	userID, _ := c.Get("userId")
	userIDStr, _ := userID.(string)

	res, err := h.usecase.ListMerchantTenants(c.Request.Context(), userIDStr)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Daftar tenant merchant berhasil diambil", res)
}

// Refresh godoc
// @Summary Refresh token
// @Description Refresh access token menggunakan refresh token yang valid
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh payload"
// @Success 200 {object} response.Envelope
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var validationErrs validator.ValidationErrors
		details := err.Error()
		if errors.As(err, &validationErrs) {
			details = validationErrs.Error()
		}
		apperrors.Write(c, apperrors.New("VALIDATION_ERROR", "Payload refresh token tidak valid", http.StatusBadRequest, details))
		return
	}

	res, err := h.usecase.Refresh(c.Request.Context(), req)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Refresh token berhasil", res)
}

// Logout godoc
// @Summary Logout user
// @Description Logout user dari sesi aktif
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.LogoutRequest false "Logout payload"
// @Success 200 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, _ := c.Get("userId")
	userIDStr, _ := userID.(string)

	var req dto.LogoutRequest
	_ = c.ShouldBindJSON(&req)

	if err := h.usecase.Logout(c.Request.Context(), userIDStr, req); err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Logout berhasil", gin.H{"loggedOut": true})
}

// Me godoc
// @Summary Current user info
// @Description Ambil profil user yang sedang login
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get("userId")
	userIDStr, _ := userID.(string)

	res, err := h.usecase.Me(c.Request.Context(), userIDStr)
	if err != nil {
		apperrors.Write(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Profil user berhasil diambil", res)
}
