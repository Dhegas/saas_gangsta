package dto

type PublicTenantResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	LogoURL     string `json:"logoUrl"`
	BannerURL   string `json:"bannerUrl"`
	Address     string `json:"address"`
	OpenHours   string `json:"openHours"`
	IsOpen      bool   `json:"isOpen"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}

type ListPublicTenantsResponse struct {
	Data []PublicTenantResponse `json:"data"`
	Meta PaginationMeta         `json:"meta"`
}

type PublicTenantDetailResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	LogoURL     string `json:"logoUrl"`
	BannerURL   string `json:"bannerUrl"`
	Description string `json:"description"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phoneNumber"`
	OpenHours   string `json:"openHours"`
}
