package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/dhegas/saas_gangsta/internal/domains/tenant/domain"
	"gorm.io/gorm"
)

type publicTenantRepository struct {
	db *gorm.DB
}

func NewPublicTenantRepository(db *gorm.DB) domain.PublicTenantRepository {
	return &publicTenantRepository{db: db}
}

func (r *publicTenantRepository) ListPublicTenants(ctx context.Context, search string, limit, offset int) ([]domain.PublicTenant, int64, error) {
	if r.db == nil {
		return nil, 0, fmt.Errorf("database is not initialized")
	}

	query := r.db.WithContext(ctx).Table("tenants").
		Where("status = 'active'").
		Where("is_public = true").
		Where("deleted_at IS NULL")

	if strings.TrimSpace(search) != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ? OR LOWER(address) LIKE ?", searchPattern, searchPattern, searchPattern)
	}

	var totalItems int64
	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, fmt.Errorf("count public tenants: %w", err)
	}

	var tenants []domain.PublicTenant
	err := query.Select("id::text AS id, name, slug, status, description, address, phone_number, open_hours, logo_url, banner_url, is_public").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(&tenants).Error
	if err != nil {
		return nil, 0, fmt.Errorf("list public tenants: %w", err)
	}

	return tenants, totalItems, nil
}

func (r *publicTenantRepository) FindTenantBySlug(ctx context.Context, slug string) (*domain.PublicTenant, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	var tenant domain.PublicTenant
	err := r.db.WithContext(ctx).Table("tenants").
		Select("id::text AS id, name, slug, status, description, address, phone_number, open_hours, logo_url, banner_url, is_public").
		Where("slug = ?", slug).
		Where("status = 'active'").
		Where("is_public = true").
		Where("deleted_at IS NULL").
		Scan(&tenant).Error
	if err != nil {
		return nil, fmt.Errorf("find tenant by slug: %w", err)
	}

	if tenant.ID == "" {
		return nil, gorm.ErrRecordNotFound
	}

	return &tenant, nil
}
