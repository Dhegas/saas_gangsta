package dto

import "time"

type OrderItemResponse struct {
	ID        string  `json:"id"`
	MenuID    string  `json:"menu_id"`
	MenuName  string  `json:"menu_name"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Subtotal  float64 `json:"subtotal"`
	Notes     string  `json:"notes"`
}

type OrderResponse struct {
	ID             string              `json:"id"`
	TenantID       string              `json:"tenant_id"`
	UserID         *string             `json:"user_id,omitempty"`
	DiningTablesID *string             `json:"dining_tables_id,omitempty"`
	Status         string              `json:"status"`
	TotalPrice     float64             `json:"total_price"`
	QueueNumber    string              `json:"queue_number"`
	PaymentMethod  string              `json:"payment_method"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
	Items          []OrderItemResponse `json:"items,omitempty"`
	CustomerName   string              `json:"customer_name,omitempty"`
	AccessToken    string              `json:"access_token,omitempty"`
}
