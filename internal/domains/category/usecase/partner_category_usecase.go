package usecase

import (
	"context"
	"net/http"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/domains/category/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/category/dto"
	"gorm.io/gorm"
	"errors"
)

type partnerCategoryUsecase struct {
	repo domain.PartnerCategoryRepository
}

func NewPartnerCategoryUsecase(repo domain.PartnerCategoryRepository) domain.PartnerCategoryUsecase {
	return &partnerCategoryUsecase{repo: repo}
}

func (u *partnerCategoryUsecase) GetAllCategories(ctx context.Context, tenantID string) ([]dto.CategoryResponse, error) {
	categories, err := u.repo.FindAllByTenant(ctx, tenantID)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil data kategori", http.StatusInternalServerError, err)
	}

	result := make([]dto.CategoryResponse, 0, len(categories))
	for _, c := range categories {
		result = append(result, toCategoryResponse(&c))
	}

	return result, nil
}

func (u *partnerCategoryUsecase) GetCategoryByID(ctx context.Context, tenantID, categoryID string) (*dto.CategoryResponse, error) {
	category, err := u.repo.FindByIDAndTenant(ctx, tenantID, categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New("NOT_FOUND", "Kategori tidak ditemukan", http.StatusNotFound, nil)
		}
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil data kategori", http.StatusInternalServerError, err)
	}

	response := toCategoryResponse(category)
	return &response, nil
}

func (u *partnerCategoryUsecase) CreateCategory(ctx context.Context, tenantID string, req dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	exists, err := u.repo.CheckNameExists(ctx, tenantID, req.Name, "")
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal memvalidasi nama kategori", http.StatusInternalServerError, err)
	}
	if exists {
		return nil, apperrors.New("CONFLICT", "Nama kategori sudah ada", http.StatusConflict, nil)
	}

	entity := &domain.CategoryEntity{
		TenantID:    tenantID,
		Name:        req.Name,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		IsActive:    true, // default
	}

	if err := u.repo.Create(ctx, entity); err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal menyimpan kategori", http.StatusInternalServerError, err)
	}

	response := toCategoryResponse(entity)
	return &response, nil
}

func (u *partnerCategoryUsecase) UpdateCategory(ctx context.Context, tenantID, categoryID string, req dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	category, err := u.repo.FindByIDAndTenant(ctx, tenantID, categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New("NOT_FOUND", "Kategori tidak ditemukan", http.StatusNotFound, nil)
		}
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil data kategori", http.StatusInternalServerError, err)
	}

	if req.Name != "" && req.Name != category.Name {
		exists, err := u.repo.CheckNameExists(ctx, tenantID, req.Name, categoryID)
		if err != nil {
			return nil, apperrors.New("INTERNAL_ERROR", "Gagal memvalidasi nama kategori", http.StatusInternalServerError, err)
		}
		if exists {
			return nil, apperrors.New("CONFLICT", "Nama kategori sudah ada", http.StatusConflict, nil)
		}
		category.Name = req.Name
	}

	if req.Description != "" {
		category.Description = req.Description
	}

	if err := u.repo.Update(ctx, category); err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal memperbarui kategori", http.StatusInternalServerError, err)
	}

	response := toCategoryResponse(category)
	return &response, nil
}

func (u *partnerCategoryUsecase) SoftDeleteCategory(ctx context.Context, tenantID, categoryID string) error {
	err := u.repo.SoftDelete(ctx, tenantID, categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.New("NOT_FOUND", "Kategori tidak ditemukan", http.StatusNotFound, nil)
		}
		return apperrors.New("INTERNAL_ERROR", "Gagal menghapus kategori", http.StatusInternalServerError, err)
	}
	return nil
}

func (u *partnerCategoryUsecase) ToggleCategoryActive(ctx context.Context, tenantID, categoryID string, isActive bool) error {
	err := u.repo.UpdateActiveStatus(ctx, tenantID, categoryID, isActive)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.New("NOT_FOUND", "Kategori tidak ditemukan", http.StatusNotFound, nil)
		}
		return apperrors.New("INTERNAL_ERROR", "Gagal memperbarui status kategori", http.StatusInternalServerError, err)
	}
	return nil
}

func (u *partnerCategoryUsecase) ReorderCategories(ctx context.Context, tenantID string, req dto.ReorderCategoryRequest) error {
	err := u.repo.UpdateSortOrderBulk(ctx, tenantID, req.Items)
	if err != nil {
		return apperrors.New("INTERNAL_ERROR", "Gagal mengurutkan kategori", http.StatusInternalServerError, err)
	}
	return nil
}

func toCategoryResponse(entity *domain.CategoryEntity) dto.CategoryResponse {
	return dto.CategoryResponse{
		ID:          entity.ID,
		TenantID:    entity.TenantID,
		Name:        entity.Name,
		Description: entity.Description,
		SortOrder:   entity.SortOrder,
		IsActive:    entity.IsActive,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		DeletedAt:   entity.DeletedAt,
	}
}
