package repository

import (
	"context"
	"time"

	tabledomain "github.com/dhegas/saas_gangsta/internal/domains/table/domain"
	"gorm.io/gorm"
)

type merchantTableRepository struct {
	db *gorm.DB
}

func NewMerchantTableRepository(db *gorm.DB) tabledomain.MerchantTableRepository {
	return &merchantTableRepository{db: db}
}

func (r *merchantTableRepository) FindAllByTenant(ctx context.Context, tenantID string) ([]tabledomain.DiningTableEntity, error) {
	var tables []tabledomain.DiningTableEntity
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID).
		Order("created_at DESC").
		Find(&tables).Error
	return tables, err
}

func (r *merchantTableRepository) FindByIDAndTenant(ctx context.Context, tenantID, tableID string) (*tabledomain.DiningTableEntity, error) {
	var table tabledomain.DiningTableEntity
	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", tableID, tenantID).
		First(&table).Error
	if err != nil {
		return nil, err
	}
	return &table, nil
}

func (r *merchantTableRepository) Create(ctx context.Context, entity *tabledomain.DiningTableEntity) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *merchantTableRepository) Update(ctx context.Context, entity *tabledomain.DiningTableEntity) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *merchantTableRepository) SoftDelete(ctx context.Context, tenantID, tableID string) error {
	now := time.Now()
	res := r.db.WithContext(ctx).
		Model(&tabledomain.DiningTableEntity{}).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", tableID, tenantID).
		Update("deleted_at", &now)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *merchantTableRepository) CheckNameExists(ctx context.Context, tenantID, tableName string, excludeID string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&tabledomain.DiningTableEntity{}).
		Where("tenant_id = ? AND table_name = ? AND deleted_at IS NULL", tenantID, tableName)

	if excludeID != "" {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *merchantTableRepository) CheckTableOccupied(ctx context.Context, tableID string) (bool, error) {
	var count int64
	// Mengecek apakah ada order yang menggunakan meja ini dan statusnya bukan COMPLETED atau CANCELLED
	err := r.db.WithContext(ctx).Table("orders").
		Where("dining_tables_id = ? AND status NOT IN (?, ?) AND deleted_at IS NULL", tableID, "COMPLETED", "CANCELLED").
		Count(&count).Error
	
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
