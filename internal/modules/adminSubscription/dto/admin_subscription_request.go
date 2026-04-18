package dto

type CreateSubscriptionPlanRequest struct {
	Name         string  `json:"name" binding:"required"`
	Description  string  `json:"description" binding:"required"`
	Price        float64 `json:"price" binding:"required"`
	BillingCycle string  `json:"billingCycle" binding:"required"` // contoh: monthly, yearly
	IsActive     bool    `json:"isActive"`
}

type UpdateSubscriptionPlanRequest struct {
	Name         *string  `json:"name"`
	Description  *string  `json:"description"`
	Price        *float64 `json:"price"`
	BillingCycle *string  `json:"billingCycle"`
	IsActive     *bool    `json:"isActive"`
}
