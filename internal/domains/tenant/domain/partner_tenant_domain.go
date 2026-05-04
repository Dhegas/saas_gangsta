package domain

import (
	"context"

	"github.com/dhegas/saas_gangsta/internal/domains/tenant/dto"
)

type PartnerTenantUsecase interface {
	CreatePartnerTenant(ctx context.Context, userID string, req dto.CreatePartnerTenantRequest) (*dto.CreatePartnerTenantResponse, error)
	ListPartnerTenants(ctx context.Context, userID string) (*dto.ListPartnerTenantsResponse, error)
}

type PartnerTenantRepository interface {
	FindPartnerByID(ctx context.Context, userID string) (*PartnerUser, error)
	CreateTenantForPartner(ctx context.Context, input CreatePartnerTenantInput) (*PartnerTenant, error)
	ListTenantsByPartner(ctx context.Context, userID string) ([]PartnerTenant, error)
}

type PartnerUser struct {
	ID       string
	Role     string
	IsActive bool
}

type PartnerTenant struct {
	ID      string
	Name    string
	Slug    string
	Status  string
	IsOwner bool
}

type CreatePartnerTenantInput struct {
	UserID string
	Name   string
}
