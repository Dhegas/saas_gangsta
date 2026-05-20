package domain

import (
	"context"

	"github.com/dhegas/saas_gangsta/internal/domains/tenant/dto"
)

type PublicTenantUsecase interface {
	ListPublicTenants(ctx context.Context, req dto.ListPublicTenantsRequest) (*dto.ListPublicTenantsResponse, error)
	GetPublicTenantBySlug(ctx context.Context, slug string) (*dto.PublicTenantDetailResponse, error)
}

type PublicTenantRepository interface {
	ListPublicTenants(ctx context.Context, search string, limit, offset int) ([]PublicTenant, int64, error)
	FindTenantBySlug(ctx context.Context, slug string) (*PublicTenant, error)
}

type PublicTenant struct {
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
	IsPublic    bool
}
