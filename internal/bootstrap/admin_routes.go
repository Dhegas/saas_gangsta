package bootstrap

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	// 1. Import Modul Tenant
	tenantHttp "github.com/dhegas/saas_gangsta/internal/modules/adminTenant/delivery/http"
	tenantRepo "github.com/dhegas/saas_gangsta/internal/modules/adminTenant/repository"
	tenantUc "github.com/dhegas/saas_gangsta/internal/modules/adminTenant/usecase"

	// 2. Import Modul Dashboard
	dashboardHttp "github.com/dhegas/saas_gangsta/internal/modules/adminDashboard/delivery/http"
	dashboardRepo "github.com/dhegas/saas_gangsta/internal/modules/adminDashboard/repository"
	dashboardUc "github.com/dhegas/saas_gangsta/internal/modules/adminDashboard/usecase"

	// 3. Import Modul Subscription (TAMBAHAN BARU)
	subsHttp "github.com/dhegas/saas_gangsta/internal/modules/adminSubscription/delivery/http"
	subsRepo "github.com/dhegas/saas_gangsta/internal/modules/adminSubscription/repository"
	subsUc "github.com/dhegas/saas_gangsta/internal/modules/adminSubscription/usecase"
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
