package bootstrap

import (
	"github.com/dhegas/saas_gangsta/internal/config"
	"github.com/dhegas/saas_gangsta/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	// Modul Tenant
	tenantHttp "github.com/dhegas/saas_gangsta/internal/domains/tenant/delivery/http"
	tenantRepo "github.com/dhegas/saas_gangsta/internal/domains/tenant/repository"
	tenantUc   "github.com/dhegas/saas_gangsta/internal/domains/tenant/usecase"

	// Modul Dashboard
	dashboardHttp "github.com/dhegas/saas_gangsta/internal/domains/report/delivery/http"
	dashboardRepo "github.com/dhegas/saas_gangsta/internal/domains/report/repository"
	dashboardUc   "github.com/dhegas/saas_gangsta/internal/domains/report/usecase"

	// Modul Subscription
	subsHttp "github.com/dhegas/saas_gangsta/internal/domains/subscription/delivery/http"
	subsRepo "github.com/dhegas/saas_gangsta/internal/domains/subscription/repository"
	subsUc   "github.com/dhegas/saas_gangsta/internal/domains/subscription/usecase"
)

// RegisterAdminRoutes mendaftarkan semua endpoint khusus Admin di bawah /api/v1/admin/
// Seluruh route di sini dilindungi JWT + RoleGuard("admin").
//
// ─── PANDUAN DEVELOPER ────────────────────────────────────────────────────────
// Setiap kali membuat feature baru untuk Admin, daftarkan routenya di sini
// mengikuti pola yang sudah ada:
//
//  1. Inisialisasi repo, usecase, dan handler modul baru
//  2. Tambahkan route di bawah adminGroup dengan method & path yang sesuai
//     Contoh pola endpoint:
//       adminGroup.POST("/nama-resource",          handler.CreateX)
//       adminGroup.GET("/nama-resource",           handler.GetAllX)
//       adminGroup.GET("/nama-resource/:id",       handler.GetXByID)
//       adminGroup.PUT("/nama-resource/:id",       handler.UpdateX)
//       adminGroup.DELETE("/nama-resource/:id",    handler.SoftDeleteX)
//       adminGroup.PATCH("/nama-resource/:id/status", handler.UpdateXStatus)
// ──────────────────────────────────────────────────────────────────────────────
func RegisterAdminRoutes(apiV1 *gin.RouterGroup, cfg *config.Config, db *gorm.DB) {

	// Semua route admin dilindungi oleh JWT Auth + RoleGuard("admin")
	adminGroup := apiV1.Group("/admin")
	adminGroup.Use(
		middleware.JWTAuth(cfg),
		middleware.RoleGuard("ADMIN"),
	)

	// ══════════════════════════════════════════════════════════════════
	// 1. Tenant Management
	//    POST   /api/v1/admin/tenants          → Registrasi tenant baru
	//    GET    /api/v1/admin/tenants          → List semua tenant
	//    GET    /api/v1/admin/tenants/:id      → Detail tenant
	//    PUT    /api/v1/admin/tenants/:id      → Update tenant
	//    DELETE /api/v1/admin/tenants/:id      → Soft delete tenant
	//    PATCH  /api/v1/admin/tenants/:id/status → Update status tenant
	// ══════════════════════════════════════════════════════════════════
	adminTenantRepo    := tenantRepo.NewAdminTenantRepository(db)
	adminTenantUsecase := tenantUc.NewAdminTenantUsecase(adminTenantRepo)
	adminTenantHandler := tenantHttp.NewTenantHandler(adminTenantUsecase)

	adminGroup.POST("/tenants",              adminTenantHandler.CreateTenant)
	adminGroup.GET("/tenants",               adminTenantHandler.GetAllTenants)
	adminGroup.GET("/tenants/:id",           adminTenantHandler.GetTenantByID)
	adminGroup.PUT("/tenants/:id",           adminTenantHandler.UpdateTenant)
	adminGroup.DELETE("/tenants/:id",        adminTenantHandler.SoftDeleteTenant)
	adminGroup.PATCH("/tenants/:id/status",  adminTenantHandler.UpdateTenantStatus)

	// ══════════════════════════════════════════════════════════════════
	// 2. Admin Dashboard
	//    GET /api/v1/admin/dashboard → Overview platform
	// ══════════════════════════════════════════════════════════════════
	adminDashboardRepo    := dashboardRepo.NewAdminDashboardRepository(db)
	adminDashboardUsecase := dashboardUc.NewAdminDashboardUsecase(adminDashboardRepo)
	adminDashboardHandler := dashboardHttp.NewDashboardHandler(adminDashboardUsecase)

	adminDashboardHandler.RegisterRoutes(adminGroup)

	// ══════════════════════════════════════════════════════════════════
	// 3. Subscription Management
	//    (route didelegasikan ke RegisterRoutes milik handler)
	// ══════════════════════════════════════════════════════════════════
	adminSubsRepo    := subsRepo.NewAdminSubscriptionRepository(db)
	adminSubsUsecase := subsUc.NewAdminSubscriptionUsecase(adminSubsRepo)
	adminSubsHandler := subsHttp.NewSubscriptionHandler(adminSubsUsecase)

	adminSubsHandler.RegisterRoutes(adminGroup)
}
