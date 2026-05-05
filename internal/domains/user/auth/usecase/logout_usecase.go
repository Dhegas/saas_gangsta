package usecase

import (
	"context"
	"net/http"
	"strings"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/domains/user/auth/dto"
)

func (u *authUsecase) Logout(_ context.Context, userID string, _ dto.LogoutRequest) error {
	if strings.TrimSpace(userID) == "" {
		return apperrors.New("UNAUTHORIZED", "User tidak valid", http.StatusUnauthorized, nil)
	}

	// Stateless JWT logout: client removes token locally.
	return nil
}
