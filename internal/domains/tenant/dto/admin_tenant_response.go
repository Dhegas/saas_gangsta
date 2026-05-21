package dto

type TenantUserResponse struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
}

type AdminTenantResponse struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Slug        string              `json:"slug"`
	Status      string              `json:"status"`
	Description string              `json:"description"`
	Address     string              `json:"address"`
	PhoneNumber string              `json:"phone_number"`
	OpenHours   string              `json:"open_hours"`
	LogoURL     string              `json:"logo_url"`
	BannerURL   string              `json:"banner_url"`
	UserID      string              `json:"user_id"` // Target partner owner's ID
	User        *TenantUserResponse `json:"user"`
}

type CreateAdminTenantResponse struct {
	Tenant AdminTenantResponse `json:"tenant"`
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

type ListAllTenantsResponse struct {
	Tenants    []AdminTenantResponse `json:"tenants"`
	Pagination PaginationResponse    `json:"pagination"`
}

type GetTenantsByUserIDResponse struct {
	Tenants []AdminTenantResponse `json:"tenants"`
}
