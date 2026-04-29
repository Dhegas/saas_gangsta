package bootstrap

import (
	"github.com/dhegas/saas_gangsta/internal/config"
	tablehttp "github.com/dhegas/saas_gangsta/internal/domains/table/delivery/http"
	tablerepo "github.com/dhegas/saas_gangsta/internal/domains/table/repository"
	tableusecase "github.com/dhegas/saas_gangsta/internal/domains/table/usecase"
	"github.com/dhegas/saas_gangsta/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterTableRoutes(api *gin.RouterGroup, cfg *config.Config, db *gorm.DB) {
	tableRepo := tablerepo.NewMerchantTableRepository(db)
	tableUC := tableusecase.NewMerchantTableUsecase(tableRepo)
	tableHandler := tablehttp.NewMerchantTableHandler(tableUC)

	tableRoutes := api.Group("/dining-tables")
	tableRoutes.Use(
		middleware.JWTAuth(cfg),
		middleware.RoleGuard("MITRA"),
		middleware.TenantGuard(),
	)

	tableRoutes.POST("", tableHandler.CreateTable)
	tableRoutes.GET("", tableHandler.GetAllTables)
	tableRoutes.GET("/:id", tableHandler.GetTableByID)
	tableRoutes.GET("/:id/status", tableHandler.GetTableStatus)
	tableRoutes.PUT("/:id", tableHandler.UpdateTable)
	tableRoutes.DELETE("/:id", tableHandler.SoftDeleteTable)
}
