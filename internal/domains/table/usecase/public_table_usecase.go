package usecase

import (
	"context"
	"net/http"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/domains/table/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/table/dto"
)

type publicTableUsecase struct {
	repo domain.PublicTableRepository
}

func NewPublicTableUsecase(repo domain.PublicTableRepository) domain.PublicTableUsecase {
	return &publicTableUsecase{repo: repo}
}

func (u *publicTableUsecase) GetPublicTables(ctx context.Context, tenantID string) ([]dto.PublicTableResponse, error) {
	tables, err := u.repo.FindPublicTables(ctx, tenantID)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil data meja publik", http.StatusInternalServerError, err)
	}
	return tables, nil
}
