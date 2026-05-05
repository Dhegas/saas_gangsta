package usecase

import (
	"context"
	"net/http"
	"strings"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	commonauth "github.com/dhegas/saas_gangsta/internal/domains/user/auth"
	"github.com/dhegas/saas_gangsta/internal/domains/user/auth/dto"
)

func (u *authUsecase) Refresh(ctx context.Context, req dto.RefreshTokenRequest) (*dto.LoginResponse, error) {
	refreshToken := strings.TrimSpace(req.RefreshToken)
	if refreshToken == "" {
		return nil, apperrors.New("VALIDATION_ERROR", "Refresh token wajib diisi", http.StatusBadRequest, nil)
	}

	claims, err := commonauth.ParseRefreshToken(refreshToken, u.cfg.JWTSecret)
	if err != nil {
		return nil, apperrors.New("UNAUTHORIZED", "Refresh token tidak valid atau expired", http.StatusUnauthorized, nil)
	}

	user, err := u.repo.FindByID(ctx, claims.Subject)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal memproses refresh token", http.StatusInternalServerError, nil)
	}
	if user == nil || !user.IsActive {
		return nil, apperrors.New("UNAUTHORIZED", "User tidak valid", http.StatusUnauthorized, nil)
	}

	if err := validateTenantState(user); err != nil {
		return nil, err
	}

	return u.buildLoginResponse(user)
}
