package bootstrap

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	// 1. Import Modul Tenant
	tenantHttp "github.com/dhegas/saas_gangsta/internal/domains/tenant/delivery/http"
	tenantRepo "github.com/dhegas/saas_gangsta/internal/domains/tenant/repository"
	tenantUc "github.com/dhegas/saas_gangsta/internal/domains/tenant/usecase"

	// 2. Import Modul Dashboard
	dashboardHttp "github.com/dhegas/saas_gangsta/internal/domains/report/delivery/http"
	dashboardRepo "github.com/dhegas/saas_gangsta/internal/domains/report/repository"
	dashboardUc "github.com/dhegas/saas_gangsta/internal/domains/report/usecase"

	// 3. Import Modul Subscription (TAMBAHAN BARU)
	subsHttp "github.com/dhegas/saas_gangsta/internal/domains/subscription/delivery/http"
	subsRepo "github.com/dhegas/saas_gangsta/internal/domains/subscription/repository"
	subsUc "github.com/dhegas/saas_gangsta/internal/domains/subscription/usecase"
)

// RegisterAdminRoutes adalah titik masuk untuk mendaftarkan semua endpoint khusus Admin.
func RegisterAdminRoutes(apiV1 *gin.RouterGroup, db *gorm.DB) {

	// ==========================================
	// 1. Modul Admin Tenant
	// ==========================================
	adminTenantRepo := tenantRepo.NewAdminTenantRepository(db)
	adminTenantUsecase := tenantUc.NewAdminTenantUsecase(adminTenantRepo)
	adminTenantHandler := tenantHttp.NewTenantHandler(adminTenantUsecase)
	adminTenantHandler.RegisterRoutes(apiV1)

	// ==========================================
	// 2. Modul Admin Dashboard
	// ==========================================
	adminDashboardRepo := dashboardRepo.NewAdminDashboardRepository(db)
	adminDashboardUsecase := dashboardUc.NewAdminDashboardUsecase(adminDashboardRepo)
	adminDashboardHandler := dashboardHttp.NewDashboardHandler(adminDashboardUsecase)
	adminDashboardHandler.RegisterRoutes(apiV1)

	// ==========================================
	// 3. Modul Admin Subscription (TAMBAHAN BARU)
	// ==========================================
	adminSubsRepo := subsRepo.NewAdminSubscriptionRepository(db)
	adminSubsUsecase := subsUc.NewAdminSubscriptionUsecase(adminSubsRepo)
	adminSubsHandler := subsHttp.NewSubscriptionHandler(adminSubsUsecase)
	adminSubsHandler.RegisterRoutes(apiV1)
}
