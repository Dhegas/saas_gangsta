package usecase

import (
	"net/http"
	"strings"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	commonauth "github.com/dhegas/saas_gangsta/internal/domains/user/auth"
	"github.com/dhegas/saas_gangsta/internal/domains/user/auth/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/user/auth/dto"
)

func validateTenantState(user *domain.User) error {
	if user.Role == "PARTNER" {
		if user.TenantID != "" && strings.TrimSpace(user.TenantStatus) != "active" {
			return apperrors.New("TENANT_INACTIVE", "Tenant tidak aktif", http.StatusForbidden)
		}
	}

	return nil
}

func (u *authUsecase) buildLoginResponse(user *domain.User) (*dto.LoginResponse, error) {
	accessToken, err := commonauth.GenerateAccessToken(
		user.ID,
		user.Role,
		user.TenantID,
		u.cfg.JWTAccessTokenExpiry,
		u.cfg.JWTSecret,
	)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal membuat access token", http.StatusInternalServerError)
	}

	refreshToken, err := commonauth.GenerateRefreshToken(
		user.ID,
		user.Role,
		user.TenantID,
		u.cfg.JWTRefreshTokenExpiry,
		u.cfg.JWTSecret,
	)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal membuat refresh token", http.StatusInternalServerError)
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			Role:     user.Role,
			TenantID: user.TenantID,
		},
	}, nil
}
