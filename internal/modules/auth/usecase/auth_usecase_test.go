package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/dhegas/saas_gangsta/internal/common/auth"
	"github.com/dhegas/saas_gangsta/internal/common/config"
	"github.com/dhegas/saas_gangsta/internal/modules/auth/domain"
	"github.com/dhegas/saas_gangsta/internal/modules/auth/dto"
	"golang.org/x/crypto/bcrypt"
)

type mockAuthRepo struct {
	userByEmail *domain.User
	userByID    *domain.User
	createErr   error
	errEmail    error
	errUserID   error
}

func (m *mockAuthRepo) FindByEmail(_ context.Context, _ string) (*domain.User, error) {
	return m.userByEmail, m.errEmail
}

func (m *mockAuthRepo) FindByID(_ context.Context, _ string) (*domain.User, error) {
	return m.userByID, m.errUserID
}

func (m *mockAuthRepo) CreateUser(_ context.Context, user *domain.User) error {
	if m.createErr != nil {
		return m.createErr
	}
	if user.ID == "" {
		user.ID = "created-user-id"
	}
	return nil
}

func mustHash(t *testing.T, plain string) string {
	t.Helper()
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	return string(h)
}

func newAuthUsecaseForTest(repo *mockAuthRepo) AuthUsecase {
	cfg := &config.Config{
		JWTSecret:             "12345678901234567890123456789012",
		JWTAccessTokenExpiry:  15 * time.Minute,
		JWTRefreshTokenExpiry: 7 * 24 * time.Hour,
	}
	return NewAuthUsecase(repo, cfg)
}

func TestLoginSuccess(t *testing.T) {
	repo := &mockAuthRepo{
		userByEmail: &domain.User{
			ID:           "u-1",
			TenantID:     "t-1",
			Email:        "merchant@test.local",
			PasswordHash: mustHash(t, "secret123"),
			Role:         "merchant",
			IsActive:     true,
			TenantStatus: "active",
		},
	}

	uc := newAuthUsecaseForTest(repo)
	res, err := uc.Login(context.Background(), dto.LoginRequest{Email: "merchant@test.local", Password: "secret123"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if res.AccessToken == "" || res.RefreshToken == "" {
		t.Fatalf("expected tokens to be generated")
	}
	if res.User.Role != "merchant" {
		t.Fatalf("expected role merchant, got %s", res.User.Role)
	}
}

func TestRegisterSuccess(t *testing.T) {
	repo := &mockAuthRepo{}
	uc := newAuthUsecaseForTest(repo)

	res, err := uc.Register(context.Background(), dto.RegisterRequest{
		Email:    "new@test.local",
		Password: "secret123",
		FullName: "New User",
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if res.User.ID == "" {
		t.Fatalf("expected created user id")
	}
}

func TestRegisterConflictEmail(t *testing.T) {
	repo := &mockAuthRepo{
		userByEmail: &domain.User{ID: "u-1", Email: "exists@test.local"},
	}
	uc := newAuthUsecaseForTest(repo)

	if _, err := uc.Register(context.Background(), dto.RegisterRequest{
		Email:    "exists@test.local",
		Password: "secret123",
		FullName: "Exists",
	}); err == nil {
		t.Fatalf("expected conflict error")
	}
}

func TestLoginUnauthorizedInvalidPassword(t *testing.T) {
	repo := &mockAuthRepo{
		userByEmail: &domain.User{
			ID:           "u-1",
			TenantID:     "t-1",
			Email:        "merchant@test.local",
			PasswordHash: mustHash(t, "secret123"),
			Role:         "merchant",
			IsActive:     true,
			TenantStatus: "active",
		},
	}

	uc := newAuthUsecaseForTest(repo)
	if _, err := uc.Login(context.Background(), dto.LoginRequest{Email: "merchant@test.local", Password: "wrongpass"}); err == nil {
		t.Fatalf("expected unauthorized error")
	}
}

func TestLoginTenantInactive(t *testing.T) {
	repo := &mockAuthRepo{
		userByEmail: &domain.User{
			ID:           "u-1",
			TenantID:     "t-1",
			Email:        "merchant@test.local",
			PasswordHash: mustHash(t, "secret123"),
			Role:         "merchant",
			IsActive:     true,
			TenantStatus: "inactive",
		},
	}

	uc := newAuthUsecaseForTest(repo)
	if _, err := uc.Login(context.Background(), dto.LoginRequest{Email: "merchant@test.local", Password: "secret123"}); err == nil {
		t.Fatalf("expected tenant inactive error")
	}
}

func TestRefreshSuccess(t *testing.T) {
	repo := &mockAuthRepo{
		userByID: &domain.User{
			ID:           "u-1",
			TenantID:     "t-1",
			Email:        "merchant@test.local",
			Role:         "merchant",
			IsActive:     true,
			TenantStatus: "active",
		},
	}

	uc := newAuthUsecaseForTest(repo)
	refreshToken, err := auth.GenerateRefreshToken(
		"u-1",
		"merchant",
		"t-1",
		7*24*time.Hour,
		"12345678901234567890123456789012",
	)
	if err != nil {
		t.Fatalf("generate refresh token: %v", err)
	}

	res, err := uc.Refresh(context.Background(), dto.RefreshTokenRequest{RefreshToken: refreshToken})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if res.AccessToken == "" || res.RefreshToken == "" {
		t.Fatalf("expected refreshed tokens")
	}
}

func TestMeSuccess(t *testing.T) {
	repo := &mockAuthRepo{
		userByID: &domain.User{
			ID:           "u-1",
			TenantID:     "t-1",
			Email:        "merchant@test.local",
			Role:         "merchant",
			IsActive:     true,
			TenantStatus: "active",
		},
	}

	uc := newAuthUsecaseForTest(repo)
	res, err := uc.Me(context.Background(), "u-1")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if res.User.Email != "merchant@test.local" {
		t.Fatalf("unexpected email: %s", res.User.Email)
	}
}

func TestLogoutSuccess(t *testing.T) {
	uc := newAuthUsecaseForTest(&mockAuthRepo{})
	if err := uc.Logout(context.Background(), "u-1", dto.LogoutRequest{}); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}
