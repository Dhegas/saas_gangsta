package repository

import (
	"context"
	"errors"

	orderdomain "github.com/dhegas/saas_gangsta/internal/domains/order/domain"
	"gorm.io/gorm"
)

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) orderdomain.CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(ctx context.Context, customer *orderdomain.CustomerEntity) error {
	return r.db.WithContext(ctx).Create(customer).Error
}

func (r *customerRepository) FindByOrderID(ctx context.Context, tenantID, orderID string) (*orderdomain.CustomerEntity, error) {
	var customer orderdomain.CustomerEntity
	err := r.db.WithContext(ctx).
		Where("order_id = ? AND tenant_id = ? AND deleted_at IS NULL", orderID, tenantID).
		First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) Update(ctx context.Context, customer *orderdomain.CustomerEntity) error {
	res := r.db.WithContext(ctx).
		Model(customer).
		Where("id = ? AND deleted_at IS NULL", customer.ID).
		Updates(map[string]interface{}{
			"full_name":    customer.FullName,
			"phone_number": customer.PhoneNumber,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *customerRepository) OrderExists(ctx context.Context, tenantID, orderID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("orders").
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", orderID, tenantID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CustomerAlreadyExistsError digunakan ketika customer sudah ada untuk order tertentu
type CustomerAlreadyExistsError struct{}

func (e *CustomerAlreadyExistsError) Error() string {
	return "customer already exists for this order"
}

func IsCustomerAlreadyExistsError(err error) bool {
	var target *CustomerAlreadyExistsError
	return errors.As(err, &target)
}
