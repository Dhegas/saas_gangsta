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

type partnerTableUsecase struct {
	repo domain.PartnerTableRepository
}

func NewPartnerTableUsecase(repo domain.PartnerTableRepository) domain.PartnerTableUsecase {
	return &partnerTableUsecase{repo: repo}
}

func (u *partnerTableUsecase) GetAllTables(ctx context.Context, tenantID string) ([]dto.TableResponse, error) {
	tables, err := u.repo.FindAllByTenant(ctx, tenantID)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil data meja", http.StatusInternalServerError)
	}

	result := make([]dto.TableResponse, 0, len(tables))
	for _, t := range tables {
		isOccupied, err := u.repo.CheckTableOccupied(ctx, t.ID)
		status := "kosong"
		if err == nil && isOccupied {
			status = "occupied"
		}

		res := toTableResponse(&t)
		res.Status = status
		result = append(result, res)
	}

	return result, nil
}

func (u *partnerTableUsecase) GetTableByID(ctx context.Context, tenantID, tableID string) (*dto.TableResponse, error) {
	table, err := u.repo.FindByIDAndTenant(ctx, tenantID, tableID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New("NOT_FOUND", "Meja tidak ditemukan", http.StatusNotFound)
		}
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil data meja", http.StatusInternalServerError)
	}

	isOccupied, err := u.repo.CheckTableOccupied(ctx, tableID)
	status := "kosong"
	if err == nil && isOccupied {
		status = "occupied"
	}

	response := toTableResponse(table)
	response.Status = status
	return &response, nil
}

func (u *partnerTableUsecase) GetTableStatus(ctx context.Context, tenantID, tableID string) (*dto.TableStatusResponse, error) {
	table, err := u.repo.FindByIDAndTenant(ctx, tenantID, tableID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New("NOT_FOUND", "Meja tidak ditemukan", http.StatusNotFound)
		}
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal memverifikasi meja", http.StatusInternalServerError)
	}

	isOccupied, err := u.repo.CheckTableOccupied(ctx, tableID)
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengecek status meja", http.StatusInternalServerError)
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

func (u *partnerTableUsecase) CreateTable(ctx context.Context, tenantID string, req dto.CreateTableRequest) (*dto.TableResponse, error) {
	exists, err := u.repo.CheckNameExists(ctx, tenantID, req.TableName, "")
	if err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal memvalidasi nama meja", http.StatusInternalServerError)
	}
	if exists {
		return nil, apperrors.New("CONFLICT", "Nama meja sudah ada", http.StatusConflict)
	}

	entity := &domain.DiningTableEntity{
		TenantID:  tenantID,
		Name:      req.TableName,
	}

	if err := u.repo.Create(ctx, entity); err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal menyimpan meja", http.StatusInternalServerError)
	}

	response := toTableResponse(entity)
	return &response, nil
}

func (u *partnerTableUsecase) UpdateTable(ctx context.Context, tenantID, tableID string, req dto.UpdateTableRequest) (*dto.TableResponse, error) {
	table, err := u.repo.FindByIDAndTenant(ctx, tenantID, tableID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.New("NOT_FOUND", "Meja tidak ditemukan", http.StatusNotFound)
		}
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal mengambil data meja", http.StatusInternalServerError)
	}

	if req.TableName != "" && req.TableName != table.Name {
		exists, err := u.repo.CheckNameExists(ctx, tenantID, req.TableName, tableID)
		if err != nil {
			return nil, apperrors.New("INTERNAL_ERROR", "Gagal memvalidasi nama meja", http.StatusInternalServerError)
		}
		if exists {
			return nil, apperrors.New("CONFLICT", "Nama meja sudah ada", http.StatusConflict)
		}
		table.Name = req.TableName
	}

	if err := u.repo.Update(ctx, table); err != nil {
		return nil, apperrors.New("INTERNAL_ERROR", "Gagal memperbarui meja", http.StatusInternalServerError)
	}

	response := toTableResponse(table)
	return &response, nil
}

func (u *partnerTableUsecase) SoftDeleteTable(ctx context.Context, tenantID, tableID string) error {
	err := u.repo.SoftDelete(ctx, tenantID, tableID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.New("NOT_FOUND", "Meja tidak ditemukan", http.StatusNotFound)
		}
		return apperrors.New("INTERNAL_ERROR", "Gagal menghapus meja", http.StatusInternalServerError)
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
