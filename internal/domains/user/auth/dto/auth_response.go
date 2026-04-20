package dto

type UserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	TenantID string `json:"tenantId,omitempty"`
}

type LoginResponse struct {
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
	User         UserResponse `json:"user"`
}

type RegisterResponse struct {
	User UserResponse `json:"user"`
}

type MeResponse struct {
	User UserResponse `json:"user"`
}

type SubscribeResponse struct {
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
	User         UserResponse `json:"user"`
}

type MerchantTenantResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	Status  string `json:"status"`
	IsOwner bool   `json:"isOwner"`
}

type CreateMerchantTenantResponse struct {
	Tenant MerchantTenantResponse `json:"tenant"`
}

type ListMerchantTenantsResponse struct {
	Tenants []MerchantTenantResponse `json:"tenants"`
}
