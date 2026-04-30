package usecase

import (
	"context"
	"fmt"
	"net/http"
	"time"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/domains/report/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/report/dto"
)

type merchantReportUsecase struct {
	repo domain.MerchantReportRepository
}

func NewMerchantReportUsecase(repo domain.MerchantReportRepository) domain.MerchantReportUsecase {
	return &merchantReportUsecase{repo: repo}
}

func (u *merchantReportUsecase) GetRevenue(ctx context.Context, tenantID string, params dto.RevenueFilterParams) (*dto.RevenueResponse, error) {
	if err := validateDateRange(params.From, params.To); err != nil {
		return nil, apperrors.New("BAD_REQUEST", err.Error(), http.StatusBadRequest, nil)
	}

	totalRevenue, totalOrders, err := u.repo.FetchRevenue(ctx, tenantID, params.From, params.To)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil data revenue", http.StatusInternalServerError, err)
	}

	return &dto.RevenueResponse{
		From:         params.From,
		To:           params.To,
		TotalRevenue: totalRevenue,
		TotalOrders:  totalOrders,
		GeneratedAt:  time.Now().UTC(),
	}, nil
}

func (u *merchantReportUsecase) GetTopMenus(ctx context.Context, tenantID string, params dto.TopMenusFilterParams) (*dto.TopMenusResponse, error) {
	rows, err := u.repo.FetchTopMenus(ctx, tenantID, params.From, params.To, params.Limit)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil data menu terlaris", http.StatusInternalServerError, err)
	}

	entries := make([]dto.TopMenuEntry, 0, len(rows))
	for i, r := range rows {
		entries = append(entries, dto.TopMenuEntry{
			Rank:      i + 1,
			MenuID:    r.MenuID,
			MenuName:  r.MenuName,
			TotalQty:  r.TotalQty,
			TotalSold: r.TotalSold,
		})
	}

	return &dto.TopMenusResponse{
		From:  params.From,
		To:    params.To,
		Menus: entries,
	}, nil
}

func (u *merchantReportUsecase) GetOrdersByTable(ctx context.Context, tenantID string, params dto.OrdersByTableFilterParams) (*dto.OrdersByTableResponse, error) {
	rows, err := u.repo.FetchOrdersByTable(ctx, tenantID, params.From, params.To, params.Limit)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil data order per meja", http.StatusInternalServerError, err)
	}

	entries := make([]dto.OrdersByTableEntry, 0, len(rows))
	for i, r := range rows {
		entries = append(entries, dto.OrdersByTableEntry{
			Rank:         i + 1,
			TableID:      r.TableID,
			TableNumber:  r.TableNumber,
			TotalOrders:  r.TotalOrders,
			TotalRevenue: r.TotalRevenue,
		})
	}

	return &dto.OrdersByTableResponse{
		From:   params.From,
		To:     params.To,
		Tables: entries,
	}, nil
}

func (u *merchantReportUsecase) GetDailySummary(ctx context.Context, tenantID string, params dto.DailySummaryFilterParams) (*dto.DailySummaryResponse, error) {
	// Default: 7 hari ke belakang jika tidak ada parameter
	from := params.From
	to := params.To
	if from == "" {
		from = time.Now().AddDate(0, 0, -6).Format("2006-01-02")
	}
	if to == "" {
		to = time.Now().Format("2006-01-02")
	}

	if err := validateDateRange(from, to); err != nil {
		return nil, apperrors.New("BAD_REQUEST", err.Error(), http.StatusBadRequest, nil)
	}

	rows, err := u.repo.FetchDailySummary(ctx, tenantID, from, to)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil ringkasan harian", http.StatusInternalServerError, err)
	}

	entries := make([]dto.DailySummaryEntry, 0, len(rows))
	for _, r := range rows {
		avg := 0.0
		if r.TotalOrders > 0 {
			avg = r.TotalRevenue / float64(r.TotalOrders)
		}
		entries = append(entries, dto.DailySummaryEntry{
			Date:          r.Date,
			TotalOrders:   r.TotalOrders,
			TotalRevenue:  r.TotalRevenue,
			AvgOrderValue: avg,
		})
	}

	return &dto.DailySummaryResponse{
		From:    from,
		To:      to,
		Summary: entries,
	}, nil
}

// validateDateRange memastikan format tanggal valid dan from <= to
func validateDateRange(from, to string) error {
	layout := "2006-01-02"
	f, err := time.Parse(layout, from)
	if err != nil {
		return fmt.Errorf("format tanggal 'from' tidak valid, gunakan YYYY-MM-DD")
	}
	t, err := time.Parse(layout, to)
	if err != nil {
		return fmt.Errorf("format tanggal 'to' tidak valid, gunakan YYYY-MM-DD")
	}
	if f.After(t) {
		return fmt.Errorf("tanggal 'from' tidak boleh lebih besar dari 'to'")
	}
	return nil
}
