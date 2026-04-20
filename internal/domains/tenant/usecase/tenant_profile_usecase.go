package usecase

import (
	"context"
	"net/http"
	"strings"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant/dto"
)

type tenantProfileUsecase struct {
	repo domain.TenantProfileRepository
}

// NewTenantProfileUsecase membuat usecase tenant profile.
func NewTenantProfileUsecase(repo domain.TenantProfileRepository) domain.TenantProfileUsecase {
	return &tenantProfileUsecase{repo: repo}
}

func (u *tenantProfileUsecase) CreateProfile(ctx context.Context, tenantID string, req dto.CreateTenantProfileRequest) (*dto.TenantProfileResponse, error) {
	if strings.TrimSpace(tenantID) == "" {
		return nil, apperrors.New("TENANT_NOT_FOUND", "Tenant context is required", http.StatusUnauthorized, nil)
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, apperrors.New("VALIDATION_ERROR", "name wajib diisi", http.StatusBadRequest, nil)
	}

	sortOrder := 0
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	entity := &domain.TenantProfileEntity{
		TenantID:    tenantID,
		Name:        name,
		Description: strings.TrimSpace(req.Description),
		SortOrder:   sortOrder,
		IsActive:    isActive,
	}

	if err := u.repo.Create(ctx, entity); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return nil, apperrors.New("CONFLICT", "Nama tenant profile sudah digunakan", http.StatusConflict, nil)
		}
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal membuat tenant profile", http.StatusInternalServerError, nil)
	}

	return &dto.TenantProfileResponse{
		ID:          entity.ID,
		TenantID:    entity.TenantID,
		Name:        entity.Name,
		Description: entity.Description,
		SortOrder:   entity.SortOrder,
		IsActive:    entity.IsActive,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}, nil
}

func (u *tenantProfileUsecase) ListProfiles(ctx context.Context, tenantID string) ([]dto.TenantProfileResponse, error) {
	if strings.TrimSpace(tenantID) == "" {
		return nil, apperrors.New("TENANT_NOT_FOUND", "Tenant context is required", http.StatusUnauthorized, nil)
	}

	entities, err := u.repo.ListByTenantID(ctx, tenantID)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil tenant profile", http.StatusInternalServerError, nil)
	}

	responses := make([]dto.TenantProfileResponse, 0, len(entities))
	for _, entity := range entities {
		responses = append(responses, dto.TenantProfileResponse{
			ID:          entity.ID,
			TenantID:    entity.TenantID,
			Name:        entity.Name,
			Description: entity.Description,
			SortOrder:   entity.SortOrder,
			IsActive:    entity.IsActive,
			CreatedAt:   entity.CreatedAt,
			UpdatedAt:   entity.UpdatedAt,
		})
	}

	return responses, nil
}
