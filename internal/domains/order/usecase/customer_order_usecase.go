package usecase

import (
	"context"
	"net/http"
	"strings"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/domains/order/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/order/dto"
)

type customerOrderUsecase struct {
	repo domain.CustomerOrderRepository
}

// NewCustomerOrderUsecase konstruktor untuk customerOrderUsecase
func NewCustomerOrderUsecase(repo domain.CustomerOrderRepository) domain.CustomerOrderUsecase {
	return &customerOrderUsecase{repo: repo}
}

func (u *customerOrderUsecase) CreateCustomerOrder(ctx context.Context, tenantID string, req dto.CreateCustomerOrderRequest) (*dto.CreateCustomerOrderResponse, error) {
	// 1. Validasi meja makan milik tenant
	tableExists, err := u.repo.CheckTableExists(ctx, tenantID, req.DiningTableID)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal memvalidasi meja makan", http.StatusInternalServerError, err)
	}
	if !tableExists {
		return nil, apperrors.New("BAD_REQUEST", "Meja makan tidak ditemukan atau bukan milik tenant ini", http.StatusBadRequest, nil)
	}

	// 2. Kumpulkan ID menu yang dipesan
	menuIDs := make([]string, 0, len(req.Items))
	for _, item := range req.Items {
		menuIDs = append(menuIDs, item.MenuID)
	}

	// 3. Ambil detail harga dan nama menu dari database untuk validasi & kalkulasi harga aman
	menuDetails, err := u.repo.GetMenuDetails(ctx, tenantID, menuIDs)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal memvalidasi item pesanan", http.StatusInternalServerError, err)
	}

	var totalOrderPrice float64
	var orderItems []domain.OrderItemEntity

	// 4. Bangun order_items dan hitung total harga
	for _, reqItem := range req.Items {
		detail, exists := menuDetails[reqItem.MenuID]
		if !exists {
			return nil, apperrors.New("BAD_REQUEST", "Salah satu menu yang dipesan tidak valid, tidak tersedia, atau bukan milik tenant ini", http.StatusBadRequest, nil)
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

	// 5. Bangun entitas Order
	orderEntity := &domain.OrderEntity{
		TenantID:       tenantID,
		DiningTablesID: req.DiningTableID,
		Status:         "PENDING",
		TotalPrice:     totalOrderPrice,
	}

	// 6. Bangun entitas Customer
	customerEntity := &domain.CustomerEntity{
		TenantID:    tenantID,
		FullName:    req.Customer.FullName,
		PhoneNumber: req.Customer.PhoneNumber,
	}

	// 7. Simpan secara transaksional
	if err := u.repo.CreateOrderWithCustomer(ctx, orderEntity, orderItems, customerEntity); err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal menyimpan pesanan", http.StatusInternalServerError, err)
	}

	// 8. Bentuk response (status disesuaikan ke lowercase sesuai spesifikasi)
	return &dto.CreateCustomerOrderResponse{
		OrderID:    orderEntity.ID,
		Status:     strings.ToLower(orderEntity.Status),
		TotalPrice: orderEntity.TotalPrice,
	}, nil
}
