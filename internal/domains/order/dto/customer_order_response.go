package dto

// CreateCustomerOrderResponse response data setelah berhasil membuat order secara publik
type CreateCustomerOrderResponse struct {
	OrderID    string  `json:"orderId"`
	Status     string  `json:"status"`
	TotalPrice float64 `json:"totalPrice"`
}
