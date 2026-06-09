package usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/dhegas/saas_gangsta/internal/common/cache"
	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant/dto"
	"gorm.io/gorm"
)

type publicTenantUsecase struct {
	repo  domain.PublicTenantRepository
	cache *cache.LocalCache
}

func NewPublicTenantUsecase(repo domain.PublicTenantRepository, cache *cache.LocalCache) domain.PublicTenantUsecase {
	return &publicTenantUsecase{
		repo:  repo,
		cache: cache,
	}
}

func (u *publicTenantUsecase) ListPublicTenants(ctx context.Context, req dto.ListPublicTenantsRequest) (*dto.ListPublicTenantsResponse, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	cacheKey := fmt.Sprintf("public:tenants:search:%s:page:%d:limit:%d", req.Search, page, limit)
	if cached, found := u.cache.Get(cacheKey); found {
		if cachedResponse, ok := cached.(*dto.ListPublicTenantsResponse); ok {
			return cachedResponse, nil
		}
	}

	offset := (page - 1) * limit

	tenants, totalItems, err := u.repo.ListPublicTenants(ctx, req.Search, limit, offset)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil daftar tenant", http.StatusInternalServerError)
	}

	items := make([]dto.PublicTenantResponse, 0, len(tenants))
	for _, t := range tenants {
		items = append(items, dto.PublicTenantResponse{
			ID:        t.ID,
			Name:      t.Name,
			Slug:      t.Slug,
			LogoURL:   t.LogoURL,
			BannerURL: t.BannerURL,
			Address:   t.Address,
			OpenHours: t.OpenHours,
			IsOpen:    t.Status == "active",
		})
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	res := &dto.ListPublicTenantsResponse{
		Data: items,
		Meta: dto.PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      totalItems,
			TotalPages: totalPages,
		},
	}

	u.cache.Set(cacheKey, res, 15*time.Minute)

	return res, nil
}

func (u *publicTenantUsecase) GetPublicTenantBySlug(ctx context.Context, slug string) (*dto.PublicTenantDetailResponse, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return nil, apperrors.New("VALIDATION_ERROR", "Slug tenant wajib diisi", http.StatusBadRequest)
	}

	cacheKey := fmt.Sprintf("public:tenant:slug:%s", slug)
	if cached, found := u.cache.Get(cacheKey); found {
		if cachedTenant, ok := cached.(*dto.PublicTenantDetailResponse); ok {
			return cachedTenant, nil
		}
	}

	tenant, err := u.repo.FindTenantBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New("NOT_FOUND", "Tenant tidak ditemukan", http.StatusNotFound)
		}
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil detail tenant", http.StatusInternalServerError)
	}

	res := &dto.PublicTenantDetailResponse{
		ID:          tenant.ID,
		Name:        tenant.Name,
		Slug:        tenant.Slug,
		LogoURL:     tenant.LogoURL,
		BannerURL:   tenant.BannerURL,
		Description: tenant.Description,
		Address:     tenant.Address,
		PhoneNumber: tenant.PhoneNumber,
		OpenHours:   tenant.OpenHours,
	}

	u.cache.Set(cacheKey, res, 15*time.Minute)

	return res, nil
}
