package usecase

import (
	"context"

	"github.com/dhegas/saas_gangsta/internal/config"
	"github.com/dhegas/saas_gangsta/internal/domains/user/auth/dto"
	"github.com/dhegas/saas_gangsta/internal/domains/user/auth/repository"
)

type AuthUsecase interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	Refresh(ctx context.Context, req dto.RefreshTokenRequest) (*dto.LoginResponse, error)
	Logout(ctx context.Context, userID string, req dto.LogoutRequest) error
	Me(ctx context.Context, userID string) (*dto.MeResponse, error)
}

type authUsecase struct {
	repo repository.AuthRepository
	cfg  *config.Config
}

func NewAuthUsecase(repo repository.AuthRepository, cfg *config.Config) AuthUsecase {
	return &authUsecase{repo: repo, cfg: cfg}
}
