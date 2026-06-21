package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dhegas/saas_gangsta/internal/common/cache"
	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant/dto"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant/repository"
)

type adminTenantUsecase struct {
	repo  domain.AdminTenantRepository
	cache *cache.LocalCache
}

func NewAdminTenantUsecase(repo domain.AdminTenantRepository, cache *cache.LocalCache) domain.AdminTenantUsecase {
	return &adminTenantUsecase{
		repo:  repo,
		cache: cache,
	}
}

func (u *adminTenantUsecase) CreateAdminTenant(ctx context.Context, req dto.CreateAdminTenantRequest) (*dto.CreateAdminTenantResponse, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, apperrors.New("VALIDATION_ERROR", "Nama tenant wajib diisi", http.StatusBadRequest)
	}

	tenant, err := u.repo.CreateTenantForAdmin(ctx, domain.CreateAdminTenantInput{
		UserID:      req.UserID,
		Name:        name,
		Status:      req.Status,
		Description: strings.TrimSpace(req.Description),
		Address:     strings.TrimSpace(req.Address),
		PhoneNumber: strings.TrimSpace(req.PhoneNumber),
		OpenHours:   strings.TrimSpace(req.OpenHours),
		LogoURL:     req.LogoURL,
		BannerURL:   req.BannerURL,
	})
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			return nil, apperrors.New("NOT_FOUND", "Partner user tidak ditemukan atau tidak aktif", http.StatusNotFound)
		case errors.Is(err, repository.ErrUserNotPartner):
			return nil, apperrors.New("FORBIDDEN", "Tenant hanya dapat diasosiasikan kepada user dengan role PARTNER", http.StatusForbidden)
		default:
			return nil, apperrors.New("INTERNAL_ERROR", "Gagal membuat tenant oleh admin", http.StatusInternalServerError)
		}
	}

	if u.cache != nil {
		u.cache.DeleteByPrefix("admin:tenants:")
	}

	return &dto.CreateAdminTenantResponse{
		Tenant: dto.AdminTenantResponse{
			ID:          tenant.ID,
			Name:        tenant.Name,
			Slug:        tenant.Slug,
			Status:      tenant.Status,
			Description: tenant.Description,
			Address:     tenant.Address,
			PhoneNumber: tenant.PhoneNumber,
			OpenHours:   tenant.OpenHours,
			LogoURL:     tenant.LogoURL,
			BannerURL:   tenant.BannerURL,
			UserID:      tenant.UserID,
			User: func() *dto.TenantUserResponse {
				if tenant.UserID == "" {
					return nil
				}
				return &dto.TenantUserResponse{
					ID:       tenant.UserID,
					FullName: tenant.OwnerName,
				}
			}(),
		},
	}, nil
}

func (u *adminTenantUsecase) ListAllTenants(ctx context.Context, req dto.ListAllTenantsRequest) (*dto.ListAllTenantsResponse, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10 // default as requested
	}

	cacheKey := fmt.Sprintf("admin:tenants:page:%d:limit:%d", page, limit)
	if u.cache != nil {
		if cached, found := u.cache.Get(cacheKey); found {
			if cachedResponse, ok := cached.(*dto.ListAllTenantsResponse); ok {
				return cachedResponse, nil
			}
		}
	}

	offset := (page - 1) * limit

	tenants, totalItems, err := u.repo.ListAllTenants(ctx, limit, offset)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil daftar tenant oleh admin", http.StatusInternalServerError)
	}

	items := make([]dto.AdminTenantResponse, 0, len(tenants))
	for _, tenant := range tenants {
		items = append(items, dto.AdminTenantResponse{
			ID:          tenant.ID,
			Name:        tenant.Name,
			Slug:        tenant.Slug,
			Status:      tenant.Status,
			Description: tenant.Description,
			Address:     tenant.Address,
			PhoneNumber: tenant.PhoneNumber,
			OpenHours:   tenant.OpenHours,
			LogoURL:     tenant.LogoURL,
			BannerURL:   tenant.BannerURL,
			UserID:      tenant.UserID,
			User: func() *dto.TenantUserResponse {
				if tenant.UserID == "" {
					return nil
				}
				return &dto.TenantUserResponse{
					ID:       tenant.UserID,
					FullName: tenant.OwnerName,
				}
			}(),
		})
	}

	totalPages := 0
	if totalItems > 0 {
		totalPages = int((totalItems + int64(limit) - 1) / int64(limit))
	}

	res := &dto.ListAllTenantsResponse{
		Tenants: items,
		Pagination: dto.PaginationResponse{
			Page:       page,
			Limit:      limit,
			TotalItems: totalItems,
			TotalPages: totalPages,
		},
	}

	if u.cache != nil {
		u.cache.Set(cacheKey, res, 15*time.Minute)
	}

	return res, nil
}

func (u *adminTenantUsecase) SoftDeleteTenant(ctx context.Context, tenantID string) error {
	tenantID = strings.TrimSpace(tenantID)
	if tenantID == "" {
		return apperrors.New("VALIDATION_ERROR", "ID tenant wajib diisi", http.StatusBadRequest)
	}

	err := u.repo.SoftDeleteTenant(ctx, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrTenantNotFound) {
			return apperrors.New("NOT_FOUND", "Tenant ID tidak ditemukan", http.StatusNotFound)
		}
		return apperrors.New("INTERNAL_ERROR", "Gagal menghapus tenant oleh admin", http.StatusInternalServerError)
	}

	if u.cache != nil {
		u.cache.DeleteByPrefix("admin:tenants:")
	}

	return nil
}

func (u *adminTenantUsecase) GetTenantsByUserID(ctx context.Context, userID string) (*dto.GetTenantsByUserIDResponse, error) {
	userID = strings.TrimSpace(userID)
	if userID == "" {
		return nil, apperrors.New("VALIDATION_ERROR", "User ID partner wajib diisi", http.StatusBadRequest)
	}

	tenants, err := u.repo.GetTenantsByUserID(ctx, userID)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil data tenant berdasarkan User ID", http.StatusInternalServerError)
	}

	items := make([]dto.AdminTenantResponse, 0, len(tenants))
	for _, tenant := range tenants {
		items = append(items, dto.AdminTenantResponse{
			ID:          tenant.ID,
			Name:        tenant.Name,
			Slug:        tenant.Slug,
			Status:      tenant.Status,
			Description: tenant.Description,
			Address:     tenant.Address,
			PhoneNumber: tenant.PhoneNumber,
			OpenHours:   tenant.OpenHours,
			LogoURL:     tenant.LogoURL,
			BannerURL:   tenant.BannerURL,
			UserID:      tenant.UserID,
			User: func() *dto.TenantUserResponse {
				if tenant.UserID == "" {
					return nil
				}
				return &dto.TenantUserResponse{
					ID:       tenant.UserID,
					FullName: tenant.OwnerName,
				}
			}(),
		})
	}

	return &dto.GetTenantsByUserIDResponse{
		Tenants: items,
	}, nil
}

func (u *adminTenantUsecase) GetTenantByID(ctx context.Context, tenantID string) (*dto.AdminTenantResponse, error) {
	tenantID = strings.TrimSpace(tenantID)
	if tenantID == "" {
		return nil, apperrors.New("VALIDATION_ERROR", "ID Tenant wajib diisi", http.StatusBadRequest)
	}

	tenant, err := u.repo.GetTenantByID(ctx, tenantID)
	if err != nil {
		if errors.Is(err, repository.ErrTenantNotFound) {
			return nil, apperrors.New("NOT_FOUND", "Tenant ID tidak ditemukan", http.StatusNotFound)
		}
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil detail tenant oleh admin", http.StatusInternalServerError)
	}

	return &dto.AdminTenantResponse{
		ID:          tenant.ID,
		Name:        tenant.Name,
		Slug:        tenant.Slug,
		Status:      tenant.Status,
		Description: tenant.Description,
		Address:     tenant.Address,
		PhoneNumber: tenant.PhoneNumber,
		OpenHours:   tenant.OpenHours,
		LogoURL:     tenant.LogoURL,
		BannerURL:   tenant.BannerURL,
		UserID:      tenant.UserID,
		User: func() *dto.TenantUserResponse {
			if tenant.UserID == "" {
				return nil
			}
			return &dto.TenantUserResponse{
				ID:       tenant.UserID,
				FullName: tenant.OwnerName,
			}
		}(),
	}, nil
}
