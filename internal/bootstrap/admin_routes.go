package bootstrap

import (
	"github.com/dhegas/saas_gangsta/internal/config"
	tenanthttp "github.com/dhegas/saas_gangsta/internal/domains/tenant/delivery/http"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant/repository"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant/usecase"
	"github.com/dhegas/saas_gangsta/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAdminRoutes(api *gin.RouterGroup, cfg *config.Config, db *gorm.DB) {
	tenantRepo := repository.NewAdminTenantRepository(db)
	tenantUsecase := usecase.NewAdminTenantUsecase(tenantRepo)
	tenantHandler := tenanthttp.NewAdminTenantHandler(tenantUsecase)

	adminRoutes := api.Group("/admin")
	adminRoutes.Use(
		middleware.JWTAuth(cfg),
		middleware.RoleGuard("ADMIN"),
	)

	adminRoutes.POST("/tenants", tenantHandler.CreateAdminTenant)
	adminRoutes.GET("/tenants", tenantHandler.ListAllTenants)
}
