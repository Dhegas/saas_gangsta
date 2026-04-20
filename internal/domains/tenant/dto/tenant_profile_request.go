package dto

// CreateTenantProfileRequest adalah payload untuk membuat tenant profile.
type CreateTenantProfileRequest struct {
	Name        string `json:"name" binding:"required,max=120"`
	Description string `json:"description"`
	SortOrder   *int   `json:"sortOrder"`
	IsActive    *bool  `json:"isActive"`
}
