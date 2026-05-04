package dto

type CreatePartnerTenantRequest struct {
	Name string `json:"name" binding:"required,min=2,max=150"`
}
