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

type AdminUserResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"fullName"`
	Role     string `json:"role"`
	IsActive bool   `json:"isActive"`
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
