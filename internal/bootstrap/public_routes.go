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

func RegisterPublicRoutes(api *gin.RouterGroup, cfg *config.Config, db *gorm.DB) {
	tenantRepo := repository.NewPublicTenantRepository(db)
	tenantUsecase := usecase.NewPublicTenantUsecase(tenantRepo)
	tenantHandler := tenanthttp.NewPublicTenantHandler(tenantUsecase)
	resourcesHandler := tenanthttp.NewPublicResourcesHandler(db)

	publicRoutes := api.Group("/public")
	{
		publicRoutes.GET("/tenants", tenantHandler.GetPublicTenantList)
		publicRoutes.GET("/tenants/:slug", tenantHandler.GetPublicTenantDetail)

		// Tenant-resolved routes
		resolvedRoutes := publicRoutes.Group("/tenant/:tenantSlug", middleware.TenantResolver(db))
		{
			resolvedRoutes.GET("/categories", resourcesHandler.GetPublicCategories)
			resolvedRoutes.GET("/menus", resourcesHandler.GetPublicMenus)
		}
	}
}
