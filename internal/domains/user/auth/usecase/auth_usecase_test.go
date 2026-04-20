package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/dhegas/saas_gangsta/internal/config"
	"github.com/dhegas/saas_gangsta/internal/domains/user/auth"
	"github.com/dhegas/saas_gangsta/internal/domains/user/auth/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/user/auth/dto"
	"github.com/dhegas/saas_gangsta/internal/domains/user/auth/repository"
	"golang.org/x/crypto/bcrypt"
)

type mockAuthRepo struct {
	userByEmail     *domain.User
	userByID        *domain.User
	createdTenant   *domain.MerchantTenant
	merchantTenants []domain.MerchantTenant
	createErr       error
	createTenantErr error
	listTenantsErr  error
	errEmail        error
	errUserID       error
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

func (m *mockAuthRepo) CreateTenantForMerchant(_ context.Context, _ repository.CreateMerchantTenantInput) (*domain.MerchantTenant, error) {
	if m.createTenantErr != nil {
		return nil, m.createTenantErr
	}
	return m.createdTenant, nil
}

func (m *mockAuthRepo) ListTenantsByMerchant(_ context.Context, _ string) ([]domain.MerchantTenant, error) {
	if m.listTenantsErr != nil {
		return nil, m.listTenantsErr
	}
	return m.merchantTenants, nil
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
			Role:         "MITRA",
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
	if res.User.Role != "MITRA" {
		t.Fatalf("expected role MITRA, got %s", res.User.Role)
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
			Role:         "MITRA",
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
			Role:         "MITRA",
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
			Role:         "MITRA",
			IsActive:     true,
			TenantStatus: "active",
		},
	}

	uc := newAuthUsecaseForTest(repo)
	refreshToken, err := auth.GenerateRefreshToken(
		"u-1",
		"MITRA",
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
			Role:         "MITRA",
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

func TestCreateMerchantTenantSuccess(t *testing.T) {
	repo := &mockAuthRepo{
		createdTenant: &domain.MerchantTenant{ID: "t-1", Name: "Warung A", Slug: "warung-a", Status: "active", IsOwner: true},
	}
	uc := newAuthUsecaseForTest(repo)

	res, err := uc.CreateMerchantTenant(context.Background(), "u-1", dto.CreateMerchantTenantRequest{Name: "Warung A"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if res.Tenant.ID == "" || !res.Tenant.IsOwner {
		t.Fatalf("expected created tenant with owner flag")
	}
}

func TestListMerchantTenantsSuccess(t *testing.T) {
	repo := &mockAuthRepo{
		userByID:        &domain.User{ID: "u-1", Role: "MITRA", IsActive: true},
		merchantTenants: []domain.MerchantTenant{{ID: "t-1", Name: "Warung A", Slug: "warung-a", Status: "active", IsOwner: true}},
	}
	uc := newAuthUsecaseForTest(repo)

	res, err := uc.ListMerchantTenants(context.Background(), "u-1")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(res.Tenants) != 1 {
		t.Fatalf("expected 1 tenant, got %d", len(res.Tenants))
	}
}
