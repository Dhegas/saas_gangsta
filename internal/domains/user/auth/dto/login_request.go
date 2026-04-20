package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type SubscribeRequest struct {
	PlanID string `json:"planId" binding:"required,uuid"`
}

type CreateMerchantTenantRequest struct {
	Name string `json:"name" binding:"required,min=2,max=150"`
}
