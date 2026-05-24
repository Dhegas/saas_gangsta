package repository

import (
	"context"

	tabledomain "github.com/dhegas/saas_gangsta/internal/domains/table/domain"
	"github.com/dhegas/saas_gangsta/internal/domains/table/dto"
	"gorm.io/gorm"
)

type publicTableRepository struct {
	db *gorm.DB
}

func NewPublicTableRepository(db *gorm.DB) tabledomain.PublicTableRepository {
	return &publicTableRepository{db: db}
}

func (r *publicTableRepository) FindPublicTables(ctx context.Context, tenantID string) ([]dto.PublicTableResponse, error) {
	var tables []dto.PublicTableResponse
	err := r.db.WithContext(ctx).
		Table("dining_tables dt").
		Select(`
			dt.id::text AS id, 
			dt.tenant_id::text AS tenant_id, 
			dt.table_name,
			CASE 
				WHEN EXISTS (
					SELECT 1 FROM orders o 
					WHERE o.dining_tables_id = dt.id 
					  AND o.status NOT IN ('COMPLETED', 'CANCELLED') 
					  AND o.deleted_at IS NULL
				) THEN 'occupied'
				ELSE 'kosong'
			END AS status
		`).
		Where("dt.tenant_id = ?", tenantID).
		Where("dt.deleted_at IS NULL").
		Order("dt.table_name ASC").
		Scan(&tables).Error

	return tables, err
}
