package dto

// CreateOrderItemRequest payload untuk detail item saat order
type CreateOrderItemRequest struct {
	MenuID   string `json:"menu_id" binding:"required,uuid"`
	Quantity int    `json:"quantity" binding:"required,min=1"`
	Notes    string `json:"notes" binding:"omitempty,max=255"`
}

// CreateOrderRequest payload untuk POST /api/orders (diakses oleh CUSTOMER)
type CreateOrderRequest struct {
	UserID          *string                  `json:"-"`
	DiningTablesID  *string                  `json:"dining_tables_id" binding:"omitempty,uuid"`
	DiningTableName *string                  `json:"dining_table_name" binding:"omitempty,max=50"`
	PaymentMethod   string                   `json:"payment_method" binding:"required,oneof=QRIS TRANSFER_BANK CASH E_WALLET KARTU_KREDIT MINIMARKET"`
	CustomerName    *string                  `json:"customer_name" binding:"omitempty,max=100"`
	Items           []CreateOrderItemRequest `json:"items" binding:"required,min=1,dive"`
}

// UpdateOrderStatusRequest payload untuk PATCH /api/orders/:id/status
type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=PENDING PROCESSING COMPLETED CANCELLED"`
}

// OrderFilterParams parameter kueri untuk GET /api/orders
type OrderFilterParams struct {
	Status  string `form:"status" binding:"omitempty,oneof=PENDING PROCESSING COMPLETED CANCELLED"`
	TableID string `form:"table_id" binding:"omitempty,uuid"`
	UserID  string `form:"-"`
}
