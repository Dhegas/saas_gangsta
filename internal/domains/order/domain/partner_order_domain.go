package domain

import (
	"context"
	"time"

	"github.com/dhegas/saas_gangsta/internal/domains/order/dto"
)

// PartnerOrderUsecase kontrak logika bisnis
type PartnerOrderUsecase interface {
	GetAllOrders(ctx context.Context, tenantID string, filter dto.OrderFilterParams) ([]dto.OrderResponse, error)
	GetOrderByID(ctx context.Context, tenantID, orderID string) (*dto.OrderResponse, error)
	CreateOrder(ctx context.Context, tenantID string, req dto.CreateOrderRequest) (*dto.OrderResponse, error)
	UpdateOrderStatus(ctx context.Context, tenantID, orderID string, req dto.UpdateOrderStatusRequest) (*dto.OrderResponse, error)
	SoftDeleteOrder(ctx context.Context, tenantID, orderID string) error
}

// PartnerOrderRepository kontrak interaksi database
type PartnerOrderRepository interface {
	FindAll(ctx context.Context, tenantID string, filter dto.OrderFilterParams) ([]OrderEntity, error)
	FindByID(ctx context.Context, tenantID, orderID string) (*OrderEntity, error)
	CreateWithItems(ctx context.Context, order *OrderEntity, items []OrderItemEntity) error
	UpdateStatus(ctx context.Context, tenantID, orderID, status string) error
	SoftDelete(ctx context.Context, tenantID, orderID string) error
	GetMenuDetails(ctx context.Context, menuIDs []string) (map[string]MenuDetail, error)
}

type MenuDetail struct {
	ID    string
	Name  string
	Price float64
}

// OrderEntity merepresentasikan tabel orders
type OrderEntity struct {
	ID             string            `gorm:"primaryKey;default:gen_random_uuid()"`
	TenantID       string            `gorm:"index;not null"`
	UserID         *string           `gorm:"index"`
	DiningTablesID string            `gorm:"index"`
	Status         string            `gorm:"not null"`
	TotalPrice     float64           `gorm:"not null"`
	CreatedAt      time.Time         `gorm:"autoCreateTime"`
	UpdatedAt      time.Time         `gorm:"autoUpdateTime"`
	DeletedAt      *time.Time        `gorm:"index"`
	Items          []OrderItemEntity `gorm:"foreignKey:OrderID"`
}

func (OrderEntity) TableName() string {
	return "orders"
}

// OrderItemEntity merepresentasikan tabel order_items
type OrderItemEntity struct {
	ID        string     `gorm:"primaryKey;default:gen_random_uuid()"`
	OrderID   string     `gorm:"index;not null"`
	MenuID    string     `gorm:"index;not null"`
	MenuName  string     `gorm:"not null"`
	Quantity  int        `gorm:"not null"`
	UnitPrice float64    `gorm:"not null"`
	Subtotal  float64    `gorm:"not null"`
	Notes     string     ``
	DeletedAt *time.Time `gorm:"index"`
}

func (OrderItemEntity) TableName() string {
	return "order_items"
}
