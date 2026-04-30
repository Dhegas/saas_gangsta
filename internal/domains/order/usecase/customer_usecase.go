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

type customerUsecase struct {
	repo domain.CustomerRepository
}

func NewCustomerUsecase(repo domain.CustomerRepository) domain.CustomerUsecase {
	return &customerUsecase{repo: repo}
}

func (u *customerUsecase) CreateCustomer(ctx context.Context, tenantID, orderID string, req dto.CreateCustomerRequest) (*dto.CustomerResponse, error) {
	// Pastikan order ada dan milik tenant ini
	exists, err := u.repo.OrderExists(ctx, tenantID, orderID)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal memverifikasi order", http.StatusInternalServerError, err)
	}
	if !exists {
		return nil, apperrors.New("NOT_FOUND", "Order tidak ditemukan", http.StatusNotFound, nil)
	}

	// Cek apakah customer sudah ada untuk order ini
	existing, err := u.repo.FindByOrderID(ctx, tenantID, orderID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal memeriksa data customer", http.StatusInternalServerError, err)
	}
	if existing != nil {
		return nil, apperrors.New("CONFLICT", "Customer sudah terdaftar untuk order ini", http.StatusConflict, nil)
	}

	customer := &domain.CustomerEntity{
		OrderID:     orderID,
		TenantID:    tenantID,
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
	}

	if err := u.repo.Create(ctx, customer); err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal menyimpan data customer", http.StatusInternalServerError, err)
	}

	return toCustomerResponse(customer), nil
}

func (u *customerUsecase) GetCustomerByOrderID(ctx context.Context, tenantID, orderID string) (*dto.CustomerResponse, error) {
	// Pastikan order ada
	exists, err := u.repo.OrderExists(ctx, tenantID, orderID)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal memverifikasi order", http.StatusInternalServerError, err)
	}
	if !exists {
		return nil, apperrors.New("NOT_FOUND", "Order tidak ditemukan", http.StatusNotFound, nil)
	}

	customer, err := u.repo.FindByOrderID(ctx, tenantID, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New("NOT_FOUND", "Data customer untuk order ini belum ada", http.StatusNotFound, nil)
		}
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil data customer", http.StatusInternalServerError, err)
	}

	return toCustomerResponse(customer), nil
}

func (u *customerUsecase) UpdateCustomer(ctx context.Context, tenantID, orderID string, req dto.UpdateCustomerRequest) (*dto.CustomerResponse, error) {
	// Pastikan order ada
	exists, err := u.repo.OrderExists(ctx, tenantID, orderID)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal memverifikasi order", http.StatusInternalServerError, err)
	}
	if !exists {
		return nil, apperrors.New("NOT_FOUND", "Order tidak ditemukan", http.StatusNotFound, nil)
	}

	customer, err := u.repo.FindByOrderID(ctx, tenantID, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New("NOT_FOUND", "Data customer untuk order ini belum ada", http.StatusNotFound, nil)
		}
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil data customer", http.StatusInternalServerError, err)
	}

	customer.FullName = req.FullName
	customer.PhoneNumber = req.PhoneNumber

	if err := u.repo.Update(ctx, customer); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New("NOT_FOUND", "Data customer tidak ditemukan", http.StatusNotFound, nil)
		}
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal memperbarui data customer", http.StatusInternalServerError, err)
	}

	return toCustomerResponse(customer), nil
}

// helper untuk mapping response
func toCustomerResponse(entity *domain.CustomerEntity) *dto.CustomerResponse {
	return &dto.CustomerResponse{
		ID:          entity.ID,
		OrderID:     entity.OrderID,
		TenantID:    entity.TenantID,
		FullName:    entity.FullName,
		PhoneNumber: entity.PhoneNumber,
		CreatedAt:   entity.CreatedAt,
		DeletedAt:   entity.DeletedAt,
	}
}
