package dto

// UpdateTenantStatusRequest adalah payload untuk PATCH /api/v1/admin/tenants/:id/status
type UpdateTenantStatusRequest struct {
	// validator.v10 akan memastikan status tidak kosong dan hanya berisi nilai tertentu
	Status string `json:"status" binding:"required,oneof=active inactive suspended"`
}
