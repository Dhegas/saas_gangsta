package usecase

import (
	"context"

	// Sesuaikan path import jika perlu
	"github.com/dhegas/saas_gangsta/internal/domains/subscription/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/subscription/dto"
)

type adminSubscriptionUsecase struct {
	repo domain.AdminSubscriptionRepository
}

// Constructor
func NewAdminSubscriptionUsecase(repo domain.AdminSubscriptionRepository) domain.AdminSubscriptionUsecase {
	return &adminSubscriptionUsecase{repo: repo}
}

func (u *adminSubscriptionUsecase) GetAllPlans(ctx context.Context) ([]dto.SubscriptionPlanResponse, error) {
	// 1. Ambil data asli dari database
	entities, err := u.repo.GetAllPlans(ctx)
	if err != nil {
		return nil, err
	}

	// 2. Terjemahkan Entity ke DTO
	var response []dto.SubscriptionPlanResponse
	for _, entity := range entities {
		response = append(response, dto.SubscriptionPlanResponse{
			ID:           entity.ID,
			Name:         entity.Name,
			Description:  entity.Description,
			Price:        entity.Price,
			BillingCycle: entity.BillingCycle,
			IsActive:     entity.IsActive,
		})
	}

	// Pastikan mengembalikan array kosong [] jika data belum ada (bukan null)
	if response == nil {
		response = []dto.SubscriptionPlanResponse{}
	}

	return response, nil
}

func (u *adminSubscriptionUsecase) CreatePlan(ctx context.Context, req dto.CreateSubscriptionPlanRequest) error {
	entity := &domain.SubscriptionPlanEntity{
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		BillingCycle: req.BillingCycle,
		IsActive:     req.IsActive,
	}
	return u.repo.CreatePlan(ctx, entity)
}

func (u *adminSubscriptionUsecase) UpdatePlan(ctx context.Context, id string, req dto.UpdateSubscriptionPlanRequest) error {
	updateData := make(map[string]interface{})

	if req.Name != nil {
		updateData["name"] = *req.Name
	}
	if req.Description != nil {
		updateData["description"] = *req.Description
	}
	if req.Price != nil {
		updateData["price"] = *req.Price
	}
	if req.BillingCycle != nil {
		updateData["billing_cycle"] = *req.BillingCycle
	}
	if req.IsActive != nil {
		updateData["is_active"] = *req.IsActive
	}

	return u.repo.UpdatePlan(ctx, id, updateData)
}
