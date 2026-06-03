package repository

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/dhegas/saas_gangsta/internal/domains/tenant/domain"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound               = errors.New("user not found")
	ErrUserNotPartner             = errors.New("user is not partner")
)

type partnerTenantRepository struct {
	db *gorm.DB
}

type partnerUserRow struct {
	ID       string `gorm:"column:id"`
	Role     string `gorm:"column:role"`
	IsActive bool   `gorm:"column:is_active"`
}

func NewPartnerTenantRepository(db *gorm.DB) domain.PartnerTenantRepository {
	return &partnerTenantRepository{db: db}
}

func (r *partnerTenantRepository) FindPartnerByID(ctx context.Context, userID string) (*domain.PartnerUser, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	var row partnerUserRow
	err := r.db.WithContext(ctx).
		Table("users").
		Select("id::text AS id, role, is_active").
		Where("id = NULLIF(?, '')::uuid", userID).
		Take(&row).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("find partner by id: %w", err)
	}

	return &domain.PartnerUser{
		ID:       row.ID,
		Role:     decodeRole(row.Role),
		IsActive: row.IsActive,
	}, nil
}

func (r *partnerTenantRepository) CreateTenantForPartner(ctx context.Context, input domain.CreatePartnerTenantInput) (*domain.PartnerTenant, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	tenantSlug := generateTenantSlug(input.Name)

	createdTenant := &domain.PartnerTenant{}
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
			 VALUES (?, ?, 'active', NULLIF(?, '')::uuid, ?, ?, ?, ?, ?, ?, NOW(), NOW())
			 RETURNING id::text, name, slug, status, description, address, phone_number, open_hours, logo_url, banner_url`,
			input.Name,
			tenantSlug,
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
		createdTenant.IsOwner = true

		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	return createdTenant, nil
}

func (r *partnerTenantRepository) ListTenantsByPartner(ctx context.Context, userID string) ([]domain.PartnerTenant, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	tenants := make([]domain.PartnerTenant, 0)
	err := r.db.WithContext(ctx).Raw(
		`SELECT t.id::text AS id, t.name, t.slug, t.status, t.description, t.address, t.phone_number, t.open_hours, t.logo_url, t.banner_url, TRUE AS is_owner
		 FROM tenants t
		 WHERE t.user_id = NULLIF(?, '')::uuid AND t.deleted_at IS NULL
		 ORDER BY t.created_at DESC`,
		userID,
	).Scan(&tenants).Error
	if err != nil {
		return nil, fmt.Errorf("list partner tenants: %w", err)
	}

	return tenants, nil
}

func (r *partnerTenantRepository) SoftDeleteTenant(ctx context.Context, userID string, tenantID string) error {
	if r.db == nil {
		return fmt.Errorf("database is not initialized")
	}

	result := r.db.WithContext(ctx).Exec(
		`UPDATE tenants 
		 SET deleted_at = NOW(), updated_at = NOW() 
		 WHERE id = NULLIF(?, '')::uuid 
		   AND user_id = NULLIF(?, '')::uuid 
		   AND deleted_at IS NULL`,
		tenantID,
		userID,
	)

	if result.Error != nil {
		return fmt.Errorf("soft delete tenant: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("tenant not found or already deleted")
	}

	return nil
}

func decodeRole(role string) string {
	return strings.ToUpper(strings.TrimSpace(role))
}

func generateTenantSlug(name string) string {
	normalized := strings.ToLower(strings.TrimSpace(name))
	re := regexp.MustCompile(`[^a-z0-9]+`)
	slug := strings.Trim(re.ReplaceAllString(normalized, "-"), "-")
	if slug == "" {
		slug = "tenant"
	}

	return fmt.Sprintf("%s-%d", slug, time.Now().UTC().UnixNano())
}

func (r *partnerTenantRepository) GetTenantByIDAndPartner(ctx context.Context, userID string, tenantID string) (*domain.PartnerTenant, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	var tenant domain.PartnerTenant
	err := r.db.WithContext(ctx).Raw(
		`SELECT t.id::text AS id, t.name, t.slug, t.status, t.description, t.address, t.phone_number, t.open_hours, t.logo_url, t.banner_url, TRUE AS is_owner
		 FROM tenants t
		 WHERE t.id = NULLIF(?, '')::uuid 
		   AND t.user_id = NULLIF(?, '')::uuid 
		   AND t.deleted_at IS NULL`,
		tenantID,
		userID,
	).Scan(&tenant).Error
	if err != nil {
		return nil, fmt.Errorf("get tenant by id and partner: %w", err)
	}

	if tenant.ID == "" {
		return nil, ErrTenantNotFound
	}

	return &tenant, nil
}

func (r *partnerTenantRepository) UpdateTenant(ctx context.Context, userID string, tenantID string, name string, description string, address string, phoneNumber string) (*domain.PartnerTenant, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	var tenant domain.PartnerTenant
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Exec(
			`UPDATE tenants 
			 SET name = ?, description = ?, address = ?, phone_number = ?, updated_at = NOW() 
			 WHERE id = NULLIF(?, '')::uuid 
			   AND user_id = NULLIF(?, '')::uuid 
			   AND deleted_at IS NULL`,
			name,
			description,
			address,
			phoneNumber,
			tenantID,
			userID,
		)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrTenantNotFound
		}

		return tx.Raw(
			`SELECT t.id::text AS id, t.name, t.slug, t.status, t.description, t.address, t.phone_number, t.open_hours, t.logo_url, t.banner_url, TRUE AS is_owner
			 FROM tenants t
			 WHERE t.id = NULLIF(?, '')::uuid 
			   AND t.user_id = NULLIF(?, '')::uuid 
			   AND t.deleted_at IS NULL`,
			tenantID,
			userID,
		).Scan(&tenant).Error
	})

	if err != nil {
		return nil, err
	}

	return &tenant, nil
}

