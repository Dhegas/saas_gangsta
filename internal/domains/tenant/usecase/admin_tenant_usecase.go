package usecase

import (
	"context"

	"github.com/dhegas/saas_gangsta/internal/domains/tenant/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant/dto"
)

type adminTenantUsecase struct {
	repo domain.AdminTenantRepository
}

// Constructor
func NewAdminTenantUsecase(repo domain.AdminTenantRepository) domain.AdminTenantUsecase {
	return &adminTenantUsecase{repo: repo}
}

func (u *adminTenantUsecase) GetAllTenants(ctx context.Context) ([]dto.TenantResponse, error) {
	// 1. Ambil data dari database melalui repository
	entities, err := u.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// 2. Mapping dari Entity ke DTO (camelCase untuk JSON frontend)
	var responses []dto.TenantResponse
	for _, e := range entities {
		responses = append(responses, dto.TenantResponse{
			ID:        e.ID,
			Name:      e.Name,
			Slug:      e.Slug,
			Status:    e.Status,
			CreatedAt: e.CreatedAt,
		})
	}

	return responses, nil
}

func (u *adminTenantUsecase) UpdateTenantStatus(ctx context.Context, tenantID string, status string) error {
	// Di sini kamu bisa tambahkan validasi, misal: status hanya boleh "active" atau "suspended"
	return u.repo.UpdateStatus(ctx, tenantID, status)
}
