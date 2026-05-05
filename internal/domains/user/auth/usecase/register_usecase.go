package usecase

import (
	"context"
	"net/http"
	"strings"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/domains/user/auth/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/user/auth/dto"
	"golang.org/x/crypto/bcrypt"
)

func (u *authUsecase) Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	email := strings.TrimSpace(req.Email)
	password := strings.TrimSpace(req.Password)
	fullName := strings.TrimSpace(req.FullName)

	existing, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal memproses registrasi", http.StatusInternalServerError, nil)
	}
	if existing != nil {
		return nil, apperrors.New("CONFLICT", "Email sudah terdaftar", http.StatusConflict, nil)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal menyiapkan akun", http.StatusInternalServerError, nil)
	}

	user := &domain.User{
		TenantID:     "",
		Email:        email,
		FullName:     fullName,
		PasswordHash: string(passwordHash),
		Role:         "CUSTOMER",
		IsActive:     true,
		TenantStatus: "active",
	}

	if err := u.repo.CreateUser(ctx, user); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate email") {
			return nil, apperrors.New("CONFLICT", "Email sudah terdaftar", http.StatusConflict, nil)
		}
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal menyimpan akun", http.StatusInternalServerError, nil)
	}

	return &dto.RegisterResponse{
		User: dto.UserResponse{
			ID:       user.ID,
			Email:    user.Email,
			Role:     user.Role,
			TenantID: user.TenantID,
		},
	}, nil
}
