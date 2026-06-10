package bootstrap

import (
	"net/http"

	"github.com/dhegas/saas_gangsta/internal/common/cache"
	"github.com/dhegas/saas_gangsta/internal/common/response"
	"github.com/dhegas/saas_gangsta/internal/config"
	categoryhttp "github.com/dhegas/saas_gangsta/internal/domains/category/delivery/http"
	categoryrepo "github.com/dhegas/saas_gangsta/internal/domains/category/repository"
	categoryusecase "github.com/dhegas/saas_gangsta/internal/domains/category/usecase"
	menuhttp "github.com/dhegas/saas_gangsta/internal/domains/menu/delivery/http"
	menurepo "github.com/dhegas/saas_gangsta/internal/domains/menu/repository"
	menuusecase "github.com/dhegas/saas_gangsta/internal/domains/menu/usecase"
	orderhttp "github.com/dhegas/saas_gangsta/internal/domains/order/delivery/http"
	orderrepo "github.com/dhegas/saas_gangsta/internal/domains/order/repository"
	orderusecase "github.com/dhegas/saas_gangsta/internal/domains/order/usecase"
	reporthttp "github.com/dhegas/saas_gangsta/internal/domains/report/delivery/http"
	reportrepo "github.com/dhegas/saas_gangsta/internal/domains/report/repository"
	reportusecase "github.com/dhegas/saas_gangsta/internal/domains/report/usecase"
	tablehttp "github.com/dhegas/saas_gangsta/internal/domains/table/delivery/http"
	tablerepo "github.com/dhegas/saas_gangsta/internal/domains/table/repository"
	tableusecase "github.com/dhegas/saas_gangsta/internal/domains/table/usecase"
	tenanthttp "github.com/dhegas/saas_gangsta/internal/domains/tenant/delivery/http"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant/repository"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant/usecase"
	authrepo "github.com/dhegas/saas_gangsta/internal/domains/user/auth/repository"
	wallethttp "github.com/dhegas/saas_gangsta/internal/domains/wallet/delivery/http"
	walletrepo "github.com/dhegas/saas_gangsta/internal/domains/wallet/repository"
	walletusecase "github.com/dhegas/saas_gangsta/internal/domains/wallet/usecase"
	"github.com/dhegas/saas_gangsta/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterPartnerRoutes(api *gin.RouterGroup, cfg *config.Config, db *gorm.DB, localCache *cache.LocalCache) {
	// === Dependency Init ===

	// Tenant
	tenantRepo := repository.NewPartnerTenantRepository(db)
	tenantUsecase := usecase.NewPartnerTenantUsecase(tenantRepo)
	tenantHandler := tenanthttp.NewPartnerTenantHandler(tenantUsecase)

	// Order
	orderRepo := orderrepo.NewPartnerOrderRepository(db)
	authRepo := authrepo.NewAuthRepository(db)
	orderUC := orderusecase.NewPartnerOrderUsecase(orderRepo, authRepo, cfg)
	orderHandler := orderhttp.NewPartnerOrderHandler(orderUC)

	// Menu
	menuRepo := menurepo.NewPartnerMenuRepository(db)
	menuUC := menuusecase.NewPartnerMenuUsecase(menuRepo)
	menuHandler := menuhttp.NewPartnerMenuHandler(menuUC)

	// Category
	categoryRepo := categoryrepo.NewPartnerCategoryRepository(db)
	categoryUC := categoryusecase.NewPartnerCategoryUsecase(categoryRepo)
	categoryHandler := categoryhttp.NewPartnerCategoryHandler(categoryUC)

	// Table
	tableRepo := tablerepo.NewPartnerTableRepository(db)
	tableUC := tableusecase.NewPartnerTableUsecase(tableRepo)
	tableHandler := tablehttp.NewPartnerTableHandler(tableUC)

	// Report
	reportRepo := reportrepo.NewPartnerReportRepository(db)
	reportUC := reportusecase.NewPartnerReportUsecase(reportRepo, localCache)
	reportHandler := reporthttp.NewPartnerReportHandler(reportUC)

	// Wallet
	walletRepo := walletrepo.NewWalletRepository(db)
	walletUC := walletusecase.NewPartnerWalletUsecase(walletRepo)
	walletHandler := wallethttp.NewPartnerWalletHandler(walletUC)

	// === Base Route Group ===
	partnerRoutes := api.Group("/partner")
	partnerRoutes.Use(
		middleware.JWTAuth(cfg),
		middleware.RoleGuard("PARTNER"),
	)

	// Ping / context validation
	partnerRoutes.GET("/me", func(c *gin.Context) {
		response.Success(c, http.StatusOK, "Partner context valid", gin.H{
			"role":     "PARTNER",
			"tenantId": "",
		})
	})

	// === Tenant Management (tanpa TenantGuard — partner mengelola tenant miliknya) ===
	partnerRoutes.POST("/tenants", tenantHandler.CreatePartnerTenant)
	partnerRoutes.GET("/tenants", tenantHandler.ListPartnerTenants)
	partnerRoutes.GET("/tenants/:id", tenantHandler.GetPartnerTenantByID)
	partnerRoutes.PUT("/tenants/:id", tenantHandler.UpdatePartnerTenant)
	partnerRoutes.DELETE("/tenants/:id", tenantHandler.SoftDeletePartnerTenant)

	// === Tenant-Scoped Routes (dengan TenantGuard) ===
	tenantScoped := partnerRoutes.Group("")
	tenantScoped.Use(middleware.TenantGuard(db))

	// Order Management
	tenantScoped.GET("/orders", orderHandler.GetAllOrders)
	tenantScoped.GET("/orders/:id", orderHandler.GetOrderByID)
	tenantScoped.POST("/orders", orderHandler.CreateOrder)
	tenantScoped.PATCH("/orders/:id/status", orderHandler.UpdateOrderStatus)
	tenantScoped.DELETE("/orders/:id", orderHandler.SoftDeleteOrder)

	// Menu Management
	tenantScoped.GET("/menus", menuHandler.GetAllMenus)
	tenantScoped.GET("/menus/:id", menuHandler.GetMenuByID)
	tenantScoped.POST("/menus", menuHandler.CreateMenu)
	tenantScoped.PUT("/menus/:id", menuHandler.UpdateMenu)
	tenantScoped.DELETE("/menus/:id", menuHandler.SoftDeleteMenu)
	tenantScoped.PATCH("/menus/:id/toggle-available", menuHandler.ToggleMenuAvailable)

	// Category Management
	tenantScoped.POST("/categories", categoryHandler.CreateCategory)
	tenantScoped.GET("/categories", categoryHandler.GetAllCategories)
	tenantScoped.GET("/categories/:id", categoryHandler.GetCategoryByID)
	tenantScoped.PUT("/categories/:id", categoryHandler.UpdateCategory)
	tenantScoped.DELETE("/categories/:id", categoryHandler.SoftDeleteCategory)
	tenantScoped.PATCH("/categories/:id/toggle-active", categoryHandler.ToggleCategoryActive)
	tenantScoped.PATCH("/categories/reorder", categoryHandler.ReorderCategories)

	// Table Management
	tenantScoped.POST("/dining-tables", tableHandler.CreateTable)
	tenantScoped.GET("/dining-tables", tableHandler.GetAllTables)
	tenantScoped.GET("/dining-tables/:id", tableHandler.GetTableByID)
	tenantScoped.GET("/dining-tables/:id/status", tableHandler.GetTableStatus)
	tenantScoped.PUT("/dining-tables/:id", tableHandler.UpdateTable)
	tenantScoped.DELETE("/dining-tables/:id", tableHandler.SoftDeleteTable)

	// Reports
	tenantScoped.GET("/reports/revenue", reportHandler.GetRevenue)
	tenantScoped.GET("/reports/top-menus", reportHandler.GetTopMenus)
	tenantScoped.GET("/reports/orders-by-table", reportHandler.GetOrdersByTable)
	tenantScoped.GET("/reports/daily-summary", reportHandler.GetDailySummary)

	// Wallet — tidak perlu TenantGuard karena scope by userID dari JWT
	walletHandler.RegisterRoutes(partnerRoutes)
}
