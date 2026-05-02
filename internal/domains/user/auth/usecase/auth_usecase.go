package usecase

import (
	"context"
	"errors"
	"net/http"
	"strings"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/config"
	commonauth "github.com/dhegas/saas_gangsta/internal/domains/user/auth"
	"github.com/dhegas/saas_gangsta/internal/domains/user/auth/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/user/auth/dto"
	"github.com/dhegas/saas_gangsta/internal/domains/user/auth/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	CreatePartnerTenant(ctx context.Context, userID string, req dto.CreatePartnerTenantRequest) (*dto.CreatePartnerTenantResponse, error)
	ListPartnerTenants(ctx context.Context, userID string) (*dto.ListPartnerTenantsResponse, error)
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

	if err := validateTenantState(user); err != nil {
		return nil, err
	}

	return u.buildLoginResponse(user)
}

func (u *authUsecase) CreatePartnerTenant(ctx context.Context, userID string, req dto.CreatePartnerTenantRequest) (*dto.CreatePartnerTenantResponse, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, apperrors.New("UNAUTHORIZED", "User tidak valid", http.StatusUnauthorized, nil)
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, apperrors.New("VALIDATION_ERROR", "Nama tenant wajib diisi", http.StatusBadRequest, nil)
	}

	tenant, err := u.repo.CreateTenantForPartner(ctx, repository.CreatePartnerTenantInput{
		UserID: userID,
		Name:   name,
	})
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			return nil, apperrors.New("NOT_FOUND", "User tidak ditemukan", http.StatusNotFound, nil)
		case errors.Is(err, repository.ErrUserNotPartner):
			return nil, apperrors.New("FORBIDDEN", "Hanya PARTNER yang dapat membuat tenant", http.StatusForbidden, nil)
		case errors.Is(err, repository.ErrPartnerSubscriptionMissing):
			return nil, apperrors.New("FORBIDDEN", "Subscription PARTNER tidak ditemukan", http.StatusForbidden, nil)
		case errors.Is(err, repository.ErrTenantLimitReached):
			return nil, apperrors.New("FORBIDDEN", "Batas jumlah tenant pada paket subscription sudah tercapai", http.StatusForbidden, nil)
		default:
			return nil, apperrors.New("INTERNAL_ERROR", "Gagal membuat tenant partner", http.StatusInternalServerError, nil)
		}
	}

	return &dto.CreatePartnerTenantResponse{
		Tenant: dto.PartnerTenantResponse{
			ID:      tenant.ID,
			Name:    tenant.Name,
			Slug:    tenant.Slug,
			Status:  tenant.Status,
			IsOwner: tenant.IsOwner,
		},
	}, nil
}

func (u *authUsecase) ListPartnerTenants(ctx context.Context, userID string) (*dto.ListPartnerTenantsResponse, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, apperrors.New("UNAUTHORIZED", "User tidak valid", http.StatusUnauthorized, nil)
	}

	partner, err := u.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil data user", http.StatusInternalServerError, nil)
	}
	if partner == nil {
		return nil, apperrors.New("UNAUTHORIZED", "User tidak ditemukan", http.StatusUnauthorized, nil)
	}
	if partner.Role != "PARTNER" {
		return nil, apperrors.New("FORBIDDEN", "Hanya PARTNER yang dapat melihat tenant", http.StatusForbidden, nil)
	}

	tenants, err := u.repo.ListTenantsByPartner(ctx, userID)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil daftar tenant partner", http.StatusInternalServerError, nil)
	}

	items := make([]dto.PartnerTenantResponse, 0, len(tenants))
	for _, tenant := range tenants {
		items = append(items, dto.PartnerTenantResponse{
			ID:      tenant.ID,
			Name:    tenant.Name,
			Slug:    tenant.Slug,
			Status:  tenant.Status,
			IsOwner: tenant.IsOwner,
		})
	}

	return &dto.ListPartnerTenantsResponse{Tenants: items}, nil
}

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

func (u *authUsecase) Logout(_ context.Context, userID string, _ dto.LogoutRequest) error {
	if strings.TrimSpace(userID) == "" {
		return apperrors.New("UNAUTHORIZED", "User tidak valid", http.StatusUnauthorized, nil)
	}

	// Stateless JWT logout: client removes token locally.
	return nil
}

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

func validateTenantState(user *domain.User) error {
	if user.Role == "PARTNER" {
		if user.TenantID != "" && strings.TrimSpace(user.TenantStatus) != "active" {
			return apperrors.New("TENANT_INACTIVE", "Tenant tidak aktif", http.StatusForbidden, nil)
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
