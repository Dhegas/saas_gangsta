package usecase

import (
	"context"
	"net/http"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/domains/menu/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/menu/dto"
)

type publicMenuUsecase struct {
	repo domain.PublicMenuRepository
}

func NewPublicMenuUsecase(repo domain.PublicMenuRepository) domain.PublicMenuUsecase {
	return &publicMenuUsecase{repo: repo}
}

func (u *publicMenuUsecase) GetPublicMenus(ctx context.Context, tenantID string, categoryID string, search string, isAvailable *bool) ([]dto.MenuResponse, error) {
	menus, err := u.repo.FindPublicMenus(ctx, tenantID, categoryID, search, isAvailable)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil daftar menu publik", http.StatusInternalServerError, err)
	}

	result := make([]dto.MenuResponse, 0, len(menus))
	for _, m := range menus {
		result = append(result, toMenuResponse(&m))
	}

	return result, nil
}
