package repository

import (
	"context"

	"github.com/dhegas/saas_gangsta/internal/domains/report/domain"
	"gorm.io/gorm"
)

type partnerReportRepository struct {
	db *gorm.DB
}

func NewPartnerReportRepository(db *gorm.DB) domain.PartnerReportRepository {
	return &partnerReportRepository{db: db}
}

// FetchRevenue menghitung total pendapatan dan jumlah order COMPLETED dalam rentang tanggal.
func (r *partnerReportRepository) FetchRevenue(ctx context.Context, tenantID, from, to string) (float64, int, error) {
	var result struct {
		TotalRevenue float64
		TotalOrders  int
	}
	err := r.db.WithContext(ctx).
		Table("orders").
		Select("COALESCE(SUM(total_price), 0) AS total_revenue, COUNT(*) AS total_orders").
		Where("tenant_id = ? AND status = 'COMPLETED' AND deleted_at IS NULL", tenantID).
		Where("created_at::date BETWEEN ? AND ?", from, to).
		Scan(&result).Error
	return result.TotalRevenue, result.TotalOrders, err
}

// FetchTopMenus mengambil menu terlaris berdasarkan total qty terjual.
func (r *partnerReportRepository) FetchTopMenus(ctx context.Context, tenantID, from, to string, limit int) ([]domain.TopMenuRow, error) {
	query := r.db.WithContext(ctx).
		Table("order_items oi").
		Select("oi.menu_id, oi.menu_name, SUM(oi.quantity) AS total_qty, SUM(oi.subtotal) AS total_sold").
		Joins("JOIN orders o ON o.id = oi.order_id").
		Where("o.tenant_id = ? AND o.status = 'COMPLETED' AND o.deleted_at IS NULL AND oi.deleted_at IS NULL", tenantID).
		Group("oi.menu_id, oi.menu_name").
		Order("total_qty DESC")

	if from != "" && to != "" {
		query = query.Where("o.created_at::date BETWEEN ? AND ?", from, to)
	}
	if limit > 0 {
		query = query.Limit(limit)
	} else {
		query = query.Limit(10) // default
	}

	var rows []domain.TopMenuRow
	err := query.Scan(&rows).Error
	return rows, err
}

// FetchOrdersByTable mengambil meja dengan jumlah order terbanyak.
func (r *partnerReportRepository) FetchOrdersByTable(ctx context.Context, tenantID, from, to string, limit int) ([]domain.OrdersByTableRow, error) {
	query := r.db.WithContext(ctx).
		Table("orders o").
		Select(`o.dining_tables_id AS table_id,
			COALESCE(dt.table_number, 'N/A') AS table_number,
			COUNT(o.id) AS total_orders,
			COALESCE(SUM(o.total_price), 0) AS total_revenue`).
		Joins("LEFT JOIN dining_tables dt ON dt.id = o.dining_tables_id").
		Where("o.tenant_id = ? AND o.status = 'COMPLETED' AND o.deleted_at IS NULL AND o.dining_tables_id IS NOT NULL", tenantID).
		Group("o.dining_tables_id, dt.table_number").
		Order("total_orders DESC")

	if from != "" && to != "" {
		query = query.Where("o.created_at::date BETWEEN ? AND ?", from, to)
	}
	if limit > 0 {
		query = query.Limit(limit)
	} else {
		query = query.Limit(10) // default
	}

	var rows []domain.OrdersByTableRow
	err := query.Scan(&rows).Error
	return rows, err
}

// FetchDailySummary mengambil ringkasan order dan revenue per hari.
func (r *partnerReportRepository) FetchDailySummary(ctx context.Context, tenantID, from, to string) ([]domain.DailySummaryRow, error) {
	var rows []domain.DailySummaryRow
	err := r.db.WithContext(ctx).
		Table("orders").
		Select(`created_at::date AS date,
			COUNT(*) AS total_orders,
			COALESCE(SUM(total_price), 0) AS total_revenue`).
		Where("tenant_id = ? AND status = 'COMPLETED' AND deleted_at IS NULL", tenantID).
		Where("created_at::date BETWEEN ? AND ?", from, to).
		Group("created_at::date").
		Order("date ASC").
		Scan(&rows).Error
	return rows, err
}
