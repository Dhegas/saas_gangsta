package dto

type PartnerTenantResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Status      string `json:"status"`
	Description string `json:"description"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	OpenHours   string `json:"open_hours"`
	LogoURL     string `json:"logo_url"`
	BannerURL   string `json:"banner_url"`
	IsOwner     bool   `json:"isOwner"`
}

type CreatePartnerTenantResponse struct {
	Tenant PartnerTenantResponse `json:"tenant"`
}

type ListPartnerTenantsResponse struct {
	Tenants []PartnerTenantResponse `json:"tenants"`
}
