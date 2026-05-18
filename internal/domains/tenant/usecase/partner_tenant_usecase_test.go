package usecase

import (
	"context"
	"testing"

	"github.com/dhegas/saas_gangsta/internal/domains/tenant/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant/dto"
)

type mockPartnerTenantRepo struct {
	partner        *domain.PartnerUser
	createdTenant  *domain.PartnerTenant
	partnerTenants []domain.PartnerTenant
	findErr        error
	createErr      error
	listErr        error
}

func (m *mockPartnerTenantRepo) FindPartnerByID(_ context.Context, _ string) (*domain.PartnerUser, error) {
	return m.partner, m.findErr
}

func (m *mockPartnerTenantRepo) CreateTenantForPartner(_ context.Context, _ domain.CreatePartnerTenantInput) (*domain.PartnerTenant, error) {
	return m.createdTenant, m.createErr
}

func (m *mockPartnerTenantRepo) ListTenantsByPartner(_ context.Context, _ string) ([]domain.PartnerTenant, error) {
	return m.partnerTenants, m.listErr
}

func (m *mockPartnerTenantRepo) SoftDeleteTenant(_ context.Context, _ string, _ string) error {
	return nil
}

func TestCreatePartnerTenantSuccess(t *testing.T) {
	repo := &mockPartnerTenantRepo{
		createdTenant: &domain.PartnerTenant{ID: "t-1", Name: "Warung A", Slug: "warung-a", Status: "active", IsOwner: true},
	}
	uc := NewPartnerTenantUsecase(repo)

	res, err := uc.CreatePartnerTenant(context.Background(), "u-1", dto.CreatePartnerTenantRequest{Name: "Warung A"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if res.Tenant.ID == "" || !res.Tenant.IsOwner {
		t.Fatalf("expected created tenant with owner flag")
	}
}

func TestListPartnerTenantsSuccess(t *testing.T) {
	repo := &mockPartnerTenantRepo{
		partner:        &domain.PartnerUser{ID: "u-1", Role: "PARTNER", IsActive: true},
		partnerTenants: []domain.PartnerTenant{{ID: "t-1", Name: "Warung A", Slug: "warung-a", Status: "active", IsOwner: true}},
	}
	uc := NewPartnerTenantUsecase(repo)

	res, err := uc.ListPartnerTenants(context.Background(), "u-1")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(res.Tenants) != 1 {
		t.Fatalf("expected 1 tenant, got %d", len(res.Tenants))
	}
}
