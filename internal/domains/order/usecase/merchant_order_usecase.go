package usecase

import (
	"context"
	"errors"
	"net/http"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/domains/order/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/order/dto"
	"gorm.io/gorm"
)

type merchantOrderUsecase struct {
	repo domain.MerchantOrderRepository
}

func NewMerchantOrderUsecase(repo domain.MerchantOrderRepository) domain.MerchantOrderUsecase {
	return &merchantOrderUsecase{repo: repo}
}

func (u *merchantOrderUsecase) GetAllOrders(ctx context.Context, tenantID string, filter dto.OrderFilterParams) ([]dto.OrderResponse, error) {
	orders, err := u.repo.FindAll(ctx, tenantID, filter)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil daftar pesanan", http.StatusInternalServerError, err)
	}

	result := make([]dto.OrderResponse, 0, len(orders))
	for _, o := range orders {
		result = append(result, toOrderResponse(&o))
	}
	return result, nil
}

func (u *merchantOrderUsecase) GetOrderByID(ctx context.Context, tenantID, orderID string) (*dto.OrderResponse, error) {
	order, err := u.repo.FindByID(ctx, tenantID, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New("NOT_FOUND", "Pesanan tidak ditemukan", http.StatusNotFound, nil)
		}
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil detail pesanan", http.StatusInternalServerError, err)
	}

	response := toOrderResponse(order)
	return &response, nil
}

func (u *merchantOrderUsecase) CreateOrder(ctx context.Context, tenantID string, req dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	// Kumpulkan ID menu yang dipesan
	menuIDs := make([]string, 0, len(req.Items))
	for _, item := range req.Items {
		menuIDs = append(menuIDs, item.MenuID)
	}

	// Ambil detail harga dan nama menu dari database untuk validasi & kalkulasi
	menuDetails, err := u.repo.GetMenuDetails(ctx, menuIDs)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal memvalidasi item pesanan", http.StatusInternalServerError, err)
	}

	var totalOrderPrice float64
	var orderItems []domain.OrderItemEntity

	// Bangun order_items dan kalkulasi subtotal menggunakan data valid dari DB
	for _, reqItem := range req.Items {
		detail, exists := menuDetails[reqItem.MenuID]
		if !exists {
			return nil, apperrors.New("BAD_REQUEST", "Salah satu menu yang dipesan tidak valid atau tidak tersedia", http.StatusBadRequest, nil)
		}

		subtotal := float64(reqItem.Quantity) * detail.Price
		totalOrderPrice += subtotal

		orderItems = append(orderItems, domain.OrderItemEntity{
			MenuID:    reqItem.MenuID,
			MenuName:  detail.Name,
			Quantity:  reqItem.Quantity,
			UnitPrice: detail.Price,
			Subtotal:  subtotal,
			Notes:     reqItem.Notes,
		})
	}

	orderEntity := &domain.OrderEntity{
		TenantID:       tenantID,
		DiningTablesID: req.DiningTablesID,
		Status:         "PENDING",
		TotalPrice:     totalOrderPrice,
	}

	// Simpan ke database dengan transaksi (Order + Items)
	if err := u.repo.CreateWithItems(ctx, orderEntity, orderItems); err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal menyimpan pesanan", http.StatusInternalServerError, err)
	}

	response := toOrderResponse(orderEntity)
	return &response, nil
}

func (u *merchantOrderUsecase) UpdateOrderStatus(ctx context.Context, tenantID, orderID string, req dto.UpdateOrderStatusRequest) (*dto.OrderResponse, error) {
	err := u.repo.UpdateStatus(ctx, tenantID, orderID, req.Status)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New("NOT_FOUND", "Pesanan tidak ditemukan", http.StatusNotFound, nil)
		}
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal memperbarui status pesanan", http.StatusInternalServerError, err)
	}

	return u.GetOrderByID(ctx, tenantID, orderID)
}

func (u *merchantOrderUsecase) SoftDeleteOrder(ctx context.Context, tenantID, orderID string) error {
	err := u.repo.SoftDelete(ctx, tenantID, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.New("NOT_FOUND", "Pesanan tidak ditemukan", http.StatusNotFound, nil)
		}
		return apperrors.New("INTERNAL_ERROR", "Gagal menghapus pesanan", http.StatusInternalServerError, err)
	}
	return nil
}

// helper untuk mapping response
func toOrderResponse(entity *domain.OrderEntity) dto.OrderResponse {
	var itemsResp []dto.OrderItemResponse
	if entity.Items != nil {
		for _, item := range entity.Items {
			itemsResp = append(itemsResp, dto.OrderItemResponse{
				ID:        item.ID,
				MenuID:    item.MenuID,
				MenuName:  item.MenuName,
				Quantity:  item.Quantity,
				UnitPrice: item.UnitPrice,
				Subtotal:  item.Subtotal,
				Notes:     item.Notes,
			})
		}
	}

	return dto.OrderResponse{
		ID:             entity.ID,
		TenantID:       entity.TenantID,
		UserID:         entity.UserID,
		DiningTablesID: entity.DiningTablesID,
		Status:         entity.Status,
		TotalPrice:     entity.TotalPrice,
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
		Items:          itemsResp,
	}
}
