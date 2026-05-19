package domain

import (
	"context"

	"github.com/dhegas/saas_gangsta/internal/domains/tenant/dto"
)

type AdminTenantUsecase interface {
	CreateAdminTenant(ctx context.Context, req dto.CreateAdminTenantRequest) (*dto.CreateAdminTenantResponse, error)
	ListAllTenants(ctx context.Context, req dto.ListAllTenantsRequest) (*dto.ListAllTenantsResponse, error)
	SoftDeleteTenant(ctx context.Context, tenantID string) error
	GetTenantsByUserID(ctx context.Context, userID string) (*dto.GetTenantsByUserIDResponse, error)
	GetTenantByID(ctx context.Context, tenantID string) (*dto.AdminTenantResponse, error)
}

type AdminTenantRepository interface {
	CreateTenantForAdmin(ctx context.Context, input CreateAdminTenantInput) (*AdminTenant, error)
	ListAllTenants(ctx context.Context, limit, offset int) ([]AdminTenant, int64, error)
	SoftDeleteTenant(ctx context.Context, tenantID string) error
	GetTenantsByUserID(ctx context.Context, userID string) ([]AdminTenant, error)
	GetTenantByID(ctx context.Context, tenantID string) (*AdminTenant, error)
}

type AdminTenant struct {
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
	UserID      string
}

type CreateAdminTenantInput struct {
	UserID      string
	Name        string
	Status      string
	Description string
	Address     string
	PhoneNumber string
	OpenHours   string
	LogoURL     string
	BannerURL   string
}
