package dto

type CreatePartnerTenantRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	OpenHours   string `json:"open_hours"`
	LogoURL     string `json:"logo_url"`
	BannerURL   string `json:"banner_url"`
}

type UpdatePartnerTenantRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}
