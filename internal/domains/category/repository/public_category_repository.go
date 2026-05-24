package repository

import (
	"context"

	categorydomain "github.com/dhegas/saas_gangsta/internal/domains/category/domain"
	"gorm.io/gorm"
)

type publicCategoryRepository struct {
	db *gorm.DB
}

func NewPublicCategoryRepository(db *gorm.DB) categorydomain.PublicCategoryRepository {
	return &publicCategoryRepository{db: db}
}

func (r *publicCategoryRepository) FindPublicCategories(ctx context.Context, tenantID string) ([]categorydomain.CategoryEntity, error) {
	var categories []categorydomain.CategoryEntity
	err := r.db.WithContext(ctx).Model(&categorydomain.CategoryEntity{}).
		Where(
			"tenant_id = ?", tenantID).
		Where("is_active = true").
		Where("deleted_at IS NULL").
		Order("sort_order ASC, name ASC").
		Find(&categories).Error
	return categories, err
}
