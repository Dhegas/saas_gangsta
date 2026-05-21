package database

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  databaseURL,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}

	// Safely auto-migrate schema additions
	_ = db.Exec(`ALTER TABLE tenants ADD COLUMN IF NOT EXISTS is_public BOOLEAN DEFAULT true`)
	_ = db.Exec(`UPDATE tenants SET is_public = true WHERE is_public IS NULL`)
	_ = db.Exec(`CREATE INDEX IF NOT EXISTS idx_tenants_is_public ON tenants (is_public) WHERE deleted_at IS NULL`)

	return db, nil
}

func IsReady(ctx context.Context, db *gorm.DB) error {
	if db == nil {
		return fmt.Errorf("database client is not initialized")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("open sql db: %w", err)
	}

	if deadline, ok := ctx.Deadline(); !ok || deadline.Before(time.Now()) {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	return nil
}
