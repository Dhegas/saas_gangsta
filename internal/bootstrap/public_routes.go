package bootstrap

import (
	categoryhttp "github.com/dhegas/saas_gangsta/internal/domains/category/delivery/http"
	categoryrepo "github.com/dhegas/saas_gangsta/internal/domains/category/repository"
	categoryusecase "github.com/dhegas/saas_gangsta/internal/domains/category/usecase"
	menuhttp "github.com/dhegas/saas_gangsta/internal/domains/menu/delivery/http"
	menurepo "github.com/dhegas/saas_gangsta/internal/domains/menu/repository"
	menuusecase "github.com/dhegas/saas_gangsta/internal/domains/menu/usecase"
	tablehttp "github.com/dhegas/saas_gangsta/internal/domains/table/delivery/http"
	tablerepo "github.com/dhegas/saas_gangsta/internal/domains/table/repository"
	tableusecase "github.com/dhegas/saas_gangsta/internal/domains/table/usecase"
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

	// Category
	categoryRepo := categoryrepo.NewPublicCategoryRepository(db)
	categoryUC := categoryusecase.NewPublicCategoryUsecase(categoryRepo)
	publicCategoryHandler := categoryhttp.NewPublicCategoryHandler(categoryUC)

	// Menu
	menuRepo := menurepo.NewPublicMenuRepository(db)
	menuUC := menuusecase.NewPublicMenuUsecase(menuRepo)
	publicMenuHandler := menuhttp.NewPublicMenuHandler(menuUC)

	// Table
	tableRepo := tablerepo.NewPublicTableRepository(db)
	tableUC := tableusecase.NewPublicTableUsecase(tableRepo)
	publicTableHandler := tablehttp.NewPublicTableHandler(tableUC)

	publicRoutes := api.Group("/public")
	{
		publicRoutes.GET("/tenants", tenantHandler.GetPublicTenantList)
		publicRoutes.GET("/tenants/:slug", tenantHandler.GetPublicTenantDetail)

		// Tenant-resolved routes
		resolvedRoutes := publicRoutes.Group("/tenant/:tenantSlug", middleware.TenantResolver(db))
		{
			resolvedRoutes.GET("/categories", publicCategoryHandler.GetPublicCategories)
			resolvedRoutes.GET("/menus", publicMenuHandler.GetPublicMenus)
			resolvedRoutes.GET("/tables", publicTableHandler.GetPublicTables)
			resolvedRoutes.GET("/dining-tables", publicTableHandler.GetPublicTables)
		}
	}
}

