package dto

// CreateTenantRequest adalah payload untuk POST /api/tenants
type CreateTenantRequest struct {
	Name   string `json:"name"   binding:"required,min=2,max=150"`
	Slug   string `json:"slug"   binding:"required,min=2,max=80,alphanum"`
	Status string `json:"status" binding:"omitempty,oneof=active inactive suspended"`
}

// UpdateTenantRequest adalah payload untuk PUT /api/tenants/:id
type UpdateTenantRequest struct {
	Name   string `json:"name"   binding:"omitempty,min=2,max=150"`
	Slug   string `json:"slug"   binding:"omitempty,min=2,max=80,alphanum"`
	Status string `json:"status" binding:"omitempty,oneof=active inactive suspended"`
}

// UpdateTenantStatusRequest adalah payload untuk PATCH /api/v1/admin/tenants/:id/status
type UpdateTenantStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active inactive suspended"`
}
