package domain

import (
	"context"

	"github.com/dhegas/saas_gangsta/internal/modules/adminSubscription/dto"
)

type SubscriptionPlanEntity struct {
	ID           string `gorm:"primaryKey;default:gen_random_uuid()"`
	Name         string
	Description  string
	Price        float64
	BillingCycle string
	IsActive     bool
}

type AdminSubscriptionRepository interface {
	GetAllPlans(ctx context.Context) ([]SubscriptionPlanEntity, error)
	CreatePlan(ctx context.Context, plan *SubscriptionPlanEntity) error
	UpdatePlan(ctx context.Context, id string, updateData map[string]interface{}) error
}

type AdminSubscriptionUsecase interface {
	GetAllPlans(ctx context.Context) ([]dto.SubscriptionPlanResponse, error)
	CreatePlan(ctx context.Context, req dto.CreateSubscriptionPlanRequest) error
	UpdatePlan(ctx context.Context, id string, req dto.UpdateSubscriptionPlanRequest) error
}
