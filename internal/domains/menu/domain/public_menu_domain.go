package domain

import (
	"context"

	"github.com/dhegas/saas_gangsta/internal/domains/menu/dto"
)

type PublicMenuUsecase interface {
	GetPublicMenus(ctx context.Context, tenantID string, categoryID string, search string, isAvailable *bool) ([]dto.MenuResponse, error)
}

type PublicMenuRepository interface {
	FindPublicMenus(ctx context.Context, tenantID string, categoryID string, search string, isAvailable *bool) ([]MenuEntity, error)
}
