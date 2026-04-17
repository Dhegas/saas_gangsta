package repository

import (
	"context"
	"fmt"

	"github.com/dhegas/saas_gangsta/internal/modules/auth/domain"
	"gorm.io/gorm"
)

type AuthRepository interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
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
		Role:         row.Role,
		IsActive:     row.IsActive,
		TenantStatus: row.TenantStatus,
	}, nil
}
