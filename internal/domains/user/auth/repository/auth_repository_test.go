package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sqlmock: %v", err)
	}

	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("open gorm db: %v", err)
	}

	cleanup := func() {
		_ = sqlDB.Close()
	}
	return gdb, mock, cleanup
}

func TestFindByEmailSuccess(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewAuthRepository(db)
	rows := sqlmock.NewRows([]string{"id", "tenant_id", "email", "password_hash", "role", "is_active", "tenant_status"}).
		AddRow("u-1", "t-1", "user@test.local", "hash", "PARTNER", true, "active")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT u.id, COALESCE(u.tenant_id::text, '') AS tenant_id, u.email, u.password_hash, u.role, u.is_active, COALESCE(t.status, 'active') AS tenant_status FROM users u LEFT JOIN tenants t ON t.id = u.tenant_id WHERE LOWER(u.email) = LOWER($1) LIMIT $2`)).
		WithArgs("user@test.local", 1).
		WillReturnRows(rows)

	user, err := repo.FindByEmail(context.Background(), "user@test.local")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if user == nil || user.ID != "u-1" {
		t.Fatalf("expected user u-1")
	}
}

func TestFindByEmailNilDB(t *testing.T) {
	repo := NewAuthRepository(nil)
	if _, err := repo.FindByEmail(context.Background(), "user@test.local"); err == nil {
		t.Fatalf("expected error when db is nil")
	}
}

func TestFindByIDNilDB(t *testing.T) {
	repo := NewAuthRepository(nil)
	if _, err := repo.FindByID(context.Background(), "u-1"); err == nil {
		t.Fatalf("expected error when db is nil")
	}
}

func TestFindByIDSuccess(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewAuthRepository(db)
	rows := sqlmock.NewRows([]string{"id", "tenant_id", "email", "password_hash", "role", "is_active", "tenant_status"}).
		AddRow("u-2", "t-2", "user2@test.local", "hash", "CUSTOMER", true, "active")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT u.id, COALESCE(u.tenant_id::text, '') AS tenant_id, u.email, u.password_hash, u.role, u.is_active, COALESCE(t.status, 'active') AS tenant_status FROM users u LEFT JOIN tenants t ON t.id = u.tenant_id WHERE u.id = $1 LIMIT $2`)).
		WithArgs("u-1", 1).
		WillReturnRows(rows)

	user, err := repo.FindByID(context.Background(), "u-1")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if user == nil || user.ID == "" {
		t.Fatalf("expected user data")
	}
}

func TestFindByEmailNotFound(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	repo := NewAuthRepository(db)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT u.id, COALESCE(u.tenant_id::text, '') AS tenant_id, u.email, u.password_hash, u.role, u.is_active, COALESCE(t.status, 'active') AS tenant_status FROM users u LEFT JOIN tenants t ON t.id = u.tenant_id WHERE LOWER(u.email) = LOWER($1) LIMIT $2`)).
		WithArgs("none@test.local", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	user, err := repo.FindByEmail(context.Background(), "none@test.local")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if user != nil {
		t.Fatalf("expected nil user when not found")
	}
}
