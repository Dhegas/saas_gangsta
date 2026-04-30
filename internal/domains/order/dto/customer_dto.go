package dto

import "time"

// CreateCustomerRequest payload untuk POST /api/orders/:id/customer
type CreateCustomerRequest struct {
	FullName    string `json:"full_name" binding:"required,min=2,max=150"`
	PhoneNumber string `json:"phone_number" binding:"omitempty,max=20"`
}

// UpdateCustomerRequest payload untuk PUT /api/orders/:id/customer
type UpdateCustomerRequest struct {
	FullName    string `json:"full_name" binding:"required,min=2,max=150"`
	PhoneNumber string `json:"phone_number" binding:"omitempty,max=20"`
}

// CustomerResponse response untuk data customer
type CustomerResponse struct {
	ID          string     `json:"id"`
	OrderID     string     `json:"order_id"`
	TenantID    string     `json:"tenant_id"`
	FullName    string     `json:"full_name"`
	PhoneNumber string     `json:"phone_number,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
