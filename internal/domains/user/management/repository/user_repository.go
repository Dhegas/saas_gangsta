package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dhegas/saas_gangsta/internal/domains/user/management/domain"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrNoFieldsToUpdate  = errors.New("no fields to update")
	ErrEmailAlreadyExist = errors.New("email already exists")
)

type UpdateUserInput struct {
	Email    *string
	FullName *string
	Role     *string
}

type UserRepository interface {
	ListByTenant(ctx context.Context, tenantID string) ([]domain.User, error)
	FindByIDAndTenant(ctx context.Context, tenantID, userID string) (*domain.User, error)
	UpdateByIDAndTenant(ctx context.Context, tenantID, userID string, input UpdateUserInput) (*domain.User, error)
	SoftDeleteByIDAndTenant(ctx context.Context, tenantID, userID string) error
	ToggleActiveByIDAndTenant(ctx context.Context, tenantID, userID string) (*domain.User, error)
	ListAllUsersForAdmin(ctx context.Context, roleFilter string, limit, offset int) ([]domain.User, int64, error)
	FindByIDForAdmin(ctx context.Context, userID string) (*domain.User, error)
}

type userRepository struct {
	db *gorm.DB
}

type userRow struct {
	ID        string `gorm:"column:id"`
	TenantID  string `gorm:"column:tenant_id"`
	Email     string `gorm:"column:email"`
	FullName  string `gorm:"column:full_name"`
	Role      string `gorm:"column:role"`
	IsActive  bool   `gorm:"column:is_active"`
	CreatedAt string `gorm:"column:created_at"`
	UpdatedAt string `gorm:"column:updated_at"`
}

func decodeRole(role string) string {
	return strings.ToUpper(strings.TrimSpace(role))
}

func encodeRole(role string) string {
	return strings.ToUpper(strings.TrimSpace(role))
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) ListByTenant(ctx context.Context, tenantID string) ([]domain.User, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	var rows []struct {
		ID       string `gorm:"column:id"`
		TenantID string `gorm:"column:tenant_id"`
		Email    string `gorm:"column:email"`
		FullName string `gorm:"column:full_name"`
		Role     string `gorm:"column:role"`
		IsActive bool   `gorm:"column:is_active"`
	}

	err := r.db.WithContext(ctx).Raw(
		`SELECT u.id::text,
		        t.id::text AS tenant_id,
		        u.email,
		        u.full_name,
		        u.role,
		        u.is_active
		 FROM tenants t
		 INNER JOIN users u ON u.id = t.user_id
		 WHERE t.id = NULLIF(?, '')::uuid
		   AND t.deleted_at IS NULL
		   AND u.deleted_at IS NULL
		 ORDER BY u.created_at DESC`,
		tenantID,
	).Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("list users by tenant: %w", err)
	}

	users := make([]domain.User, 0, len(rows))
	for _, row := range rows {
		users = append(users, domain.User{
			ID:       row.ID,
			TenantID: row.TenantID,
			Email:    row.Email,
			FullName: row.FullName,
			Role:     decodeRole(row.Role),
			IsActive: row.IsActive,
		})
	}

	return users, nil
}

func (r *userRepository) FindByIDAndTenant(ctx context.Context, tenantID, userID string) (*domain.User, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	var row struct {
		ID       string `gorm:"column:id"`
		TenantID string `gorm:"column:tenant_id"`
		Email    string `gorm:"column:email"`
		FullName string `gorm:"column:full_name"`
		Role     string `gorm:"column:role"`
		IsActive bool   `gorm:"column:is_active"`
	}

	err := r.db.WithContext(ctx).Raw(
		`SELECT u.id::text,
		        t.id::text AS tenant_id,
		        u.email,
		        u.full_name,
		        u.role,
		        u.is_active
		 FROM tenants t
		 INNER JOIN users u ON u.id = t.user_id
		 WHERE u.id = NULLIF(?, '')::uuid
		   AND t.id = NULLIF(?, '')::uuid
		   AND t.deleted_at IS NULL
		   AND u.deleted_at IS NULL
		 LIMIT 1`,
		userID,
		tenantID,
	).Scan(&row).Error
	if err != nil {
		return nil, fmt.Errorf("find user by id and tenant: %w", err)
	}
	if row.ID == "" {
		return nil, ErrUserNotFound
	}

	return &domain.User{
		ID:       row.ID,
		TenantID: row.TenantID,
		Email:    row.Email,
		FullName: row.FullName,
		Role:     decodeRole(row.Role),
		IsActive: row.IsActive,
	}, nil
}

func (r *userRepository) UpdateByIDAndTenant(ctx context.Context, tenantID, userID string, input UpdateUserInput) (*domain.User, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	updates := map[string]interface{}{}
	if input.Email != nil {
		updates["email"] = strings.TrimSpace(*input.Email)
	}
	if input.FullName != nil {
		updates["full_name"] = strings.TrimSpace(*input.FullName)
	}
	if input.Role != nil {
		updates["role"] = encodeRole(*input.Role)
	}
	if len(updates) == 0 {
		return nil, ErrNoFieldsToUpdate
	}
	updates["updated_at"] = gorm.Expr("NOW()")

	result := r.db.WithContext(ctx).
		Table("users").
		Where("id = NULLIF(?, '')::uuid AND id IN (SELECT user_id FROM tenants WHERE id = NULLIF(?, '')::uuid AND deleted_at IS NULL) AND deleted_at IS NULL", userID, tenantID).
		Updates(updates)
	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrEmailAlreadyExist
		}
		return nil, fmt.Errorf("update user by id and tenant: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, ErrUserNotFound
	}

	return r.FindByIDAndTenant(ctx, tenantID, userID)
}

func (r *userRepository) SoftDeleteByIDAndTenant(ctx context.Context, tenantID, userID string) error {
	if r.db == nil {
		return fmt.Errorf("database is not initialized")
	}

	result := r.db.WithContext(ctx).
		Exec(
			`UPDATE users
			 SET deleted_at = NOW(),
			     updated_at = NOW()
			 WHERE id = NULLIF(?, '')::uuid
			   AND id IN (SELECT user_id FROM tenants WHERE id = NULLIF(?, '')::uuid AND deleted_at IS NULL)
			   AND deleted_at IS NULL`,
			userID,
			tenantID,
		)
	if result.Error != nil {
		return fmt.Errorf("soft delete user by id and tenant: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func (r *userRepository) ToggleActiveByIDAndTenant(ctx context.Context, tenantID, userID string) (*domain.User, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	var row struct {
		ID       string `gorm:"column:id"`
		TenantID string `gorm:"column:tenant_id"`
		Email    string `gorm:"column:email"`
		FullName string `gorm:"column:full_name"`
		Role     string `gorm:"column:role"`
		IsActive bool   `gorm:"column:is_active"`
	}

	err := r.db.WithContext(ctx).Raw(
		`UPDATE users
		 SET is_active = NOT is_active,
		     updated_at = NOW()
		 WHERE id = NULLIF(?, '')::uuid
		   AND id IN (SELECT user_id FROM tenants WHERE id = NULLIF(?, '')::uuid AND deleted_at IS NULL)
		   AND deleted_at IS NULL
		 RETURNING id::text,
		           NULLIF(?, '')::uuid::text AS tenant_id,
		           email,
		           full_name,
		           role,
		           is_active`,
		userID,
		tenantID,
		tenantID,
	).Scan(&row).Error
	if err != nil {
		return nil, fmt.Errorf("toggle active user by id and tenant: %w", err)
	}
	if row.ID == "" {
		return nil, ErrUserNotFound
	}

	return &domain.User{
		ID:       row.ID,
		TenantID: row.TenantID,
		Email:    row.Email,
		FullName: row.FullName,
		Role:     decodeRole(row.Role),
		IsActive: row.IsActive,
	}, nil
}

func (r *userRepository) ListAllUsersForAdmin(ctx context.Context, roleFilter string, limit, offset int) ([]domain.User, int64, error) {
	if r.db == nil {
		return nil, 0, fmt.Errorf("database is not initialized")
	}

	// 1. Hitung total items
	var totalItems int64
	countQuery := `SELECT COUNT(*) FROM users u WHERE u.deleted_at IS NULL AND u.role::text != 'ADMIN'`
	var countArgs []interface{}
	if roleFilter != "" {
		countQuery += " AND u.role::text = ?"
		countArgs = append(countArgs, encodeRole(roleFilter))
	}
	if err := r.db.WithContext(ctx).Raw(countQuery, countArgs...).Scan(&totalItems).Error; err != nil {
		return nil, 0, fmt.Errorf("count all users for admin: %w", err)
	}

	// 2. Ambil data dengan LIMIT dan OFFSET
	var rows []struct {
		ID       string `gorm:"column:id"`
		TenantID string `gorm:"column:tenant_id"`
		Email    string `gorm:"column:email"`
		FullName string `gorm:"column:full_name"`
		Role     string `gorm:"column:role"`
		IsActive bool   `gorm:"column:is_active"`
	}

	query := `SELECT u.id::text,
	                 '' AS tenant_id,
	                 u.email,
	                 u.full_name,
	                 u.role,
	                 u.is_active
	          FROM users u
	          WHERE u.deleted_at IS NULL AND u.role::text != 'ADMIN'`

	var args []interface{}
	if roleFilter != "" {
		query += " AND u.role::text = ?"
		args = append(args, encodeRole(roleFilter))
	}
	query += " ORDER BY u.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	err := r.db.WithContext(ctx).Raw(query, args...).Scan(&rows).Error
	if err != nil {
		return nil, 0, fmt.Errorf("list all users for admin: %w", err)
	}

	users := make([]domain.User, 0, len(rows))
	for _, row := range rows {
		users = append(users, domain.User{
			ID:       row.ID,
			TenantID: row.TenantID,
			Email:    row.Email,
			FullName: row.FullName,
			Role:     decodeRole(row.Role),
			IsActive: row.IsActive,
		})
	}

	return users, totalItems, nil
}

func (r *userRepository) FindByIDForAdmin(ctx context.Context, userID string) (*domain.User, error) {
	if r.db == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	var row userRow
	query := `SELECT u.id::text,
	                 '' AS tenant_id,
	                 u.email,
	                 u.full_name,
	                 u.role,
	                 u.is_active
	          FROM users u
	          WHERE u.id = ? AND u.deleted_at IS NULL AND u.role::text != 'ADMIN'`

	err := r.db.WithContext(ctx).Raw(query, userID).Scan(&row).Error
	if err != nil {
		return nil, fmt.Errorf("find user by id for admin: %w", err)
	}
	if row.ID == "" {
		return nil, ErrUserNotFound
	}

	return &domain.User{
		ID:       row.ID,
		TenantID: row.TenantID,
		Email:    row.Email,
		FullName: row.FullName,
		Role:     decodeRole(row.Role),
		IsActive: row.IsActive,
	}, nil
}
