package dto

type PartnerTenantResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	Status  string `json:"status"`
	IsOwner bool   `json:"isOwner"`
}

type CreatePartnerTenantResponse struct {
	Tenant PartnerTenantResponse `json:"tenant"`
}

type ListPartnerTenantsResponse struct {
	Tenants []PartnerTenantResponse `json:"tenants"`
}
