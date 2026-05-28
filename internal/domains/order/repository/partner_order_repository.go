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
	if filter.UserID != "" {
		query = query.Where("user_id = ?", filter.UserID)
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

func (r *partnerOrderRepository) GetMenuDetails(ctx context.Context, tenantID string, menuIDs []string) (map[string]orderdomain.MenuDetail, error) {
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

func (r *partnerOrderRepository) CheckTableExists(ctx context.Context, tenantID, tableID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("dining_tables").
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", tableID, tenantID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *partnerOrderRepository) GetPublicOrderDetails(ctx context.Context, tenantID, orderID string) (*orderdomain.OrderEntity, string, error) {
	var order orderdomain.OrderEntity
	err := r.db.WithContext(ctx).Preload("Items").
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", orderID, tenantID).
		First(&order).Error
	if err != nil {
		return nil, "", err
	}

	var tableName string
	if order.DiningTablesID != "" {
		err = r.db.WithContext(ctx).Table("dining_tables").
			Select("table_name").
			Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", order.DiningTablesID, tenantID).
			Scan(&tableName).Error
		if err != nil {
			return nil, "", err
		}
	}

	return &order, tableName, nil
}

func (r *partnerOrderRepository) FindAllPublicOrders(ctx context.Context, tenantID string, filter dto.PublicOrderFilterParams) ([]orderdomain.OrderEntity, map[string]string, error) {
	var orders []orderdomain.OrderEntity
	query := r.db.WithContext(ctx).Preload("Items").Where("tenant_id = ? AND deleted_at IS NULL", tenantID)

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.TableID != "" {
		query = query.Where("dining_tables_id = ?", filter.TableID)
	}

	err := query.Order("created_at DESC").Find(&orders).Error
	if err != nil {
		return nil, nil, err
	}

	if len(orders) == 0 {
		return nil, make(map[string]string), nil
	}

	orderIDs := make([]string, 0, len(orders))
	tableIDs := make([]string, 0, len(orders))
	for _, o := range orders {
		orderIDs = append(orderIDs, o.ID)
		if o.DiningTablesID != "" {
			tableIDs = append(tableIDs, o.DiningTablesID)
		}
	}

	tableNames := make(map[string]string)
	if len(tableIDs) > 0 {
		var tables []struct {
			ID        string
			TableName string `gorm:"column:table_name"`
		}
		err = r.db.WithContext(ctx).Table("dining_tables").
			Select("id, table_name").
			Where("id IN ? AND tenant_id = ? AND deleted_at IS NULL", tableIDs, tenantID).
			Scan(&tables).Error
		if err != nil {
			return nil, nil, err
		}
		for _, t := range tables {
			tableNames[t.ID] = t.TableName
		}
	}

	return orders, tableNames, nil
}



