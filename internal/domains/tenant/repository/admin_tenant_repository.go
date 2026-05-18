package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dhegas/saas_gangsta/internal/domains/tenant/domain"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var ErrTenantNotFound = errors.New("tenant not found or already deleted")

type adminTenantRepository struct {
	db *gorm.DB
}

func NewAdminTenantRepository(db *gorm.DB) domain.AdminTenantRepository {
	return &adminTenantRepository{db: db}
}

func (r *adminTenantRepository) CreateTenantForAdmin(ctx context.Context, input domain.CreateAdminTenantInput) (*domain.AdminTenant, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	tenantSlug := generateTenantSlug(input.Name)
	status := "active"
	if strings.TrimSpace(input.Status) != "" {
		status = input.Status
	}

	createdTenant := &domain.AdminTenant{}
	txErr := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var userRow struct {
			ID       string `gorm:"column:id"`
			Role     string `gorm:"column:role"`
			IsActive bool   `gorm:"column:is_active"`
		}

		if err := tx.Raw(
			`SELECT id::text, role, is_active
			 FROM users
			 WHERE id = NULLIF(?, '')::uuid
			 FOR UPDATE`,
			input.UserID,
		).Scan(&userRow).Error; err != nil {
			return fmt.Errorf("lock partner user: %w", err)
		}
		if userRow.ID == "" || !userRow.IsActive {
			return ErrUserNotFound
		}
		if decodeRole(userRow.Role) != "PARTNER" {
			return ErrUserNotPartner
		}

		if err := tx.Raw(
			`INSERT INTO tenants (name, slug, status, user_id, description, address, phone_number, open_hours, logo_url, banner_url, created_at, updated_at)
			 VALUES (?, ?, ?, NULLIF(?, '')::uuid, ?, ?, ?, ?, ?, ?, NOW(), NOW())
			 RETURNING id::text, name, slug, status, description, address, phone_number, open_hours, logo_url, banner_url, user_id::text AS user_id`,
			input.Name,
			tenantSlug,
			status,
			input.UserID,
			input.Description,
			input.Address,
			input.PhoneNumber,
			input.OpenHours,
			input.LogoURL,
			input.BannerURL,
		).Scan(createdTenant).Error; err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return fmt.Errorf("duplicate tenant slug: %w", err)
			}
			return fmt.Errorf("create tenant: %w", err)
		}

		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	return createdTenant, nil
}

func (r *adminTenantRepository) ListAllTenants(ctx context.Context, limit, offset int) ([]domain.AdminTenant, int64, error) {
	if r.db == nil {
		return nil, 0, fmt.Errorf("database is not initialized")
	}

	var totalItems int64
	if err := r.db.WithContext(ctx).Raw(
		`SELECT COUNT(*) FROM tenants WHERE deleted_at IS NULL`,
	).Scan(&totalItems).Error; err != nil {
		return nil, 0, fmt.Errorf("count all tenants: %w", err)
	}

	tenants := make([]domain.AdminTenant, 0)
	err := r.db.WithContext(ctx).Raw(
		`SELECT t.id::text AS id, t.name, t.slug, t.status, t.description, t.address, t.phone_number, t.open_hours, t.logo_url, t.banner_url, t.user_id::text AS user_id
		 FROM tenants t
		 WHERE t.deleted_at IS NULL
		 ORDER BY t.created_at DESC
		 LIMIT ? OFFSET ?`,
		limit,
		offset,
	).Scan(&tenants).Error
	if err != nil {
		return nil, 0, fmt.Errorf("list all tenants: %w", err)
	}

	return tenants, totalItems, nil
}

func (r *adminTenantRepository) SoftDeleteTenant(ctx context.Context, tenantID string) error {
	if r.db == nil {
		return fmt.Errorf("database is not initialized")
	}

	res := r.db.WithContext(ctx).Exec(
		`UPDATE tenants 
		 SET deleted_at = NOW(), updated_at = NOW() 
		 WHERE id = NULLIF(?, '')::uuid AND deleted_at IS NULL`,
		tenantID,
	)
	if res.Error != nil {
		return fmt.Errorf("soft delete tenant: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return ErrTenantNotFound
	}

	return nil
}

func (r *adminTenantRepository) GetTenantsByUserID(ctx context.Context, userID string) ([]domain.AdminTenant, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	tenants := make([]domain.AdminTenant, 0)
	err := r.db.WithContext(ctx).Raw(
		`SELECT t.id::text AS id, t.name, t.slug, t.status, t.description, t.address, t.phone_number, t.open_hours, t.logo_url, t.banner_url, t.user_id::text AS user_id
		 FROM tenants t
		 WHERE t.user_id = NULLIF(?, '')::uuid AND t.deleted_at IS NULL
		 ORDER BY t.created_at DESC`,
		userID,
	).Scan(&tenants).Error
	if err != nil {
		return nil, fmt.Errorf("get tenants by user id: %w", err)
	}

	return tenants, nil
}
