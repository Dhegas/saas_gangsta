package domain

import (
	"context"

	"github.com/dhegas/saas_gangsta/internal/domains/order/dto"
)

// CustomerOrderUsecase kontrak logika bisnis untuk order dari halaman publik (customer self-order)
type CustomerOrderUsecase interface {
	CreateCustomerOrder(ctx context.Context, tenantID string, req dto.CreateCustomerOrderRequest) (*dto.CreateCustomerOrderResponse, error)
}

// CustomerOrderRepository kontrak interaksi database untuk customer self-order secara transaksional
type CustomerOrderRepository interface {
	CreateOrderWithCustomer(ctx context.Context, order *OrderEntity, items []OrderItemEntity, customer *CustomerEntity) error
	CheckTableExists(ctx context.Context, tenantID, tableID string) (bool, error)
	GetMenuDetails(ctx context.Context, tenantID string, menuIDs []string) (map[string]MenuDetail, error)
}
