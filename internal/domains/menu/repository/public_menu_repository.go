package repository

import (
	"context"
	"strings"

	menudomain "github.com/dhegas/saas_gangsta/internal/domains/menu/domain"
	"gorm.io/gorm"
)

type publicMenuRepository struct {
	db *gorm.DB
}

func NewPublicMenuRepository(db *gorm.DB) menudomain.PublicMenuRepository {
	return &publicMenuRepository{db: db}
}

func (r *publicMenuRepository) FindPublicMenus(ctx context.Context, tenantID string, categoryID string, search string, isAvailable *bool) ([]menudomain.MenuEntity, error) {
	var menus []menudomain.MenuEntity
	query := r.db.WithContext(ctx).Model(&menudomain.MenuEntity{}).
		Where("tenant_id = ?", tenantID).
		Where("deleted_at IS NULL")

	// Public customers should see available menus by default
	if isAvailable != nil {
		query = query.Where("is_available = ?", *isAvailable)
	} else {
		query = query.Where("is_available = true")
	}

	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchPattern, searchPattern)
	}

	err := query.Order("name ASC").Find(&menus).Error
	return menus, err
}
