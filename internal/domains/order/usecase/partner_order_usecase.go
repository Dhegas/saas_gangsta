package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dhegas/saas_gangsta/internal/common/cache"
	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/config"
	"github.com/dhegas/saas_gangsta/internal/domains/order/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/order/dto"
	authrepo "github.com/dhegas/saas_gangsta/internal/domains/user/auth/repository"
	"github.com/dhegas/saas_gangsta/internal/infrastructure/websocket"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var orderCache = cache.NewLocalCache()

type partnerOrderUsecase struct {
	repo     domain.PartnerOrderRepository
	authRepo authrepo.AuthRepository
	cfg      *config.Config
	wsHub    *websocket.Hub
}

func NewPartnerOrderUsecase(repo domain.PartnerOrderRepository, authRepo authrepo.AuthRepository, cfg *config.Config, wsHub *websocket.Hub) domain.PartnerOrderUsecase {
	return &partnerOrderUsecase{
		repo:     repo,
		authRepo: authRepo,
		cfg:      cfg,
		wsHub:    wsHub,
	}
}

func (u *partnerOrderUsecase) GetAllOrders(ctx context.Context, tenantID string, filter dto.OrderFilterParams) ([]dto.OrderResponse, error) {
	orders, err := u.repo.FindAll(ctx, tenantID, filter)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil daftar pesanan", http.StatusInternalServerError)
	}

	result := make([]dto.OrderResponse, 0, len(orders))
	for _, o := range orders {
		result = append(result, toOrderResponse(&o))
	}
	return result, nil
}

func (u *partnerOrderUsecase) GetOrderByID(ctx context.Context, tenantID, orderID string) (*dto.OrderResponse, error) {
	order, err := u.repo.FindByID(ctx, tenantID, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New("NOT_FOUND", "Pesanan tidak ditemukan", http.StatusNotFound)
		}
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil detail pesanan", http.StatusInternalServerError)
	}

	response := toOrderResponse(order)
	return &response, nil
}

func (u *partnerOrderUsecase) CreateOrder(ctx context.Context, tenantID string, req dto.CreateOrderRequest) (*dto.OrderResponse, error) {
	// 1. Validasi / dapatkan meja makan milik tenant jika ada
	var diningTableID *string
	if req.DiningTableName != nil && *req.DiningTableName != "" {
		tableID, err := u.repo.GetTableByName(ctx, tenantID, *req.DiningTableName)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "record not found" {
				return nil, apperrors.New("BAD_REQUEST", "Meja makan dengan nomor "+*req.DiningTableName+" tidak ditemukan", http.StatusBadRequest)
			}
			return nil, apperrors.New("INTERNAL_ERROR", "Gagal memproses meja makan", http.StatusInternalServerError)
		}
		diningTableID = &tableID
	} else if req.DiningTablesID != nil && *req.DiningTablesID != "" {
		tableExists, err := u.repo.CheckTableExists(ctx, tenantID, *req.DiningTablesID)
		if err != nil {
			return nil, apperrors.New("INTERNAL_ERROR", "Gagal memvalidasi meja makan", http.StatusInternalServerError)
		}
		if !tableExists {
			return nil, apperrors.New("BAD_REQUEST", "Meja makan tidak ditemukan atau bukan milik tenant ini", http.StatusBadRequest)
		}
		diningTableID = req.DiningTablesID
	}

	// 2. Kumpulkan ID menu yang dipesan
	menuIDs := make([]string, 0, len(req.Items))
	for _, item := range req.Items {
		menuIDs = append(menuIDs, item.MenuID)
	}

	// 3. Ambil detail harga dan nama menu dari database untuk validasi & kalkulasi harga aman (tenant-isolated)
	menuDetails, err := u.repo.GetMenuDetails(ctx, tenantID, menuIDs)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal memvalidasi item pesanan", http.StatusInternalServerError)
	}

	var totalOrderPrice float64
	var orderItems []domain.OrderItemEntity

	// 4. Bangun order_items dan kalkulasi subtotal menggunakan data valid dari DB
	for _, reqItem := range req.Items {
		detail, exists := menuDetails[reqItem.MenuID]
		if !exists {
			return nil, apperrors.New("BAD_REQUEST", "Salah satu menu yang dipesan tidak valid, tidak tersedia, atau bukan milik tenant ini", http.StatusBadRequest)
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

	if req.UserID == nil || *req.UserID == "" {
		return nil, apperrors.New("UNAUTHORIZED", "User ID diperlukan untuk membuat pesanan", http.StatusUnauthorized)
	}

	// Cek apakah user ada di database dan statusnya aktif
	authUser, err := u.authRepo.FindByID(ctx, *req.UserID)
	if err != nil || authUser == nil {
		return nil, apperrors.New("UNAUTHORIZED", "Pengguna tidak ditemukan di database", http.StatusUnauthorized)
	}
	if !authUser.IsActive {
		return nil, apperrors.New("FORBIDDEN", "Akun pengguna tidak aktif", http.StatusForbidden)
	}

	// 5. Bangun entitas Order (tanpa mengisi relasi User agar GORM tidak mencoba melakukan INSERT ke tabel users)
	orderEntity := &domain.OrderEntity{
		TenantID:       tenantID,
		UserID:         req.UserID,
		DiningTablesID: diningTableID,
		Status:         "PENDING",
		TotalPrice:     totalOrderPrice,
		CustomerName:   authUser.FullName,
	}

	var saveErr error
	for attempt := 1; attempt <= 5; attempt++ {
		// Reset ID agar GORM membuat UUID baru saat retry
		orderEntity.ID = ""
		for i := range orderItems {
			orderItems[i].ID = ""
		}

		// Ambil antrian terbesar hari ini dan tentukan antrian berikutnya
		maxQueueVal, err := u.repo.GetMaxQueueNumberToday(ctx, tenantID)
		if err != nil {
			return nil, apperrors.New("INTERNAL_ERROR", "Gagal mendapatkan nomor antrian", http.StatusInternalServerError)
		}

		orderEntity.QueueNumber = fmt.Sprintf("Q-%d", maxQueueVal+1)
		orderEntity.PaymentMethod = req.PaymentMethod

		// 6. Simpan secara transaksional
		saveErr = u.repo.CreateWithItems(ctx, orderEntity, orderItems)
		if saveErr == nil {
			break
		}

		// Periksa jika kesalahan adalah Unique Violation (tabrakan queue_number harian)
		var pgErr *pgconn.PgError
		isUniqueViolation := false
		if errors.As(saveErr, &pgErr) && pgErr.Code == "23505" {
			isUniqueViolation = true
		} else if strings.Contains(saveErr.Error(), "23505") || strings.Contains(saveErr.Error(), "unique constraint") {
			isUniqueViolation = true
		}

		if isUniqueViolation {
			// Tabrakan terdeteksi, lanjutkan ke iterasi berikutnya untuk mencoba antrian yang lebih tinggi
			continue
		}

		// Jika error lain, langsung gagalkan proses
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal menyimpan pesanan", http.StatusInternalServerError)
	}

	if saveErr != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal menyimpan pesanan setelah beberapa kali percobaan karena tabrakan nomor antrian", http.StatusInternalServerError)
	}

	// Isi relasi User setelah penyimpanan berhasil untuk keperluan format response
	orderEntity.User = &domain.UserEntity{
		ID:       authUser.ID,
		FullName: authUser.FullName,
		Email:    authUser.Email,
	}

	response := toOrderResponse(orderEntity)

	if u.wsHub != nil {
		customerID := ""
		if response.UserID != nil {
			customerID = *response.UserID
		}
		u.wsHub.SendTo(tenantID, map[string]interface{}{
			"type":        "new_order",
			"order_id":    response.ID,
			"customer_id": customerID,
			"status":      strings.ToLower(response.Status),
		})
	}

	return &response, nil
}


func (u *partnerOrderUsecase) UpdateOrderStatus(ctx context.Context, tenantID, orderID string, req dto.UpdateOrderStatusRequest) (*dto.OrderResponse, error) {
	err := u.repo.UpdateStatus(ctx, tenantID, orderID, req.Status)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New("NOT_FOUND", "Pesanan tidak ditemukan", http.StatusNotFound)
		}
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal memperbarui status pesanan", http.StatusInternalServerError)
	}

	return u.GetOrderByID(ctx, tenantID, orderID)
}

func (u *partnerOrderUsecase) SoftDeleteOrder(ctx context.Context, tenantID, orderID string) error {
	err := u.repo.SoftDelete(ctx, tenantID, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.New("NOT_FOUND", "Pesanan tidak ditemukan", http.StatusNotFound)
		}
		return apperrors.New("INTERNAL_ERROR", "Gagal menghapus pesanan", http.StatusInternalServerError)
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

	customerName := entity.CustomerName
	if customerName == "" && entity.User != nil {
		customerName = entity.User.FullName
	}

	return dto.OrderResponse{
		ID:             entity.ID,
		TenantID:       entity.TenantID,
		UserID:         entity.UserID,
		DiningTablesID: entity.DiningTablesID,
		Status:         entity.Status,
		TotalPrice:     entity.TotalPrice,
		QueueNumber:    entity.QueueNumber,
		PaymentMethod:  entity.PaymentMethod,
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
		Items:          itemsResp,
		CustomerName:   customerName,
	}
}

func (u *partnerOrderUsecase) GetPublicOrderStatus(ctx context.Context, tenantID, orderID string) (*dto.PublicOrderDetailsResponse, error) {
	order, tableName, err := u.repo.GetPublicOrderDetails(ctx, tenantID, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New("NOT_FOUND", "Pesanan tidak ditemukan", http.StatusNotFound)
		}
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil status pesanan", http.StatusInternalServerError)
	}

	customerName := order.CustomerName
	if customerName == "" && order.User != nil {
		customerName = order.User.FullName
	}
	customerResp := dto.PublicCustomerDetails{
		FullName: customerName,
	}

	itemsResp := make([]dto.PublicOrderItemResponse, 0, len(order.Items))
	for _, item := range order.Items {
		itemsResp = append(itemsResp, dto.PublicOrderItemResponse{
			MenuName: item.MenuName,
			Quantity: item.Quantity,
			Notes:    item.Notes,
			Subtotal: item.Subtotal,
		})
	}

	return &dto.PublicOrderDetailsResponse{
		ID:            order.ID,
		Status:        order.Status,
		TotalPrice:    order.TotalPrice,
		QueueNumber:   order.QueueNumber,
		PaymentMethod: order.PaymentMethod,
		CreatedAt:     order.CreatedAt,
		UserID:        order.UserID,
		Customer:      customerResp,
		DiningTable: dto.PublicDiningTableDetails{
			TableName: tableName,
		},
		Items: itemsResp,
	}, nil
}

func (u *partnerOrderUsecase) GetPublicOrdersList(ctx context.Context, tenantID string, filter dto.PublicOrderFilterParams) ([]dto.PublicOrderDetailsResponse, error) {
	cacheKey := fmt.Sprintf("customer:orders:tenant:%s:status:%s:table:%s:user:%s", tenantID, filter.Status, filter.TableID, filter.UserID)
	if cached, found := orderCache.Get(cacheKey); found {
		if cachedOrders, ok := cached.([]dto.PublicOrderDetailsResponse); ok {
			return cachedOrders, nil
		}
	}

	orders, tableNames, err := u.repo.FindAllPublicOrders(ctx, tenantID, filter)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil daftar pesanan", http.StatusInternalServerError)
	}

	result := make([]dto.PublicOrderDetailsResponse, 0, len(orders))
	for _, o := range orders {
		customerName := o.CustomerName
		if customerName == "" && o.User != nil {
			customerName = o.User.FullName
		}
		customerResp := dto.PublicCustomerDetails{
			FullName: customerName,
		}

		itemsResp := make([]dto.PublicOrderItemResponse, 0, len(o.Items))
		for _, item := range o.Items {
			itemsResp = append(itemsResp, dto.PublicOrderItemResponse{
				MenuName: item.MenuName,
				Quantity: item.Quantity,
				Notes:    item.Notes,
				Subtotal: item.Subtotal,
			})
		}

		var tableName string
		if o.DiningTablesID != nil {
			tableName = tableNames[*o.DiningTablesID]
		}

		result = append(result, dto.PublicOrderDetailsResponse{
			ID:            o.ID,
			Status:        o.Status,
			TotalPrice:    o.TotalPrice,
			QueueNumber:   o.QueueNumber,
			PaymentMethod: o.PaymentMethod,
			CreatedAt:     o.CreatedAt,
			UserID:        o.UserID,
			Customer:      customerResp,
			DiningTable: dto.PublicDiningTableDetails{
				TableName: tableName,
			},
			Items: itemsResp,
		})
	}

	orderCache.Set(cacheKey, result, 1*time.Minute)

	return result, nil
}


