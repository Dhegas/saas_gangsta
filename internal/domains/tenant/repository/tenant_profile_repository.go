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
		RETURNING id::text, tenant_id::text, name, COALESCE(description, '') AS description, sort_order, is_active, created_at, updated_at`,
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

func (r *tenantProfileRepository) FindByID(ctx context.Context, tenantID, profileID string) (*domain.TenantProfileEntity, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	var profile domain.TenantProfileEntity
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
		 WHERE id = NULLIF(?, '')::uuid
		   AND tenant_id = NULLIF(?, '')::uuid
		   AND deleted_at IS NULL
		 LIMIT 1`,
		profileID,
		tenantID,
	).Scan(&profile).Error
	if err != nil {
		return nil, fmt.Errorf("find tenant profile by id: %w", err)
	}

	if profile.ID == "" {
		return nil, nil
	}

	return &profile, nil
}

func (r *tenantProfileRepository) Update(ctx context.Context, profile *domain.TenantProfileEntity) error {
	if r.db == nil {
		return fmt.Errorf("database is not initialized")
	}

	err := r.db.WithContext(ctx).Raw(
		`UPDATE tenant_profiles
		 SET
			name        = ?,
			description = NULLIF(?, ''),
			sort_order  = ?,
			is_active   = ?,
			updated_at  = NOW()
		 WHERE id = NULLIF(?, '')::uuid
		   AND tenant_id = NULLIF(?, '')::uuid
		   AND deleted_at IS NULL
		 RETURNING id::text, tenant_id::text, name, COALESCE(description, '') AS description, sort_order, is_active, created_at, updated_at`,
		profile.Name,
		profile.Description,
		profile.SortOrder,
		profile.IsActive,
		profile.ID,
		profile.TenantID,
	).Scan(profile).Error
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return fmt.Errorf("duplicate tenant profile name: %w", err)
		}
		return fmt.Errorf("update tenant profile: %w", err)
	}

	return nil
}

func (r *tenantProfileRepository) SoftDelete(ctx context.Context, tenantID, profileID string) error {
	if r.db == nil {
		return fmt.Errorf("database is not initialized")
	}

	result := r.db.WithContext(ctx).Exec(
		`UPDATE tenant_profiles
		 SET deleted_at = NOW(), updated_at = NOW()
		 WHERE id = NULLIF(?, '')::uuid
		   AND tenant_id = NULLIF(?, '')::uuid
		   AND deleted_at IS NULL`,
		profileID,
		tenantID,
	)
	if result.Error != nil {
		return fmt.Errorf("soft delete tenant profile: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("not found")
	}

	return nil
}
