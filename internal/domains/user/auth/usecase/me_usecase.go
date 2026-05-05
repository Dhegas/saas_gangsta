package usecase

import (
	"context"
	"net/http"
	"strings"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/domains/user/auth/dto"
)

func (u *authUsecase) Me(ctx context.Context, userID string) (*dto.MeResponse, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, apperrors.New("UNAUTHORIZED", "User tidak valid", http.StatusUnauthorized, nil)
	}

	user, err := u.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil profil user", http.StatusInternalServerError, nil)
	}
	if user == nil {
		return nil, apperrors.New("UNAUTHORIZED", "User tidak ditemukan", http.StatusUnauthorized, nil)
	}

	if err := validateTenantState(user); err != nil {
		return nil, err
	}

	return &dto.MeResponse{
		User: dto.UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			Role:     user.Role,
			TenantID: user.TenantID,
		},
	}, nil
}
