package domain

import (
	"context"

	// Sesuaikan dengan nama modul di go.mod kamu
	"github.com/dhegas/saas_gangsta/internal/domains/report/dto"
)

// AdminDashboardUsecase adalah kontrak untuk logika manajer bisnis
type AdminDashboardUsecase interface {
	GetStats(ctx context.Context) (*dto.DashboardStatsResponse, error)
}

// AdminDashboardRepository adalah kontrak untuk pekerja database
type AdminDashboardRepository interface {
	CountTotalTenants(ctx context.Context) (int, error)
	CountActiveSubscriptions(ctx context.Context) (int, error)
	CalculateMonthlyRevenue(ctx context.Context) (float64, error)
}
