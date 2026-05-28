package dto

// CreateCustomerOrderItemRequest payload untuk detail item pada pemesanan publik
type CreateCustomerOrderItemRequest struct {
	MenuID   string `json:"menuId" binding:"required,uuid"`
	Quantity int    `json:"quantity" binding:"required,min=1"`
	Notes    string `json:"notes" binding:"omitempty,max=255"`
}

// CreateCustomerDetailsRequest payload detail customer dalam pemesanan publik
type CreateCustomerDetailsRequest struct {
	FullName    string `json:"fullName" binding:"required,min=2,max=150"`
	Email       string `json:"email" binding:"omitempty,email"`
	Password    string `json:"password" binding:"omitempty,min=6"`
	PhoneNumber string `json:"phoneNumber" binding:"omitempty,max=20"`
}

// CreateCustomerOrderRequest payload utama untuk POST /api/v1/public/tenant/:tenantSlug/orders
type CreateCustomerOrderRequest struct {
	DiningTableID string                           `json:"diningTableId" binding:"required,uuid"`
	Items         []CreateCustomerOrderItemRequest `json:"items" binding:"required,min=1,dive"`
	Customer      CreateCustomerDetailsRequest     `json:"customer" binding:"required"`
}
