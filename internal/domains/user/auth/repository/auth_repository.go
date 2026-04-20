package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/dhegas/saas_gangsta/internal/domains/user/auth/domain"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound                = errors.New("user not found")
	ErrUserNotCustomer             = errors.New("user is not customer")
	ErrUserNotMerchant             = errors.New("user is not merchant")
	ErrSubscriptionPlanNotFound    = errors.New("subscription plan not found")
	ErrSubscriptionAlreadyExists   = errors.New("subscription already exists")
	ErrMerchantSubscriptionMissing = errors.New("merchant subscription is required")
	ErrTenantLimitReached          = errors.New("tenant limit reached")
)

type AuthRepository interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
	SubscribeAndUpgradeCustomer(ctx context.Context, input SubscribeUpgradeInput) (*domain.User, error)
	CreateTenantForMerchant(ctx context.Context, input CreateMerchantTenantInput) (*domain.MerchantTenant, error)
	ListTenantsByMerchant(ctx context.Context, userID string) ([]domain.MerchantTenant, error)
}

type authRepository struct {
	db *gorm.DB
}

type authUserRow struct {
	ID           string `gorm:"column:id"`
	TenantID     string `gorm:"column:tenant_id"`
	Email        string `gorm:"column:email"`
	PasswordHash string `gorm:"column:password_hash"`
	Role         string `gorm:"column:role"`
	IsActive     bool   `gorm:"column:is_active"`
	TenantStatus string `gorm:"column:tenant_status"`
}

type SubscribeUpgradeInput struct {
	UserID string
	PlanID string
}

type CreateMerchantTenantInput struct {
	UserID string
	Name   string
}

func decodeRole(role string) string {
	trimmed := strings.TrimSpace(role)
	if strings.EqualFold(trimmed, "basic") || strings.EqualFold(trimmed, "customer") {
		return "BASIC"
	}
	if strings.EqualFold(trimmed, "mitra") || strings.EqualFold(trimmed, "merchant") {
		return "MITRA"
	}
	if strings.EqualFold(trimmed, "admin") {
		return "ADMIN"
	}

	return strings.ToUpper(trimmed)
}

func encodeRole(role string) string {
	return decodeRole(role)
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	var row authUserRow
	err := r.db.WithContext(ctx).
		Table("users u").
		Select(
			"u.id, COALESCE(u.tenant_id::text, '') AS tenant_id, u.email, u.password_hash, u.role, u.is_active, COALESCE(t.status, 'active') AS tenant_status",
		).
		Joins("LEFT JOIN tenants t ON t.id = u.tenant_id").
		Where("LOWER(u.email) = LOWER(?)", email).
		Take(&row).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("find user by email: %w", err)
	}

	return &domain.User{
		ID:           row.ID,
		TenantID:     row.TenantID,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		Role:         decodeRole(row.Role),
		IsActive:     row.IsActive,
		TenantStatus: row.TenantStatus,
	}, nil
}

func (r *authRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	var row authUserRow
	err := r.db.WithContext(ctx).
		Table("users u").
		Select(
			"u.id, COALESCE(u.tenant_id::text, '') AS tenant_id, u.email, u.password_hash, u.role, u.is_active, COALESCE(t.status, 'active') AS tenant_status",
		).
		Joins("LEFT JOIN tenants t ON t.id = u.tenant_id").
		Where("u.id = ?", id).
		Take(&row).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}

	return &domain.User{
		ID:           row.ID,
		TenantID:     row.TenantID,
		Email:        row.Email,
		FullName:     "",
		PasswordHash: row.PasswordHash,
		Role:         decodeRole(row.Role),
		IsActive:     row.IsActive,
		TenantStatus: row.TenantStatus,
	}, nil
}

func (r *authRepository) CreateUser(ctx context.Context, user *domain.User) error {
	if r.db == nil {
		return fmt.Errorf("database is not initialized")
	}

	row := struct {
		ID string `gorm:"column:id"`
	}{ID: ""}

	err := r.db.WithContext(ctx).Raw(
		`INSERT INTO users (tenant_id, email, password_hash, full_name, role, is_active, created_at, updated_at)
		 VALUES (NULLIF(?, '')::uuid, ?, ?, ?, ?, TRUE, NOW(), NOW())
		 RETURNING id::text`,
		user.TenantID,
		user.Email,
		user.PasswordHash,
		user.FullName,
		encodeRole(user.Role),
	).Scan(&row).Error
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("duplicate email: %w", err)
		}
		return fmt.Errorf("create user: %w", err)
	}

	user.ID = row.ID
	user.IsActive = true

	return nil
}

func (r *authRepository) SubscribeAndUpgradeCustomer(ctx context.Context, input SubscribeUpgradeInput) (*domain.User, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

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
			return fmt.Errorf("lock user: %w", err)
		}
		if userRow.ID == "" || !userRow.IsActive {
			return ErrUserNotFound
		}
		if decodeRole(userRow.Role) != "BASIC" {
			return ErrUserNotCustomer
		}

		var billingCycle string
		if err := tx.Raw(
			`SELECT billing_cycle
			 FROM subscription_plans
			 WHERE id = NULLIF(?, '')::uuid AND is_active = TRUE AND deleted_at IS NULL`,
			input.PlanID,
		).Scan(&billingCycle).Error; err != nil {
			return fmt.Errorf("find subscription plan: %w", err)
		}
		if strings.TrimSpace(billingCycle) == "" {
			return ErrSubscriptionPlanNotFound
		}

		now := time.Now().UTC()
		expiresAt := now.AddDate(0, 1, 0)
		if strings.TrimSpace(billingCycle) == "yearly" {
			expiresAt = now.AddDate(1, 0, 0)
		}

		if err := tx.Exec(
			`INSERT INTO subscriptions (
				tenant_id,
				subscriber_user_id,
				plan_id,
				status,
				started_at,
				expires_at,
				created_at,
				updated_at
			)
			VALUES (
				NULL,
				NULLIF(?, '')::uuid,
				NULLIF(?, '')::uuid,
				'pending_tenant',
				?,
				?,
				NOW(),
				NOW()
			)`,
			input.UserID,
			input.PlanID,
			now,
			expiresAt,
		).Error; err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return ErrSubscriptionAlreadyExists
			}
			return fmt.Errorf("create pending subscription: %w", err)
		}

		if err := tx.Exec(
			`UPDATE users
			 SET role = ?, updated_at = NOW()
			 WHERE id = NULLIF(?, '')::uuid`,
			encodeRole("MITRA"),
			input.UserID,
		).Error; err != nil {
			return fmt.Errorf("upgrade user to merchant: %w", err)
		}

		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	upgradedUser, err := r.FindByID(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("reload upgraded user: %w", err)
	}
	if upgradedUser == nil {
		return nil, ErrUserNotFound
	}

	return upgradedUser, nil
}

func (r *authRepository) CreateTenantForMerchant(ctx context.Context, input CreateMerchantTenantInput) (*domain.MerchantTenant, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	tenantSlug := generateTenantSlug(input.Name)

	createdTenant := &domain.MerchantTenant{}
	txErr := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var userRow struct {
			ID       string `gorm:"column:id"`
			Role     string `gorm:"column:role"`
			IsActive bool   `gorm:"column:is_active"`
			TenantID string `gorm:"column:tenant_id"`
		}

		if err := tx.Raw(
			`SELECT id::text, role, is_active, COALESCE(tenant_id::text, '') AS tenant_id
			 FROM users
			 WHERE id = NULLIF(?, '')::uuid
			 FOR UPDATE`,
			input.UserID,
		).Scan(&userRow).Error; err != nil {
			return fmt.Errorf("lock merchant user: %w", err)
		}
		if userRow.ID == "" || !userRow.IsActive {
			return ErrUserNotFound
		}
		if decodeRole(userRow.Role) != "MITRA" {
			return ErrUserNotMerchant
		}

		var subscriptionRow struct {
			ID         string        `gorm:"column:id"`
			MaxTenants sql.NullInt64 `gorm:"column:max_tenants"`
		}

		if err := tx.Raw(
			`SELECT s.id::text AS id, sp.max_tenants
			 FROM subscriptions s
			 INNER JOIN subscription_plans sp ON sp.id = s.plan_id
			 WHERE s.subscriber_user_id = NULLIF(?, '')::uuid
			   AND s.status IN ('active', 'pending_tenant', 'trial')
			   AND s.deleted_at IS NULL
			   AND sp.deleted_at IS NULL
			 ORDER BY
			   CASE s.status
			     WHEN 'active' THEN 0
			     WHEN 'pending_tenant' THEN 1
			     ELSE 2
			   END,
			   s.created_at DESC
			 LIMIT 1`,
			input.UserID,
		).Scan(&subscriptionRow).Error; err != nil {
			return fmt.Errorf("find merchant subscription: %w", err)
		}
		if subscriptionRow.ID == "" {
			return ErrMerchantSubscriptionMissing
		}

		var tenantCount int64
		if err := tx.Raw(
			`SELECT COUNT(*)
			 FROM merchant_tenants mt
			 INNER JOIN tenants t ON t.id = mt.tenant_id
			 WHERE mt.user_id = NULLIF(?, '')::uuid AND t.deleted_at IS NULL`,
			input.UserID,
		).Scan(&tenantCount).Error; err != nil {
			return fmt.Errorf("count merchant tenants: %w", err)
		}

		if subscriptionRow.MaxTenants.Valid && tenantCount >= subscriptionRow.MaxTenants.Int64 {
			return ErrTenantLimitReached
		}

		if err := tx.Raw(
			`INSERT INTO tenants (name, slug, status, created_at, updated_at)
			 VALUES (?, ?, 'active', NOW(), NOW())
			 RETURNING id::text, name, slug, status`,
			input.Name,
			tenantSlug,
		).Scan(createdTenant).Error; err != nil {
			return fmt.Errorf("create tenant: %w", err)
		}
		createdTenant.IsOwner = true

		if err := tx.Exec(
			`INSERT INTO merchant_tenants (user_id, tenant_id, is_owner, created_at, updated_at)
			 VALUES (NULLIF(?, '')::uuid, NULLIF(?, '')::uuid, TRUE, NOW(), NOW())`,
			input.UserID,
			createdTenant.ID,
		).Error; err != nil {
			return fmt.Errorf("link merchant tenant: %w", err)
		}

		if userRow.TenantID == "" {
			if err := tx.Exec(
				`UPDATE users
				 SET tenant_id = NULLIF(?, '')::uuid, updated_at = NOW()
				 WHERE id = NULLIF(?, '')::uuid`,
				createdTenant.ID,
				input.UserID,
			).Error; err != nil {
				return fmt.Errorf("set default tenant on user: %w", err)
			}
		}

		if err := tx.Exec(
			`UPDATE subscriptions
			 SET tenant_id = NULLIF(?, '')::uuid, status = 'active', updated_at = NOW()
			 WHERE id = NULLIF(?, '')::uuid AND status = 'pending_tenant' AND tenant_id IS NULL`,
			createdTenant.ID,
			subscriptionRow.ID,
		).Error; err != nil {
			return fmt.Errorf("activate pending subscription: %w", err)
		}

		return nil
	})
	if txErr != nil {
		return nil, txErr
	}

	return createdTenant, nil
}

func (r *authRepository) ListTenantsByMerchant(ctx context.Context, userID string) ([]domain.MerchantTenant, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	tenants := make([]domain.MerchantTenant, 0)
	err := r.db.WithContext(ctx).Raw(
		`SELECT t.id::text AS id, t.name, t.slug, t.status, mt.is_owner
		 FROM merchant_tenants mt
		 INNER JOIN tenants t ON t.id = mt.tenant_id
		 WHERE mt.user_id = NULLIF(?, '')::uuid AND t.deleted_at IS NULL
		 ORDER BY mt.created_at DESC`,
		userID,
	).Scan(&tenants).Error
	if err != nil {
		return nil, fmt.Errorf("list merchant tenants: %w", err)
	}

	return tenants, nil
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
