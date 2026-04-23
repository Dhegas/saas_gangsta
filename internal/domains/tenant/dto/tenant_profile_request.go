package dto

// CreateTenantProfileRequest adalah payload untuk membuat tenant profile.
type CreateTenantProfileRequest struct {
	Name        string `json:"name" binding:"required,max=120"`
	Description string `json:"description"`
	SortOrder   *int   `json:"sortOrder"`
	IsActive    *bool  `json:"isActive"`
}

// UpdateTenantProfileRequest adalah payload untuk mengupdate tenant profile.
// Semua field bersifat opsional (partial update).
type UpdateTenantProfileRequest struct {
	Name        *string `json:"name" binding:"omitempty,max=120"`
	Description *string `json:"description"`
	SortOrder   *int    `json:"sortOrder"`
	IsActive    *bool   `json:"isActive"`
}

