package repository

import (
	"context"

	"github.com/dhegas/saas_gangsta/internal/domains/tenant/domain"
	"gorm.io/gorm"
)

type adminTenantRepository struct {
	db *gorm.DB
}

// Constructor untuk dependency injection
func NewAdminTenantRepository(db *gorm.DB) domain.AdminTenantRepository {
	return &adminTenantRepository{db: db}
}

func (r *adminTenantRepository) FindAll(ctx context.Context) ([]domain.TenantEntity, error) {
	var tenants []domain.TenantEntity

	// Kita paksa GORM untuk langsung membaca tabel "tenants"
	err := r.db.WithContext(ctx).Table("tenants").Find(&tenants).Error
	if err != nil {
		return nil, err
	}

	return tenants, nil
}

func (r *adminTenantRepository) UpdateStatus(ctx context.Context, tenantID string, status string) error {
	// Melakukan update hanya pada kolom status berdasarkan ID
	err := r.db.WithContext(ctx).
		Table("tenants").
		Where("id = ?", tenantID).
		Update("status", status).Error
	return err
}
