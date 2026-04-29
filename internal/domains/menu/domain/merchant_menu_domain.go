package domain

import (
	"context"
	"time"

	"github.com/dhegas/saas_gangsta/internal/domains/menu/dto"
)

// MerchantMenuUsecase kontrak logika bisnis
type MerchantMenuUsecase interface {
	GetAllMenus(ctx context.Context, tenantID string, filter dto.MenuFilterParams) ([]dto.MenuResponse, error)
	GetMenuByID(ctx context.Context, tenantID, menuID string) (*dto.MenuResponse, error)
	CreateMenu(ctx context.Context, tenantID string, req dto.CreateMenuRequest) (*dto.MenuResponse, error)
	UpdateMenu(ctx context.Context, tenantID, menuID string, req dto.UpdateMenuRequest) (*dto.MenuResponse, error)
	SoftDeleteMenu(ctx context.Context, tenantID, menuID string) error
	ToggleMenuAvailable(ctx context.Context, tenantID, menuID string, isAvailable bool) error
}

// MerchantMenuRepository kontrak interaksi database
type MerchantMenuRepository interface {
	FindAllByTenant(ctx context.Context, tenantID string, filter dto.MenuFilterParams) ([]MenuEntity, error)
	FindByIDAndTenant(ctx context.Context, tenantID, menuID string) (*MenuEntity, error)
	Create(ctx context.Context, entity *MenuEntity) error
	Update(ctx context.Context, entity *MenuEntity) error
	SoftDelete(ctx context.Context, tenantID, menuID string) error
	UpdateAvailableStatus(ctx context.Context, tenantID, menuID string, isAvailable bool) error
	CheckNameExists(ctx context.Context, tenantID, name string, excludeID string) (bool, error)
	CategoryExists(ctx context.Context, tenantID, categoryID string) (bool, error)
}

// MenuEntity merepresentasikan tabel menus
type MenuEntity struct {
	ID          string     `gorm:"primaryKey;default:gen_random_uuid()"`
	TenantID    string     `gorm:"index;not null"`
	CategoryID  *string    `gorm:"index"`
	Name        string     `gorm:"not null"`
	Description string
	Price       float64    `gorm:"type:numeric(12,2);not null"`
	ImageURL    string
	IsAvailable bool       `gorm:"not null;default:true"`
	CreatedAt   time.Time  `gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime"`
	DeletedAt   *time.Time `gorm:"index"`
}

func (MenuEntity) TableName() string {
	return "menus"
}
