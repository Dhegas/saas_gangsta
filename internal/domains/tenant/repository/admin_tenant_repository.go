package repository

import (
	"context"
	"time"

	"github.com/dhegas/saas_gangsta/internal/domains/tenant/domain"
	"gorm.io/gorm"
)

type adminTenantRepository struct {
	db *gorm.DB
}

// NewAdminTenantRepository adalah constructor untuk dependency injection
func NewAdminTenantRepository(db *gorm.DB) domain.AdminTenantRepository {
	return &adminTenantRepository{db: db}
}

// FindAll mengambil semua tenant yang belum dihapus (soft delete)
func (r *adminTenantRepository) FindAll(ctx context.Context) ([]domain.TenantEntity, error) {
	var tenants []domain.TenantEntity

	err := r.db.WithContext(ctx).
		Table("tenants").
		Where("deleted_at IS NULL").
		Find(&tenants).Error
	if err != nil {
		return nil, err
	}

	return tenants, nil
}

// FindByID mengambil satu tenant berdasarkan ID, hanya yang belum dihapus
func (r *adminTenantRepository) FindByID(ctx context.Context, tenantID string) (*domain.TenantEntity, error) {
	var tenant domain.TenantEntity

	err := r.db.WithContext(ctx).
		Table("tenants").
		Where("id = ? AND deleted_at IS NULL", tenantID).
		First(&tenant).Error
	if err != nil {
		return nil, err
	}

	return &tenant, nil
}

// Create menyimpan tenant baru ke database
func (r *adminTenantRepository) Create(ctx context.Context, entity *domain.TenantEntity) error {
	return r.db.WithContext(ctx).
		Table("tenants").
		Create(entity).Error
}

// Update memperbarui data tenant (name, slug, status) berdasarkan ID
func (r *adminTenantRepository) Update(ctx context.Context, entity *domain.TenantEntity) error {
	return r.db.WithContext(ctx).
		Table("tenants").
		Where("id = ? AND deleted_at IS NULL", entity.ID).
		Updates(map[string]interface{}{
			"name":       entity.Name,
			"slug":       entity.Slug,
			"status":     entity.Status,
			"updated_at": time.Now(),
		}).Error
}

// SoftDelete mengisi kolom deleted_at sehingga tenant dianggap terhapus
func (r *adminTenantRepository) SoftDelete(ctx context.Context, tenantID string) error {
	return r.db.WithContext(ctx).
		Table("tenants").
		Where("id = ? AND deleted_at IS NULL", tenantID).
		Update("deleted_at", time.Now()).Error
}

// UpdateStatus memperbarui hanya kolom status berdasarkan ID
func (r *adminTenantRepository) UpdateStatus(ctx context.Context, tenantID string, status string) error {
	return r.db.WithContext(ctx).
		Table("tenants").
		Where("id = ? AND deleted_at IS NULL", tenantID).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}
