package domain

import (
	"context"

	"github.com/dhegas/saas_gangsta/internal/domains/category/dto"
)

type PublicCategoryUsecase interface {
	GetPublicCategories(ctx context.Context, tenantID string) ([]dto.CategoryResponse, error)
}

type PublicCategoryRepository interface {
	FindPublicCategories(ctx context.Context, tenantID string) ([]CategoryEntity, error)
}
