package usecase

import (
	"context"
	"net/http"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/domains/category/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/category/dto"
)

type publicCategoryUsecase struct {
	repo domain.PublicCategoryRepository
}

func NewPublicCategoryUsecase(repo domain.PublicCategoryRepository) domain.PublicCategoryUsecase {
	return &publicCategoryUsecase{repo: repo}
}

func (u *publicCategoryUsecase) GetPublicCategories(ctx context.Context, tenantID string) ([]dto.CategoryResponse, error) {
	categories, err := u.repo.FindPublicCategories(ctx, tenantID)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil data kategori publik", http.StatusInternalServerError, err)
	}

	result := make([]dto.CategoryResponse, 0, len(categories))
	for _, c := range categories {
		result = append(result, toCategoryResponse(&c))
	}

	return result, nil
}
