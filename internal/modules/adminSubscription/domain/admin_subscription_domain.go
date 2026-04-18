package domain

import (
	"context"

	// Sesuaikan jika path foldermu berbeda
	"github.com/dhegas/saas_gangsta/internal/modules/adminSubscription/dto"
)

// SubscriptionPlanEntity merepresentasikan bentuk asli di tabel 'subscription_plans' Supabase
type SubscriptionPlanEntity struct {
	ID           string
	Name         string
	Description  string
	Price        float64
	BillingCycle string
	IsActive     bool
}

// AdminSubscriptionRepository adalah pekerja database
type AdminSubscriptionRepository interface {
	GetAllPlans(ctx context.Context) ([]SubscriptionPlanEntity, error)
}

// AdminSubscriptionUsecase adalah manajer logika bisnis
type AdminSubscriptionUsecase interface {
	GetAllPlans(ctx context.Context) ([]dto.SubscriptionPlanResponse, error)
}
