package repository

import (
	"context"
	"time"

	"github.com/dhegas/saas_gangsta/internal/domains/category/domain"
	menudomain "github.com/dhegas/saas_gangsta/internal/domains/menu/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/menu/dto"
	"gorm.io/gorm"
)

type merchantMenuRepository struct {
	db *gorm.DB
}

func NewMerchantMenuRepository(db *gorm.DB) menudomain.MerchantMenuRepository {
	return &merchantMenuRepository{db: db}
}

func (r *merchantMenuRepository) FindAllByTenant(ctx context.Context, tenantID string, filter dto.MenuFilterParams) ([]menudomain.MenuEntity, error) {
	var menus []menudomain.MenuEntity
	query := r.db.WithContext(ctx).
		Where("tenant_id = ? AND deleted_at IS NULL", tenantID)

	if filter.CategoryID != "" {
		query = query.Where("category_id = ?", filter.CategoryID)
	}
	if filter.IsAvailable != nil {
		query = query.Where("is_available = ?", *filter.IsAvailable)
	}

	err := query.Order("created_at DESC").Find(&menus).Error
	return menus, err
}

func (r *merchantMenuRepository) FindByIDAndTenant(ctx context.Context, tenantID, menuID string) (*menudomain.MenuEntity, error) {
	var menu menudomain.MenuEntity
	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", menuID, tenantID).
		First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *merchantMenuRepository) Create(ctx context.Context, entity *menudomain.MenuEntity) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *merchantMenuRepository) Update(ctx context.Context, entity *menudomain.MenuEntity) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *merchantMenuRepository) SoftDelete(ctx context.Context, tenantID, menuID string) error {
	now := time.Now()
	res := r.db.WithContext(ctx).
		Model(&menudomain.MenuEntity{}).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", menuID, tenantID).
		Update("deleted_at", &now)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *merchantMenuRepository) UpdateAvailableStatus(ctx context.Context, tenantID, menuID string, isAvailable bool) error {
	res := r.db.WithContext(ctx).
		Model(&menudomain.MenuEntity{}).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", menuID, tenantID).
		Update("is_available", isAvailable)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *merchantMenuRepository) CheckNameExists(ctx context.Context, tenantID, name string, excludeID string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&menudomain.MenuEntity{}).
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

func (r *merchantMenuRepository) CategoryExists(ctx context.Context, tenantID, categoryID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.CategoryEntity{}).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", categoryID, tenantID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
