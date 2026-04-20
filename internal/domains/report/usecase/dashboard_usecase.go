package usecase

import (
	"context"

	// Sesuaikan path import jika berbeda
	"github.com/dhegas/saas_gangsta/internal/domains/report/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/report/dto"
)

type adminDashboardUsecase struct {
	repo domain.AdminDashboardRepository
}

// Constructor
func NewAdminDashboardUsecase(repo domain.AdminDashboardRepository) domain.AdminDashboardUsecase {
	return &adminDashboardUsecase{repo: repo}
}

func (u *adminDashboardUsecase) GetStats(ctx context.Context) (*dto.DashboardStatsResponse, error) {
	// 1. Ambil total tenant
	totalTenants, err := u.repo.CountTotalTenants(ctx)
	if err != nil {
		return nil, err
	}

	// 2. Ambil langganan aktif
	activeSubs, err := u.repo.CountActiveSubscriptions(ctx)
	if err != nil {
		return nil, err
	}

	// 3. Ambil total pendapatan (revenue) bulan ini
	revenue, err := u.repo.CalculateMonthlyRevenue(ctx)
	if err != nil {
		return nil, err
	}

	// 4. Bungkus semua data ke dalam DTO
	return &dto.DashboardStatsResponse{
		TotalTenants:        totalTenants,
		ActiveSubscriptions: activeSubs,
		MonthlyRevenue:      revenue,
	}, nil
}
