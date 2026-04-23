package repository

import (
	"context"
	"time"

	"github.com/dhegas/saas_gangsta/internal/domains/category/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/category/dto"
	"gorm.io/gorm"
)

type merchantCategoryRepository struct {
	db *gorm.DB
}

func NewMerchantCategoryRepository(db *gorm.DB) domain.MerchantCategoryRepository {
	return &merchantCategoryRepository{db: db}
}

func (r *merchantCategoryRepository) FindAllByTenant(ctx context.Context, tenantID string) ([]domain.CategoryEntity, error) {
	var categories []domain.CategoryEntity
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID).
		Order("sort_order ASC, created_at DESC").
		Find(&categories).Error
	return categories, err
}

func (r *merchantCategoryRepository) FindByIDAndTenant(ctx context.Context, tenantID, categoryID string) (*domain.CategoryEntity, error) {
	var category domain.CategoryEntity
	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", categoryID, tenantID).
		First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *merchantCategoryRepository) Create(ctx context.Context, entity *domain.CategoryEntity) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *merchantCategoryRepository) Update(ctx context.Context, entity *domain.CategoryEntity) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *merchantCategoryRepository) SoftDelete(ctx context.Context, tenantID, categoryID string) error {
	now := time.Now()
	res := r.db.WithContext(ctx).
		Model(&domain.CategoryEntity{}).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", categoryID, tenantID).
		Update("deleted_at", &now)
		
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *merchantCategoryRepository) UpdateActiveStatus(ctx context.Context, tenantID, categoryID string, isActive bool) error {
	res := r.db.WithContext(ctx).
		Model(&domain.CategoryEntity{}).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", categoryID, tenantID).
		Update("is_active", isActive)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *merchantCategoryRepository) UpdateSortOrderBulk(ctx context.Context, tenantID string, items []dto.CategoryOrder) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			res := tx.Model(&domain.CategoryEntity{}).
				Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", item.ID, tenantID).
				Update("sort_order", item.SortOrder)
			if res.Error != nil {
				return res.Error
			}
		}
		return nil
	})
}

func (r *merchantCategoryRepository) CheckNameExists(ctx context.Context, tenantID, name string, excludeID string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&domain.CategoryEntity{}).
		Where("tenant_id = ? AND name = ? AND deleted_at IS NULL", tenantID, name)

	if excludeID != "" {
		query = query.Where("id != ?", excludeID)
	}

	err := query.Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
