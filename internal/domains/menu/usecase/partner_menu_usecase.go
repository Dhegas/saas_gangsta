package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dhegas/saas_gangsta/internal/common/cache"
	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/domains/menu/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/menu/dto"
	"gorm.io/gorm"
)

type partnerMenuUsecase struct {
	repo  domain.PartnerMenuRepository
	cache *cache.LocalCache
}

func NewPartnerMenuUsecase(repo domain.PartnerMenuRepository, cache *cache.LocalCache) domain.PartnerMenuUsecase {
	return &partnerMenuUsecase{
		repo:  repo,
		cache: cache,
	}
}

func (u *partnerMenuUsecase) GetAllMenus(ctx context.Context, tenantID string, filter dto.MenuFilterParams) ([]dto.MenuResponse, error) {
	// 1. Validasi kategori jika ada di filter
	if filter.CategoryID != "" {
		isValid, err := u.repo.CategoryExists(ctx, tenantID, filter.CategoryID)
		if err != nil {
			// Jika error berasal dari pengecekan kepemilikan (misal: kategori milik orang lain)
			return nil, apperrors.New("VALIDATION_ERROR", "Kategori tidak valid atau tidak ditemukan", http.StatusBadRequest)
		}
		if !isValid {
			return nil, apperrors.New("NOT_FOUND", "Kategori tidak ditemukan atau bukan milik tenant Anda", http.StatusBadRequest)
		}
	}

	// 2. Ambil data menu
	menus, err := u.repo.FindAllByTenant(ctx, tenantID, filter)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil data menu", http.StatusInternalServerError)
	}

	result := make([]dto.MenuResponse, 0, len(menus))
	for _, m := range menus {
		result = append(result, toMenuResponse(&m))
	}

	return result, nil
}

func (u *partnerMenuUsecase) GetMenuByID(ctx context.Context, tenantID, menuID string) (*dto.MenuResponse, error) {
	menu, err := u.repo.FindByIDAndTenant(ctx, tenantID, menuID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New("NOT_FOUND", "Menu tidak ditemukan", http.StatusNotFound)
		}
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil data menu", http.StatusInternalServerError)
	}

	response := toMenuResponse(menu)
	return &response, nil
}

func (u *partnerMenuUsecase) CreateMenu(ctx context.Context, tenantID string, req dto.CreateMenuRequest) (*dto.MenuResponse, error) {
	if req.CategoryID != nil && *req.CategoryID != "" {
		isValidCategory, err := u.repo.CategoryExists(ctx, tenantID, *req.CategoryID)
		if err != nil {
			return nil, apperrors.New("VALIDATION_ERROR", "Kategori tidak valid atau tidak ditemukan", http.StatusBadRequest)
		}
		if !isValidCategory {
			return nil, apperrors.New("VALIDATION_ERROR", "Kategori tidak valid", http.StatusBadRequest)
		}
	}

	exists, err := u.repo.CheckNameExists(ctx, tenantID, req.Name, "")
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal memvalidasi nama menu", http.StatusInternalServerError)
	}
	if exists {
		return nil, apperrors.New("CONFLICT", "Nama menu sudah ada", http.StatusConflict)
	}

	entity := &domain.MenuEntity{
		TenantID:    tenantID,
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		ImageURL:    req.ImageURL,
		IsAvailable: true,
	}

	if err := u.repo.Create(ctx, entity); err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal menyimpan menu", http.StatusInternalServerError)
	}

	if u.cache != nil {
		u.cache.DeleteByPrefix(fmt.Sprintf("public:menus:tenant:%s", tenantID))
	}

	response := toMenuResponse(entity)
	return &response, nil
}

func (u *partnerMenuUsecase) UpdateMenu(ctx context.Context, tenantID, menuID string, req dto.UpdateMenuRequest) (*dto.MenuResponse, error) {
	menu, err := u.repo.FindByIDAndTenant(ctx, tenantID, menuID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New("NOT_FOUND", "Menu tidak ditemukan", http.StatusNotFound)
		}
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil data menu", http.StatusInternalServerError)
	}

	if req.CategoryID != nil {
		if *req.CategoryID != "" {
			isValidCategory, err := u.repo.CategoryExists(ctx, tenantID, *req.CategoryID)
			if err != nil {
				return nil, apperrors.New("VALIDATION_ERROR", "Kategori tidak valid atau tidak ditemukan", http.StatusBadRequest)
			}
			if !isValidCategory {
				return nil, apperrors.New("VALIDATION_ERROR", "Kategori tidak valid", http.StatusBadRequest)
			}
			menu.CategoryID = req.CategoryID
		} else {
			menu.CategoryID = nil
		}
	}

	if req.Name != "" && req.Name != menu.Name {
		exists, err := u.repo.CheckNameExists(ctx, tenantID, req.Name, menuID)
		if err != nil {
			return nil, apperrors.New("INTERNAL_ERROR", "Gagal memvalidasi nama menu", http.StatusInternalServerError)
		}
		if exists {
			return nil, apperrors.New("CONFLICT", "Nama menu sudah ada", http.StatusConflict)
		}
		menu.Name = req.Name
	}

	if req.Description != "" {
		menu.Description = req.Description
	}

	if req.Price != nil {
		menu.Price = *req.Price
	}

	if req.ImageURL != "" {
		menu.ImageURL = req.ImageURL
	}

	if err := u.repo.Update(ctx, menu); err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal memperbarui menu", http.StatusInternalServerError)
	}

	if u.cache != nil {
		u.cache.DeleteByPrefix(fmt.Sprintf("public:menus:tenant:%s", tenantID))
		u.cache.DeleteByPrefix(fmt.Sprintf("public:categories:tenant:%s", tenantID))
	}

	response := toMenuResponse(menu)
	return &response, nil
}

func (u *partnerMenuUsecase) SoftDeleteMenu(ctx context.Context, tenantID, menuID string) error {
	err := u.repo.SoftDelete(ctx, tenantID, menuID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.New("NOT_FOUND", "Menu tidak ditemukan", http.StatusNotFound)
		}
		return apperrors.New("INTERNAL_ERROR", "Gagal menghapus menu", http.StatusInternalServerError)
	}

	if u.cache != nil {
		u.cache.DeleteByPrefix(fmt.Sprintf("public:menus:tenant:%s", tenantID))
		u.cache.DeleteByPrefix(fmt.Sprintf("public:categories:tenant:%s", tenantID))
	}

	return nil
}

func (u *partnerMenuUsecase) ToggleMenuAvailable(ctx context.Context, tenantID, menuID string, isAvailable bool) error {
	err := u.repo.UpdateAvailableStatus(ctx, tenantID, menuID, isAvailable)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.New("NOT_FOUND", "Menu tidak ditemukan", http.StatusNotFound)
		}
		return apperrors.New("INTERNAL_ERROR", "Gagal memperbarui status ketersediaan menu", http.StatusInternalServerError)
	}

	if u.cache != nil {
		u.cache.DeleteByPrefix(fmt.Sprintf("public:menus:tenant:%s", tenantID))
	}

	return nil
}

func toMenuResponse(entity *domain.MenuEntity) dto.MenuResponse {
	return dto.MenuResponse{
		ID:          entity.ID,
		TenantID:    entity.TenantID,
		CategoryID:  entity.CategoryID,
		Name:        entity.Name,
		Description: entity.Description,
		Price:       entity.Price,
		ImageURL:    entity.ImageURL,
		IsAvailable: entity.IsAvailable,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		DeletedAt:   entity.DeletedAt,
	}
}
