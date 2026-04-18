package dto

// SubscriptionPlanResponse adalah bentuk data paket yang akan dikirim ke Flutter
type SubscriptionPlanResponse struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	BillingCycle string  `json:"billingCycle"`
	IsActive     bool    `json:"isActive"`
}
