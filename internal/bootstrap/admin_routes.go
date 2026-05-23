package bootstrap

import (
	"github.com/dhegas/saas_gangsta/internal/config"
	menuhttp "github.com/dhegas/saas_gangsta/internal/domains/menu/delivery/http"
	menurepo "github.com/dhegas/saas_gangsta/internal/domains/menu/repository"
	menuusecase "github.com/dhegas/saas_gangsta/internal/domains/menu/usecase"
	tenanthttp "github.com/dhegas/saas_gangsta/internal/domains/tenant/delivery/http"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant/repository"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant/usecase"
	userhttp "github.com/dhegas/saas_gangsta/internal/domains/user/management/delivery/http"
	userrepo "github.com/dhegas/saas_gangsta/internal/domains/user/management/repository"
	userusecase "github.com/dhegas/saas_gangsta/internal/domains/user/management/usecase"
	"github.com/dhegas/saas_gangsta/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAdminRoutes(api *gin.RouterGroup, cfg *config.Config, db *gorm.DB) {
	tenantRepo := repository.NewAdminTenantRepository(db)
	tenantUsecase := usecase.NewAdminTenantUsecase(tenantRepo)
	tenantHandler := tenanthttp.NewAdminTenantHandler(tenantUsecase)

	userRepo := userrepo.NewUserRepository(db)
	userUsecase := userusecase.NewUserUsecase(userRepo, tenantRepo)
	userHandler := userhttp.NewUserHandler(userUsecase)

	menuRepo := menurepo.NewPartnerMenuRepository(db)
	menuUC := menuusecase.NewPartnerMenuUsecase(menuRepo)
	adminMenuHandler := menuhttp.NewAdminMenuHandler(menuUC)

	adminRoutes := api.Group("/admin")
	adminRoutes.Use(
		middleware.JWTAuth(cfg),
		middleware.RoleGuard("ADMIN"),
	)

	adminRoutes.POST("/tenants", tenantHandler.CreateAdminTenant)
	adminRoutes.GET("/tenants", tenantHandler.ListAllTenants)
	adminRoutes.DELETE("/tenants/:id", tenantHandler.SoftDeleteTenant)
	adminRoutes.GET("/tenants/:id", tenantHandler.GetTenantByID)
	adminRoutes.GET("/tenants/users/:userId", tenantHandler.GetTenantsByUserID)

	adminRoutes.GET("/users", userHandler.ListAllUsersForAdmin)
	adminRoutes.GET("/users/:id", userHandler.GetUserDetailForAdmin)

	// Admin Menu Management — tenant ditentukan via header X-Tenant-ID
	adminRoutes.GET("/menus", adminMenuHandler.GetAllMenus)
	adminRoutes.GET("/menus/:id", adminMenuHandler.GetMenuByID)
	adminRoutes.POST("/menus", adminMenuHandler.CreateMenu)
	adminRoutes.PUT("/menus/:id", adminMenuHandler.UpdateMenu)
	adminRoutes.DELETE("/menus/:id", adminMenuHandler.SoftDeleteMenu)
	adminRoutes.PATCH("/menus/:id/toggle-available", adminMenuHandler.ToggleMenuAvailable)
}
