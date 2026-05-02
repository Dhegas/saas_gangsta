package domain

import (
	"context"
	"time"

	"github.com/dhegas/saas_gangsta/internal/domains/table/dto"
)

// PartnerTableUsecase kontrak logika bisnis
type PartnerTableUsecase interface {
	GetAllTables(ctx context.Context, tenantID string) ([]dto.TableResponse, error)
	GetTableByID(ctx context.Context, tenantID, tableID string) (*dto.TableResponse, error)
	GetTableStatus(ctx context.Context, tenantID, tableID string) (*dto.TableStatusResponse, error)
	CreateTable(ctx context.Context, tenantID string, req dto.CreateTableRequest) (*dto.TableResponse, error)
	UpdateTable(ctx context.Context, tenantID, tableID string, req dto.UpdateTableRequest) (*dto.TableResponse, error)
	SoftDeleteTable(ctx context.Context, tenantID, tableID string) error
}

// PartnerTableRepository kontrak interaksi database
type PartnerTableRepository interface {
	FindAllByTenant(ctx context.Context, tenantID string) ([]DiningTableEntity, error)
	FindByIDAndTenant(ctx context.Context, tenantID, tableID string) (*DiningTableEntity, error)
	Create(ctx context.Context, entity *DiningTableEntity) error
	Update(ctx context.Context, entity *DiningTableEntity) error
	SoftDelete(ctx context.Context, tenantID, tableID string) error
	CheckNameExists(ctx context.Context, tenantID, tableName string, excludeID string) (bool, error)
	CheckTableOccupied(ctx context.Context, tableID string) (bool, error)
}

// DiningTableEntity merepresentasikan tabel dining_tables
type DiningTableEntity struct {
	ID        string     `gorm:"primaryKey;default:gen_random_uuid()"`
	TenantID  string     `gorm:"index;not null"`
	Name      string     `gorm:"column:table_name;not null"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
	DeletedAt *time.Time `gorm:"index"`
}

func (DiningTableEntity) TableName() string {
	return "dining_tables"
}
