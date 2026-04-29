package usecase

import (
	"context"
	"errors"
	"net/http"

	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/domains/table/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/table/dto"
	"gorm.io/gorm"
)

type merchantTableUsecase struct {
	repo domain.MerchantTableRepository
}

func NewMerchantTableUsecase(repo domain.MerchantTableRepository) domain.MerchantTableUsecase {
	return &merchantTableUsecase{repo: repo}
}

func (u *merchantTableUsecase) GetAllTables(ctx context.Context, tenantID string) ([]dto.TableResponse, error) {
	tables, err := u.repo.FindAllByTenant(ctx, tenantID)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil data meja", http.StatusInternalServerError, err)
	}

	result := make([]dto.TableResponse, 0, len(tables))
	for _, t := range tables {
		result = append(result, toTableResponse(&t))
	}

	return result, nil
}

func (u *merchantTableUsecase) GetTableByID(ctx context.Context, tenantID, tableID string) (*dto.TableResponse, error) {
	table, err := u.repo.FindByIDAndTenant(ctx, tenantID, tableID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New("NOT_FOUND", "Meja tidak ditemukan", http.StatusNotFound, nil)
		}
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil data meja", http.StatusInternalServerError, err)
	}

	response := toTableResponse(table)
	return &response, nil
}

func (u *merchantTableUsecase) GetTableStatus(ctx context.Context, tenantID, tableID string) (*dto.TableStatusResponse, error) {
	table, err := u.repo.FindByIDAndTenant(ctx, tenantID, tableID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New("NOT_FOUND", "Meja tidak ditemukan", http.StatusNotFound, nil)
		}
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal memverifikasi meja", http.StatusInternalServerError, err)
	}

	isOccupied, err := u.repo.CheckTableOccupied(ctx, tableID)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengecek status meja", http.StatusInternalServerError, err)
	}

	status := "kosong"
	if isOccupied {
		status = "occupied"
	}

	return &dto.TableStatusResponse{
		ID:        table.ID,
		TableName: table.Name,
		Status:    status,
	}, nil
}

func (u *merchantTableUsecase) CreateTable(ctx context.Context, tenantID string, req dto.CreateTableRequest) (*dto.TableResponse, error) {
	exists, err := u.repo.CheckNameExists(ctx, tenantID, req.TableName, "")
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal memvalidasi nama meja", http.StatusInternalServerError, err)
	}
	if exists {
		return nil, apperrors.New("CONFLICT", "Nama meja sudah ada", http.StatusConflict, nil)
	}

	entity := &domain.DiningTableEntity{
		TenantID:  tenantID,
		Name: req.TableName,
	}

	if err := u.repo.Create(ctx, entity); err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal menyimpan meja", http.StatusInternalServerError, err)
	}

	response := toTableResponse(entity)
	return &response, nil
}

func (u *merchantTableUsecase) UpdateTable(ctx context.Context, tenantID, tableID string, req dto.UpdateTableRequest) (*dto.TableResponse, error) {
	table, err := u.repo.FindByIDAndTenant(ctx, tenantID, tableID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New("NOT_FOUND", "Meja tidak ditemukan", http.StatusNotFound, nil)
		}
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil data meja", http.StatusInternalServerError, err)
	}

	if req.TableName != "" && req.TableName != table.Name {
		exists, err := u.repo.CheckNameExists(ctx, tenantID, req.TableName, tableID)
		if err != nil {
			return nil, apperrors.New("INTERNAL_ERROR", "Gagal memvalidasi nama meja", http.StatusInternalServerError, err)
		}
		if exists {
			return nil, apperrors.New("CONFLICT", "Nama meja sudah ada", http.StatusConflict, nil)
		}
		table.Name = req.TableName
	}

	if err := u.repo.Update(ctx, table); err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal memperbarui meja", http.StatusInternalServerError, err)
	}

	response := toTableResponse(table)
	return &response, nil
}

func (u *merchantTableUsecase) SoftDeleteTable(ctx context.Context, tenantID, tableID string) error {
	err := u.repo.SoftDelete(ctx, tenantID, tableID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.New("NOT_FOUND", "Meja tidak ditemukan", http.StatusNotFound, nil)
		}
		return apperrors.New("INTERNAL_ERROR", "Gagal menghapus meja", http.StatusInternalServerError, err)
	}
	return nil
}

func toTableResponse(entity *domain.DiningTableEntity) dto.TableResponse {
	return dto.TableResponse{
		ID:        entity.ID,
		TenantID:  entity.TenantID,
		TableName: entity.Name,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: entity.DeletedAt,
	}
}
