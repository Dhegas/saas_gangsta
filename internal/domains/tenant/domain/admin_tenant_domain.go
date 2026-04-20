package domain

import (
	"context"
	"time"

	"github.com/dhegas/saas_gangsta/internal/domains/tenant/dto"
)

// AdminTenantUsecase adalah kontrak untuk logika bisnis Admin
type AdminTenantUsecase interface {
	GetAllTenants(ctx context.Context) ([]dto.TenantResponse, error)
	GetTenantByID(ctx context.Context, tenantID string) (*dto.TenantResponse, error)
	CreateTenant(ctx context.Context, req dto.CreateTenantRequest) (*dto.TenantResponse, error)
	UpdateTenant(ctx context.Context, tenantID string, req dto.UpdateTenantRequest) (*dto.TenantResponse, error)
	SoftDeleteTenant(ctx context.Context, tenantID string) error
	UpdateTenantStatus(ctx context.Context, tenantID string, status string) error
}

// AdminTenantRepository adalah kontrak untuk query ke database
type AdminTenantRepository interface {
	FindAll(ctx context.Context) ([]TenantEntity, error)
	FindByID(ctx context.Context, tenantID string) (*TenantEntity, error)
	Create(ctx context.Context, entity *TenantEntity) error
	Update(ctx context.Context, entity *TenantEntity) error
	SoftDelete(ctx context.Context, tenantID string) error
	UpdateStatus(ctx context.Context, tenantID string, status string) error
}

// TenantEntity merepresentasikan tabel tenants di database (dipakai GORM)
type TenantEntity struct {
	ID        string
	Name      string
	Slug      string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
