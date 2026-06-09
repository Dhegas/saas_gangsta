package usecase

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dhegas/saas_gangsta/internal/common/cache"
	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/domains/menu/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/menu/dto"
)

type publicMenuUsecase struct {
	repo  domain.PublicMenuRepository
	cache *cache.LocalCache
}

func NewPublicMenuUsecase(repo domain.PublicMenuRepository, cache *cache.LocalCache) domain.PublicMenuUsecase {
	return &publicMenuUsecase{
		repo:  repo,
		cache: cache,
	}
}

func (u *publicMenuUsecase) GetPublicMenus(ctx context.Context, tenantID string, categoryID string, search string, isAvailable *bool) ([]dto.MenuResponse, error) {
	availStr := "nil"
	if isAvailable != nil {
		availStr = fmt.Sprintf("%v", *isAvailable)
	}
	cacheKey := fmt.Sprintf("public:menus:tenant:%s:cat:%s:search:%s:avail:%s", tenantID, categoryID, search, availStr)

	if cached, found := u.cache.Get(cacheKey); found {
		if cachedMenus, ok := cached.([]dto.MenuResponse); ok {
			return cachedMenus, nil
		}
	}

	menus, err := u.repo.FindPublicMenus(ctx, tenantID, categoryID, search, isAvailable)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil daftar menu publik", http.StatusInternalServerError)
	}

	result := make([]dto.MenuResponse, 0, len(menus))
	for _, m := range menus {
		result = append(result, toMenuResponse(&m))
	}

	u.cache.Set(cacheKey, result, 5*time.Minute)

	return result, nil
}
