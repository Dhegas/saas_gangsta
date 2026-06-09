package usecase

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dhegas/saas_gangsta/internal/common/cache"
	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/domains/category/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/category/dto"
)

type publicCategoryUsecase struct {
	repo  domain.PublicCategoryRepository
	cache *cache.LocalCache
}

func NewPublicCategoryUsecase(repo domain.PublicCategoryRepository, cache *cache.LocalCache) domain.PublicCategoryUsecase {
	return &publicCategoryUsecase{
		repo:  repo,
		cache: cache,
	}
}

func (u *publicCategoryUsecase) GetPublicCategories(ctx context.Context, tenantID string) ([]dto.CategoryResponse, error) {
	cacheKey := fmt.Sprintf("public:categories:tenant:%s", tenantID)

	if cached, found := u.cache.Get(cacheKey); found {
		if cachedCategories, ok := cached.([]dto.CategoryResponse); ok {
			return cachedCategories, nil
		}
	}

	categories, err := u.repo.FindPublicCategories(ctx, tenantID)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil data kategori publik", http.StatusInternalServerError)
	}

	result := make([]dto.CategoryResponse, 0, len(categories))
	for _, c := range categories {
		result = append(result, toCategoryResponse(&c))
	}

	u.cache.Set(cacheKey, result, 10*time.Minute)

	return result, nil
}

