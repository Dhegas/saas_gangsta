package repository

import (
	"context"

	"gorm.io/gorm"

	// Sesuaikan path import jika perlu
	"github.com/dhegas/saas_gangsta/internal/modules/adminSubscription/domain"
)

type adminSubscriptionRepository struct {
	db *gorm.DB
}

// Constructor
func NewAdminSubscriptionRepository(db *gorm.DB) domain.AdminSubscriptionRepository {
	return &adminSubscriptionRepository{db: db}
}

func (r *adminSubscriptionRepository) GetAllPlans(ctx context.Context) ([]domain.SubscriptionPlanEntity, error) {
	var plans []domain.SubscriptionPlanEntity

	// Kita paksa GORM untuk langsung membaca tabel "subscription_plans"
	err := r.db.WithContext(ctx).Table("subscription_plans").Find(&plans).Error
	if err != nil {
		return nil, err
	}

	return plans, nil
}
