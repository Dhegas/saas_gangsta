package usecase

import (
	"context"
	"net/http"
	"strings"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/domains/user/auth/dto"
	"golang.org/x/crypto/bcrypt"
)

func (u *authUsecase) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	email := strings.TrimSpace(req.Email)
	password := strings.TrimSpace(req.Password)

	user, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal memproses login", http.StatusInternalServerError)
	}

	if user == nil || !user.IsActive {
		return nil, apperrors.New("UNAUTHORIZED", "Email atau password salah", http.StatusUnauthorized)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, apperrors.New("UNAUTHORIZED", "Email atau password salah", http.StatusUnauthorized)
	}

	if err := validateTenantState(user); err != nil {
		return nil, err
	}

	return u.buildLoginResponse(user)
}
