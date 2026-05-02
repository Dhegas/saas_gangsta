package domain

import (
	"context"
	"time"

	"github.com/dhegas/saas_gangsta/internal/domains/category/dto"
)

// PartnerCategoryUsecase adalah kontrak untuk logika bisnis Category oleh Partner (PARTNER)
type PartnerCategoryUsecase interface {
	GetAllCategories(ctx context.Context, tenantID string) ([]dto.CategoryResponse, error)
	GetCategoryByID(ctx context.Context, tenantID, categoryID string) (*dto.CategoryResponse, error)
	CreateCategory(ctx context.Context, tenantID string, req dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	UpdateCategory(ctx context.Context, tenantID, categoryID string, req dto.UpdateCategoryRequest) (*dto.CategoryResponse, error)
	SoftDeleteCategory(ctx context.Context, tenantID, categoryID string) error
	ToggleCategoryActive(ctx context.Context, tenantID, categoryID string, isActive bool) error
	ReorderCategories(ctx context.Context, tenantID string, req dto.ReorderCategoryRequest) error
}

// PartnerCategoryRepository adalah kontrak untuk query ke tabel categories
type PartnerCategoryRepository interface {
	FindAllByTenant(ctx context.Context, tenantID string) ([]CategoryEntity, error)
	FindByIDAndTenant(ctx context.Context, tenantID, categoryID string) (*CategoryEntity, error)
	Create(ctx context.Context, entity *CategoryEntity) error
	Update(ctx context.Context, entity *CategoryEntity) error
	SoftDelete(ctx context.Context, tenantID, categoryID string) error
	UpdateActiveStatus(ctx context.Context, tenantID, categoryID string, isActive bool) error
	UpdateSortOrderBulk(ctx context.Context, tenantID string, items []dto.CategoryOrder) error
	CheckNameExists(ctx context.Context, tenantID, name string, excludeID string) (bool, error)
}

// CategoryEntity merepresentasikan struktur tabel categories di database untuk GORM
type CategoryEntity struct {
	ID          string     `gorm:"primaryKey;default:gen_random_uuid()"`
	TenantID    string     `gorm:"index;not null"`
	Name        string     `gorm:"not null"`
	Description string     
	SortOrder   int        `gorm:"not null;default:0"`
	IsActive    bool       `gorm:"not null;default:true"`
	CreatedAt   time.Time  `gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime"`
	DeletedAt   *time.Time `gorm:"index"`
}

// TableName meng-override nama tabel
func (CategoryEntity) TableName() string {
	return "categories"
}
