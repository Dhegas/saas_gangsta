package domain

import (
	"context"
	"time"

	"github.com/dhegas/saas_gangsta/internal/modules/adminTenant/dto"
)

// AdminTenantUsecase adalah kontrak untuk logika bisnis Admin
type AdminTenantUsecase interface {
	GetAllTenants(ctx context.Context) ([]dto.TenantResponse, error)
	UpdateTenantStatus(ctx context.Context, tenantID string, status string) error
}

// AdminTenantRepository adalah kontrak untuk query ke database
type AdminTenantRepository interface {
	FindAll(ctx context.Context) ([]TenantEntity, error)
	UpdateStatus(ctx context.Context, tenantID string, status string) error
}

// TenantEntity merepresentasikan tabel di database (akan dipakai oleh GORM)
type TenantEntity struct {
	ID        string
	Name      string
	Slug      string
	Status    string
	CreatedAt time.Time
}
