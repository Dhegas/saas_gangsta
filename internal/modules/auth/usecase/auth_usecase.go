package usecase

import (
	"context"
	"net/http"
	"strings"

	commonauth "github.com/dhegas/saas_gangsta/internal/common/auth"
	"github.com/dhegas/saas_gangsta/internal/common/config"
	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/modules/auth/dto"
	"github.com/dhegas/saas_gangsta/internal/modules/auth/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
}

type authUsecase struct {
	repo repository.AuthRepository
	cfg  *config.Config
}

func NewAuthUsecase(repo repository.AuthRepository, cfg *config.Config) AuthUsecase {
	return &authUsecase{repo: repo, cfg: cfg}
}

func (u *authUsecase) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	email := strings.TrimSpace(req.Email)
	password := strings.TrimSpace(req.Password)

	user, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal memproses login", http.StatusInternalServerError, nil)
	}

	if user == nil || !user.IsActive {
		return nil, apperrors.New("UNAUTHORIZED", "Email atau password salah", http.StatusUnauthorized, nil)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, apperrors.New("UNAUTHORIZED", "Email atau password salah", http.StatusUnauthorized, nil)
	}

	if user.Role != "admin" {
		if user.TenantID == "" {
			return nil, apperrors.New("TENANT_NOT_FOUND", "Tenant tidak ditemukan", http.StatusNotFound, nil)
		}
		if strings.TrimSpace(user.TenantStatus) != "active" {
			return nil, apperrors.New("TENANT_INACTIVE", "Tenant tidak aktif", http.StatusForbidden, nil)
		}
	}

	accessToken, err := commonauth.GenerateAccessToken(
		user.ID,
		user.Role,
		user.TenantID,
		u.cfg.JWTAccessTokenExpiry,
		u.cfg.JWTSecret,
	)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal membuat access token", http.StatusInternalServerError, nil)
	}

	refreshToken, err := commonauth.GenerateRefreshToken(
		user.ID,
		user.Role,
		user.TenantID,
		u.cfg.JWTRefreshTokenExpiry,
		u.cfg.JWTSecret,
	)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal membuat refresh token", http.StatusInternalServerError, nil)
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
