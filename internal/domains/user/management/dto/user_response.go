package dto

type UserResponse struct {
	ID       string `json:"id"`
	TenantID string `json:"tenantId"`
	Email    string `json:"email"`
	FullName string `json:"fullName"`
	Role     string `json:"role"`
	IsActive bool   `json:"isActive"`
}

type ListUsersResponse struct {
	Users []UserResponse `json:"users"`
}

type DetailUserResponse struct {
	User UserResponse `json:"user"`
}

type UpdateUserResponse struct {
	User UserResponse `json:"user"`
}

type ToggleActiveUserResponse struct {
	User UserResponse `json:"user"`
}

type DeleteUserResponse struct {
	Deleted bool `json:"deleted"`
}

type UserTenantResponse struct {
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
}

type AdminUserResponse struct {
	ID       string                `json:"id"`
	Email    string                `json:"email"`
	FullName string                `json:"fullName"`
	Role     string                `json:"role"`
	IsActive bool                  `json:"isActive"`
	Tenants  *[]UserTenantResponse `json:"tenants,omitempty"`
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

type ListAdminUsersResponse struct {
	Users      []AdminUserResponse `json:"users"`
	Pagination PaginationResponse  `json:"pagination"`
}

type AdminUserDetailResponse struct {
	User AdminUserResponse `json:"user"`
}
