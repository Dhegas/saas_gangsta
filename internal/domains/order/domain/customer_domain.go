package domain

import (
	"context"
	"time"

	"github.com/dhegas/saas_gangsta/internal/domains/order/dto"
)

// CustomerUsecase kontrak logika bisnis customer order
type CustomerUsecase interface {
	CreateCustomer(ctx context.Context, tenantID, orderID string, req dto.CreateCustomerRequest) (*dto.CustomerResponse, error)
	GetCustomerByOrderID(ctx context.Context, tenantID, orderID string) (*dto.CustomerResponse, error)
	UpdateCustomer(ctx context.Context, tenantID, orderID string, req dto.UpdateCustomerRequest) (*dto.CustomerResponse, error)
}

// CustomerRepository kontrak interaksi database
type CustomerRepository interface {
	Create(ctx context.Context, customer *CustomerEntity) error
	FindByOrderID(ctx context.Context, tenantID, orderID string) (*CustomerEntity, error)
	Update(ctx context.Context, customer *CustomerEntity) error
	OrderExists(ctx context.Context, tenantID, orderID string) (bool, error)
}

// CustomerEntity merepresentasikan tabel customers
type CustomerEntity struct {
	ID          string     `gorm:"primaryKey;default:gen_random_uuid()"`
	OrderID     string     `gorm:"index;not null"`
	TenantID    string     `gorm:"index;not null"`
	FullName    string     `gorm:"not null"`
	PhoneNumber string     ``
	CreatedAt   time.Time  `gorm:"autoCreateTime"`
	DeletedAt   *time.Time `gorm:"index"`
}

func (CustomerEntity) TableName() string {
	return "customers"
}
