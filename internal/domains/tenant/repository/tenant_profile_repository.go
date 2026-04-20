package repository

import (
	"context"
	"fmt"

	"github.com/dhegas/saas_gangsta/internal/domains/tenant/domain"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type tenantProfileRepository struct {
	db *gorm.DB
}

// NewTenantProfileRepository membuat repository tenant profile.
func NewTenantProfileRepository(db *gorm.DB) domain.TenantProfileRepository {
	return &tenantProfileRepository{db: db}
}

func (r *tenantProfileRepository) Create(ctx context.Context, profile *domain.TenantProfileEntity) error {
	if r.db == nil {
		return fmt.Errorf("database is not initialized")
	}

	err := r.db.WithContext(ctx).Raw(
		`INSERT INTO tenant_profiles (
			tenant_id,
			name,
			description,
			sort_order,
			is_active,
			created_at,
			updated_at
		)
		VALUES (
			NULLIF(?, '')::uuid,
			?,
			NULLIF(?, ''),
			?,
			?,
			NOW(),
			NOW()
		)
		RETURNING id::text, tenant_id::text, name, COALESCE(description, ''), sort_order, is_active, created_at, updated_at`,
		profile.TenantID,
		profile.Name,
		profile.Description,
		profile.SortOrder,
		profile.IsActive,
	).Scan(profile).Error
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return fmt.Errorf("duplicate tenant profile name: %w", err)
		}
		return fmt.Errorf("create tenant profile: %w", err)
	}

	return nil
}

func (r *tenantProfileRepository) ListByTenantID(ctx context.Context, tenantID string) ([]domain.TenantProfileEntity, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	profiles := make([]domain.TenantProfileEntity, 0)
	err := r.db.WithContext(ctx).Raw(
		`SELECT
			id::text,
			tenant_id::text,
			name,
			COALESCE(description, '') AS description,
			sort_order,
			is_active,
			created_at,
			updated_at
		 FROM tenant_profiles
		 WHERE tenant_id = NULLIF(?, '')::uuid
		   AND deleted_at IS NULL
		 ORDER BY sort_order ASC, created_at DESC`,
		tenantID,
	).Scan(&profiles).Error
	if err != nil {
		return nil, fmt.Errorf("list tenant profiles: %w", err)
	}

	return profiles, nil
}
