package domain

import (
	"context"

	"github.com/dhegas/saas_gangsta/internal/domains/report/dto"
)

// MerchantReportUsecase kontrak logika bisnis laporan merchant
type MerchantReportUsecase interface {
	GetRevenue(ctx context.Context, tenantID string, params dto.RevenueFilterParams) (*dto.RevenueResponse, error)
	GetTopMenus(ctx context.Context, tenantID string, params dto.TopMenusFilterParams) (*dto.TopMenusResponse, error)
	GetOrdersByTable(ctx context.Context, tenantID string, params dto.OrdersByTableFilterParams) (*dto.OrdersByTableResponse, error)
	GetDailySummary(ctx context.Context, tenantID string, params dto.DailySummaryFilterParams) (*dto.DailySummaryResponse, error)
}

// MerchantReportRepository kontrak interaksi database untuk laporan
type MerchantReportRepository interface {
	FetchRevenue(ctx context.Context, tenantID, from, to string) (float64, int, error)
	FetchTopMenus(ctx context.Context, tenantID, from, to string, limit int) ([]TopMenuRow, error)
	FetchOrdersByTable(ctx context.Context, tenantID, from, to string, limit int) ([]OrdersByTableRow, error)
	FetchDailySummary(ctx context.Context, tenantID, from, to string) ([]DailySummaryRow, error)
}

// ---- Raw row types returned from repository ----

// TopMenuRow hasil query terlaris dari DB
type TopMenuRow struct {
	MenuID   string
	MenuName string
	TotalQty int
	TotalSold float64
}

// OrdersByTableRow hasil query order per meja dari DB
type OrdersByTableRow struct {
	TableID      string
	TableNumber  string
	TotalOrders  int
	TotalRevenue float64
}

// DailySummaryRow hasil query ringkasan harian dari DB
type DailySummaryRow struct {
	Date         string
	TotalOrders  int
	TotalRevenue float64
}
