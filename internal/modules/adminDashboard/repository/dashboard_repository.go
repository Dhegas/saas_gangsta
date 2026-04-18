package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
	// Sesuaikan path import jika berbeda
	"github.com/dhegas/saas_gangsta/internal/modules/adminDashboard/domain"
)

type adminDashboardRepository struct {
	db *gorm.DB
}

// Constructor
func NewAdminDashboardRepository(db *gorm.DB) domain.AdminDashboardRepository {
	return &adminDashboardRepository{db: db}
}

func (r *adminDashboardRepository) CountTotalTenants(ctx context.Context) (int, error) {
	var count int64
	// Menghitung jumlah baris di tabel tenants
	err := r.db.WithContext(ctx).Table("tenants").Count(&count).Error
	return int(count), err
}

func (r *adminDashboardRepository) CountActiveSubscriptions(ctx context.Context) (int, error) {
	var count int64
	// Menghitung langganan yang statusnya 'active'
	err := r.db.WithContext(ctx).Table("subscriptions").Where("status = ?", "active").Count(&count).Error
	return int(count), err
}

func (r *adminDashboardRepository) CalculateMonthlyRevenue(ctx context.Context) (float64, error) {
	var total float64

	// Mendapatkan waktu awal bulan ini (tanggal 1, jam 00:00)
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// Menjumlahkan kolom 'amount' di tabel payments yang statusnya 'paid' bulan ini
	err := r.db.WithContext(ctx).
		Table("payments").
		Where("status = ?", "paid").
		Where("paid_at >= ?", startOfMonth).
		Select("COALESCE(SUM(amount), 0)"). // COALESCE agar tidak error jika belum ada transaksi (null)
		Scan(&total).Error

	return total, err
}
