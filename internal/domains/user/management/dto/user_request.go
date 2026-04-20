package dto

type UserIDParam struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type UpdateUserRequest struct {
	Email    *string `json:"email" binding:"omitempty,email"`
	FullName *string `json:"fullName" binding:"omitempty,min=2,max=150"`
	Role     *string `json:"role" binding:"omitempty,oneof=BASIC MITRA ADMIN"`
}
