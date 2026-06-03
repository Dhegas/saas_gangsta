package domain

import (
	"context"

	"github.com/dhegas/saas_gangsta/internal/domains/tenant/dto"
)

type PartnerTenantUsecase interface {
	CreatePartnerTenant(ctx context.Context, userID string, req dto.CreatePartnerTenantRequest) (*dto.CreatePartnerTenantResponse, error)
	ListPartnerTenants(ctx context.Context, userID string) (*dto.ListPartnerTenantsResponse, error)
	SoftDeletePartnerTenant(ctx context.Context, userID string, tenantID string) error
	GetPartnerTenantByID(ctx context.Context, userID string, tenantID string) (*dto.PartnerTenantResponse, error)
	UpdatePartnerTenant(ctx context.Context, userID string, tenantID string, req dto.UpdatePartnerTenantRequest) (*dto.CreatePartnerTenantResponse, error)
}

type PartnerTenantRepository interface {
	FindPartnerByID(ctx context.Context, userID string) (*PartnerUser, error)
	CreateTenantForPartner(ctx context.Context, input CreatePartnerTenantInput) (*PartnerTenant, error)
	ListTenantsByPartner(ctx context.Context, userID string) ([]PartnerTenant, error)
	SoftDeleteTenant(ctx context.Context, userID string, tenantID string) error
	GetTenantByIDAndPartner(ctx context.Context, userID string, tenantID string) (*PartnerTenant, error)
	UpdateTenant(ctx context.Context, userID string, tenantID string, name string, description string, address string, phoneNumber string) (*PartnerTenant, error)
}

type PartnerUser struct {
	ID       string
	Role     string
	IsActive bool
}

type PartnerTenant struct {
	ID          string
	Name        string
	Slug        string
	Status      string
	Description string
	Address     string
	PhoneNumber string
	OpenHours   string
	LogoURL     string
	BannerURL   string
	IsOwner     bool
}

type CreatePartnerTenantInput struct {
	UserID      string
	Name        string
	Description string
	Address     string
	PhoneNumber string
	OpenHours   string
	LogoURL     string
	BannerURL   string
}
