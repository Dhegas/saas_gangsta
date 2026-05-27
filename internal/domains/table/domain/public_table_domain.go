package domain

import (
	"context"

	"github.com/dhegas/saas_gangsta/internal/domains/table/dto"
)

type PublicTableUsecase interface {
	GetPublicTables(ctx context.Context, tenantID string) ([]dto.PublicTableResponse, error)
}

type PublicTableRepository interface {
	FindPublicTables(ctx context.Context, tenantID string) ([]dto.PublicTableResponse, error)
}
