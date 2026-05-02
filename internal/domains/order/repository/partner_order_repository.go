package repository

import (
	"context"
	"time"

	orderdomain "github.com/dhegas/saas_gangsta/internal/domains/order/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/order/dto"
	"gorm.io/gorm"
)

type partnerOrderRepository struct {
	db *gorm.DB
}

func NewPartnerOrderRepository(db *gorm.DB) orderdomain.PartnerOrderRepository {
	return &partnerOrderRepository{db: db}
}

func (r *partnerOrderRepository) FindAll(ctx context.Context, tenantID string, filter dto.OrderFilterParams) ([]orderdomain.OrderEntity, error) {
	var orders []orderdomain.OrderEntity
	query := r.db.WithContext(ctx).Preload("Items").Where("tenant_id = ? AND deleted_at IS NULL", tenantID)

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.TableID != "" {
		query = query.Where("dining_tables_id = ?", filter.TableID)
	}

	err := query.Order("created_at DESC").Find(&orders).Error
	return orders, err
}

func (r *partnerOrderRepository) FindByID(ctx context.Context, tenantID, orderID string) (*orderdomain.OrderEntity, error) {
	var order orderdomain.OrderEntity
	err := r.db.WithContext(ctx).Preload("Items").
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", orderID, tenantID).
		First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *partnerOrderRepository) CreateWithItems(ctx context.Context, order *orderdomain.OrderEntity, items []orderdomain.OrderItemEntity) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		for i := range items {
			items[i].OrderID = order.ID
		}

		if err := tx.Create(&items).Error; err != nil {
			return err
		}

		order.Items = items
		return nil
	})
}

func (r *partnerOrderRepository) UpdateStatus(ctx context.Context, tenantID, orderID, status string) error {
	res := r.db.WithContext(ctx).Model(&orderdomain.OrderEntity{}).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", orderID, tenantID).
		Update("status", status)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *partnerOrderRepository) SoftDelete(ctx context.Context, tenantID, orderID string) error {
	now := time.Now()
	res := r.db.WithContext(ctx).Model(&orderdomain.OrderEntity{}).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", orderID, tenantID).
		Update("deleted_at", &now)

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *partnerOrderRepository) GetMenuDetails(ctx context.Context, menuIDs []string) (map[string]orderdomain.MenuDetail, error) {
	var menus []struct {
		ID    string
		Name  string
		Price float64
	}

	err := r.db.WithContext(ctx).Table("menus").
		Select("id, name, price").
		Where("id IN ?", menuIDs).
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
