package repository

import (
	"context"

	orderdomain "github.com/dhegas/saas_gangsta/internal/domains/order/domain"
	"gorm.io/gorm"
)

type customerOrderRepository struct {
	db *gorm.DB
}

// NewCustomerOrderRepository konstruktor untuk customerOrderRepository
func NewCustomerOrderRepository(db *gorm.DB) orderdomain.CustomerOrderRepository {
	return &customerOrderRepository{db: db}
}

func (r *customerOrderRepository) CreateOrderWithCustomer(ctx context.Context, order *orderdomain.OrderEntity, items []orderdomain.OrderItemEntity, customer *orderdomain.CustomerEntity) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Simpan order
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// 2. Set OrderID ke masing-masing item, lalu simpan ke database
		for i := range items {
			items[i].OrderID = order.ID
		}
		if err := tx.Create(&items).Error; err != nil {
			return err
		}

		// 3. Set OrderID ke data customer, lalu simpan ke database
		customer.OrderID = order.ID
		if err := tx.Create(customer).Error; err != nil {
			return err
		}

		order.Items = items
		return nil
	})
}

func (r *customerOrderRepository) CheckTableExists(ctx context.Context, tenantID, tableID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("dining_tables").
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", tableID, tenantID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *customerOrderRepository) GetMenuDetails(ctx context.Context, tenantID string, menuIDs []string) (map[string]orderdomain.MenuDetail, error) {
	var menus []struct {
		ID    string
		Name  string
		Price float64
	}

	err := r.db.WithContext(ctx).Table("menus").
		Select("id, name, price").
		Where("tenant_id = ? AND id IN ? AND is_available = true AND deleted_at IS NULL", tenantID, menuIDs).
		Find(&menus).Error

	if err != nil {
		return nil, err
	}

	result := make(map[string]orderdomain.MenuDetail)
	for _, m := range menus {
		result[m.ID] = orderdomain.MenuDetail{
			ID:    m.ID,
			Name:  m.Name,
			Price: m.Price,
		}
	}
	return result, nil
}
