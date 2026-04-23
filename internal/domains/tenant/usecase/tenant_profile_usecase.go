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

	return entityToResponse(entity), nil
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
	for i := range entities {
		responses = append(responses, *entityToResponse(&entities[i]))
	}

	return responses, nil
}

func (u *tenantProfileUsecase) GetProfileByID(ctx context.Context, tenantID, profileID string) (*dto.TenantProfileResponse, error) {
	if err := validateTenantAndProfile(tenantID, profileID); err != nil {
		return nil, err
	}

	entity, err := u.repo.FindByID(ctx, tenantID, profileID)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil tenant profile", http.StatusInternalServerError, nil)
	}
	if entity == nil {
		return nil, apperrors.New("NOT_FOUND", "Tenant profile tidak ditemukan", http.StatusNotFound, nil)
	}

	return entityToResponse(entity), nil
}

func (u *tenantProfileUsecase) UpdateProfile(ctx context.Context, tenantID, profileID string, req dto.UpdateTenantProfileRequest) (*dto.TenantProfileResponse, error) {
	if err := validateTenantAndProfile(tenantID, profileID); err != nil {
		return nil, err
	}

	// Ambil data yang ada terlebih dahulu
	entity, err := u.repo.FindByID(ctx, tenantID, profileID)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil tenant profile", http.StatusInternalServerError, nil)
	}
	if entity == nil {
		return nil, apperrors.New("NOT_FOUND", "Tenant profile tidak ditemukan", http.StatusNotFound, nil)
	}

	// Terapkan perubahan (partial update)
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, apperrors.New("VALIDATION_ERROR", "name tidak boleh kosong", http.StatusBadRequest, nil)
		}
		entity.Name = name
	}
	if req.Description != nil {
		entity.Description = strings.TrimSpace(*req.Description)
	}
	if req.SortOrder != nil {
		entity.SortOrder = *req.SortOrder
	}
	if req.IsActive != nil {
		entity.IsActive = *req.IsActive
	}

	if err := u.repo.Update(ctx, entity); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return nil, apperrors.New("CONFLICT", "Nama tenant profile sudah digunakan", http.StatusConflict, nil)
		}
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengupdate tenant profile", http.StatusInternalServerError, nil)
	}

	return entityToResponse(entity), nil
}

func (u *tenantProfileUsecase) DeleteProfile(ctx context.Context, tenantID, profileID string) error {
	if err := validateTenantAndProfile(tenantID, profileID); err != nil {
		return err
	}

	if err := u.repo.SoftDelete(ctx, tenantID, profileID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return apperrors.New("NOT_FOUND", "Tenant profile tidak ditemukan", http.StatusNotFound, nil)
		}
		return apperrors.New("INTERNAL_ERROR", "Gagal menghapus tenant profile", http.StatusInternalServerError, nil)
	}

	return nil
}

func (u *tenantProfileUsecase) ToggleActive(ctx context.Context, tenantID, profileID string) (*dto.TenantProfileResponse, error) {
	if err := validateTenantAndProfile(tenantID, profileID); err != nil {
		return nil, err
	}

	entity, err := u.repo.FindByID(ctx, tenantID, profileID)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil tenant profile", http.StatusInternalServerError, nil)
	}
	if entity == nil {
		return nil, apperrors.New("NOT_FOUND", "Tenant profile tidak ditemukan", http.StatusNotFound, nil)
	}

	// Balik status aktif
	entity.IsActive = !entity.IsActive

	if err := u.repo.Update(ctx, entity); err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengubah status tenant profile", http.StatusInternalServerError, nil)
	}

	return entityToResponse(entity), nil
}

// --- helpers ---

func entityToResponse(e *domain.TenantProfileEntity) *dto.TenantProfileResponse {
	return &dto.TenantProfileResponse{
		ID:          e.ID,
		TenantID:    e.TenantID,
		Name:        e.Name,
		Description: e.Description,
		SortOrder:   e.SortOrder,
		IsActive:    e.IsActive,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

func validateTenantAndProfile(tenantID, profileID string) error {
	if strings.TrimSpace(tenantID) == "" {
		return apperrors.New("TENANT_NOT_FOUND", "Tenant context is required", http.StatusUnauthorized, nil)
	}
	if strings.TrimSpace(profileID) == "" {
		return apperrors.New("VALIDATION_ERROR", "Profile ID is required", http.StatusBadRequest, nil)
	}
	return nil
}
