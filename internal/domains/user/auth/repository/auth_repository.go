package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dhegas/saas_gangsta/internal/domains/user/auth/domain"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type AuthRepository interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
	FindPhoneNumber(ctx context.Context, userID string, role string) (string, error)
}

type authRepository struct {
	db *gorm.DB
}

type authUserRow struct {
	ID           string `gorm:"column:id"`
	TenantID     string `gorm:"column:tenant_id"`
	Email        string `gorm:"column:email"`
	FullName     string `gorm:"column:full_name"`
	PasswordHash string `gorm:"column:password_hash"`
	Role         string `gorm:"column:role"`
	IsActive     bool   `gorm:"column:is_active"`
	TenantStatus string `gorm:"column:tenant_status"`
}

func decodeRole(role string) string {
	return strings.ToUpper(strings.TrimSpace(role))
}

func encodeRole(role string) string {
	return strings.ToUpper(strings.TrimSpace(role))
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
			"u.id, COALESCE(t.id::text, '') AS tenant_id, u.email, u.full_name, u.password_hash, u.role, u.is_active, COALESCE(t.status, 'active') AS tenant_status",
		).
		Joins("LEFT JOIN LATERAL (SELECT id, status FROM tenants WHERE user_id = u.id AND deleted_at IS NULL ORDER BY created_at DESC LIMIT 1) t ON true").
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
		FullName:     row.FullName,
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
			"u.id, COALESCE(t.id::text, '') AS tenant_id, u.email, u.full_name, u.password_hash, u.role, u.is_active, COALESCE(t.status, 'active') AS tenant_status",
		).
		Joins("LEFT JOIN LATERAL (SELECT id, status FROM tenants WHERE user_id = u.id AND deleted_at IS NULL ORDER BY created_at DESC LIMIT 1) t ON true").
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
		FullName:     row.FullName,
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
		`INSERT INTO users (email, password_hash, full_name, role, is_active, created_at, updated_at)
		 VALUES (?, ?, ?, ?, TRUE, NOW(), NOW())
		 RETURNING id::text`,
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

func (r *authRepository) FindPhoneNumber(ctx context.Context, userID string, role string) (string, error) {
	if r.db == nil {
		return "", fmt.Errorf("database is not initialized")
	}

	var phoneNumber string
	var err error

	if strings.ToUpper(role) == "PARTNER" {
		err = r.db.WithContext(ctx).
			Table("tenants").
			Select("phone_number").
			Where("user_id = ? AND deleted_at IS NULL", userID).
			Order("created_at DESC").
			Limit(1).
			Scan(&phoneNumber).
			Error
	} else if strings.ToUpper(role) == "CUSTOMER" {
		phoneNumber = ""
	}

	if err != nil {
		return "", err
	}
	return phoneNumber, nil
}

