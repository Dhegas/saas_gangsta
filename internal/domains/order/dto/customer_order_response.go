package dto

import "time"

// CreateCustomerOrderResponse response data setelah berhasil membuat order secara publik
type CreateCustomerOrderResponse struct {
	OrderID    string  `json:"orderId"`
	Status     string  `json:"status"`
	TotalPrice float64 `json:"totalPrice"`
}

type PublicOrderDetailsResponse struct {
	ID          string                    `json:"id"`
	Status      string                    `json:"status"`
	TotalPrice  float64                   `json:"totalPrice"`
	CreatedAt   time.Time                 `json:"createdAt"`
	Customer    PublicCustomerDetails     `json:"customer"`
	DiningTable PublicDiningTableDetails  `json:"diningTable"`
	Items       []PublicOrderItemResponse `json:"items"`
}

type PublicCustomerDetails struct {
	FullName string `json:"fullName"`
}

type PublicDiningTableDetails struct {
	TableName string `json:"tableName"`
}

type PublicOrderItemResponse struct {
	MenuName string  `json:"menuName"`
	Quantity int     `json:"quantity"`
	Notes    string  `json:"notes"`
	Subtotal float64 `json:"subtotal"`
}

type PublicOrderFilterParams struct {
	Status  string `form:"status" binding:"omitempty"`
	TableID string `form:"table_id" binding:"omitempty,uuid"`
}


