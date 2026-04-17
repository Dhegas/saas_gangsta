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
