package usecase

import (
	"context"
	"errors"
	"net/http"
	"testing"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/config"
	"github.com/dhegas/saas_gangsta/internal/domains/order/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/order/dto"
	authdomain "github.com/dhegas/saas_gangsta/internal/domains/user/auth/domain"
	"gorm.io/gorm"
)

type mockPartnerOrderRepo struct {
	findTableID    string
	findTableErr   error
	menuDetails    map[string]domain.MenuDetail
	menuErr        error
	createErr      error
	historyOrders  []domain.OrderEntity
	historyTenants map[string]domain.TenantInfo
	historyTables  map[string]string
	historyErr     error
}

func (m *mockPartnerOrderRepo) FindAll(ctx context.Context, tenantID string, filter dto.OrderFilterParams) ([]domain.OrderEntity, error) {
	return nil, nil
}
func (m *mockPartnerOrderRepo) FindByID(ctx context.Context, tenantID, orderID string) (*domain.OrderEntity, error) {
	return nil, nil
}
func (m *mockPartnerOrderRepo) CreateWithItems(ctx context.Context, order *domain.OrderEntity, items []domain.OrderItemEntity) error {
	return m.createErr
}
func (m *mockPartnerOrderRepo) UpdateStatus(ctx context.Context, tenantID, orderID, status string) error {
	return nil
}
func (m *mockPartnerOrderRepo) SoftDelete(ctx context.Context, tenantID, orderID string) error {
	return nil
}
func (m *mockPartnerOrderRepo) GetMenuDetails(ctx context.Context, tenantID string, menuIDs []string) (map[string]domain.MenuDetail, error) {
	return m.menuDetails, m.menuErr
}
func (m *mockPartnerOrderRepo) CheckTableExists(ctx context.Context, tenantID, tableID string) (bool, error) {
	return true, nil
}
func (m *mockPartnerOrderRepo) GetTableByName(ctx context.Context, tenantID, tableName string) (string, error) {
	return m.findTableID, m.findTableErr
}
func (m *mockPartnerOrderRepo) GetPublicOrderDetails(ctx context.Context, tenantID, orderID string) (*domain.OrderEntity, string, error) {
	return nil, "", nil
}
func (m *mockPartnerOrderRepo) FindAllPublicOrders(ctx context.Context, tenantID string, filter dto.PublicOrderFilterParams) ([]domain.OrderEntity, map[string]string, error) {
	return nil, nil, nil
}
func (m *mockPartnerOrderRepo) GetMaxQueueNumberToday(ctx context.Context, tenantID string) (int, error) {
	return 0, nil
}
func (m *mockPartnerOrderRepo) FindCustomerOrderHistory(ctx context.Context, userID string) ([]domain.OrderEntity, map[string]domain.TenantInfo, map[string]string, error) {
	return m.historyOrders, m.historyTenants, m.historyTables, m.historyErr
}

type mockAuthRepo struct {
	user *authdomain.User
}

func (m *mockAuthRepo) FindByID(ctx context.Context, id string) (*authdomain.User, error) {
	return m.user, nil
}
func (m *mockAuthRepo) FindByEmail(ctx context.Context, email string) (*authdomain.User, error) {
	return nil, nil
}
func (m *mockAuthRepo) CreateUser(ctx context.Context, user *authdomain.User) error {
	return nil
}
func (m *mockAuthRepo) FindPhoneNumber(ctx context.Context, userID string, role string) (string, error) {
	return "", nil
}

func TestCreateOrder_TableValidation(t *testing.T) {
	cfg := &config.Config{
		JWTSecret:            "secret",
		JWTAccessTokenExpiry: 3600,
	}

	userID := "user-123"
	userMock := &authdomain.User{
		ID:       userID,
		FullName: "Budi",
		IsActive: true,
	}

	t.Run("success takeaway order without table name", func(t *testing.T) {
		repo := &mockPartnerOrderRepo{
			menuDetails: map[string]domain.MenuDetail{
				"m-1": {ID: "m-1", Name: "Nasi Goreng", Price: 20000},
			},
		}
		authRepo := &mockAuthRepo{
			user: userMock,
		}
		uc := NewPartnerOrderUsecase(repo, authRepo, cfg, nil)

		req := dto.CreateOrderRequest{
			UserID: &userID,
			Items: []dto.CreateOrderItemRequest{
				{MenuID: "m-1", Quantity: 2},
			},
		}

		res, err := uc.CreateOrder(context.Background(), "t-1", req)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if res == nil {
			t.Fatal("expected response, got nil")
		}
	})

	t.Run("success dine-in order with valid table name", func(t *testing.T) {
		tableName := "Meja 1"
		repo := &mockPartnerOrderRepo{
			findTableID: "table-1",
			menuDetails: map[string]domain.MenuDetail{
				"m-1": {ID: "m-1", Name: "Nasi Goreng", Price: 20000},
			},
		}
		authRepo := &mockAuthRepo{
			user: userMock,
		}
		uc := NewPartnerOrderUsecase(repo, authRepo, cfg, nil)

		req := dto.CreateOrderRequest{
			UserID:          &userID,
			DiningTableName: &tableName,
			Items: []dto.CreateOrderItemRequest{
				{MenuID: "m-1", Quantity: 2},
			},
		}

		res, err := uc.CreateOrder(context.Background(), "t-1", req)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
		if res == nil {
			t.Fatal("expected response, got nil")
		}
	})

	t.Run("error dine-in order with non-existent table name", func(t *testing.T) {
		tableName := "Meja Palsu"
		repo := &mockPartnerOrderRepo{
			findTableErr: gorm.ErrRecordNotFound,
		}
		authRepo := &mockAuthRepo{
			user: userMock,
		}
		uc := NewPartnerOrderUsecase(repo, authRepo, cfg, nil)

		req := dto.CreateOrderRequest{
			UserID:          &userID,
			DiningTableName: &tableName,
			Items: []dto.CreateOrderItemRequest{
				{MenuID: "m-1", Quantity: 2},
			},
		}

		res, err := uc.CreateOrder(context.Background(), "t-1", req)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			if appErr.Code != "BAD_REQUEST" {
				t.Errorf("expected code BAD_REQUEST, got %s", appErr.Code)
			}
			if appErr.Status != http.StatusBadRequest {
				t.Errorf("expected status 400, got %d", appErr.Status)
			}
		} else {
			t.Fatalf("expected AppError, got %T", err)
		}
		if res != nil {
			t.Fatal("expected nil response, got one")
		}
	})
}

func TestGetCustomerOrderHistory(t *testing.T) {
	cfg := &config.Config{}
	userID := "user-123"

	t.Run("success fetching order history", func(t *testing.T) {
		repo := &mockPartnerOrderRepo{
			historyOrders: []domain.OrderEntity{
				{
					ID:             "order-1",
					TenantID:       "tenant-1",
					Status:         "PENDING",
					TotalPrice:     50000,
					QueueNumber:    "Q-1",
					PaymentMethod:  "CASH",
					DiningTablesID: pointerToString("table-1"),
					Items: []domain.OrderItemEntity{
						{MenuName: "Bakso", Quantity: 2, Subtotal: 50000},
					},
				},
			},
			historyTenants: map[string]domain.TenantInfo{
				"tenant-1": {ID: "tenant-1", Name: "Warung Bakso", Slug: "warung-bakso"},
			},
			historyTables: map[string]string{
				"table-1": "Meja Utama",
			},
		}

		uc := NewPartnerOrderUsecase(repo, &mockAuthRepo{}, cfg, nil)
		res, err := uc.GetCustomerOrderHistory(context.Background(), userID)
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}

		if len(res) != 1 {
			t.Fatalf("expected 1 order, got %d", len(res))
		}

		o := res[0]
		if o.ID != "order-1" || o.TenantName != "Warung Bakso" || o.TenantSlug != "warung-bakso" || o.DiningTable.TableName != "Meja Utama" {
			t.Errorf("unexpected order data mapping: %+v", o)
		}
	})

	t.Run("error fetching order history", func(t *testing.T) {
		repo := &mockPartnerOrderRepo{
			historyErr: errors.New("db error"),
		}

		uc := NewPartnerOrderUsecase(repo, &mockAuthRepo{}, cfg, nil)
		_, err := uc.GetCustomerOrderHistory(context.Background(), userID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func pointerToString(s string) *string {
	return &s
}

