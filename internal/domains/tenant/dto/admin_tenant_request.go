package dto

type CreateAdminTenantRequest struct {
	UserID      string `json:"user_id" binding:"required,uuid"` // ID of the partner who will own this tenant
	Name        string `json:"name" binding:"required"`         // Tenant name
	Status      string `json:"status" binding:"omitempty,oneof=active inactive suspended"` // Optional: default "active"
	Description string `json:"description"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	OpenHours   string `json:"open_hours"`
	LogoURL     string `json:"logo_url"`
	BannerURL   string `json:"banner_url"`
}

type ListAllTenantsRequest struct {
	Page  int `form:"page" binding:"omitempty,min=1"`
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
}
