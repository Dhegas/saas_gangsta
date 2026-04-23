package domain

import (
	"context"
	"time"

	"github.com/dhegas/saas_gangsta/internal/domains/tenant/dto"
)

// TenantProfileUsecase adalah kontrak logika bisnis tenant profile.
type TenantProfileUsecase interface {
	CreateProfile(ctx context.Context, tenantID string, req dto.CreateTenantProfileRequest) (*dto.TenantProfileResponse, error)
	ListProfiles(ctx context.Context, tenantID string) ([]dto.TenantProfileResponse, error)
	GetProfileByID(ctx context.Context, tenantID, profileID string) (*dto.TenantProfileResponse, error)
	UpdateProfile(ctx context.Context, tenantID, profileID string, req dto.UpdateTenantProfileRequest) (*dto.TenantProfileResponse, error)
	DeleteProfile(ctx context.Context, tenantID, profileID string) error
	ToggleActive(ctx context.Context, tenantID, profileID string) (*dto.TenantProfileResponse, error)
}

// TenantProfileRepository adalah kontrak akses data tenant profile.
type TenantProfileRepository interface {
	Create(ctx context.Context, profile *TenantProfileEntity) error
	ListByTenantID(ctx context.Context, tenantID string) ([]TenantProfileEntity, error)
	FindByID(ctx context.Context, tenantID, profileID string) (*TenantProfileEntity, error)
	Update(ctx context.Context, profile *TenantProfileEntity) error
	SoftDelete(ctx context.Context, tenantID, profileID string) error
}

// TenantProfileEntity merepresentasikan tabel tenant_profiles.
type TenantProfileEntity struct {
	ID          string
	TenantID    string
	Name        string
	Description string
	SortOrder   int
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

