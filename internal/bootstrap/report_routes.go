package bootstrap

import (
	"github.com/dhegas/saas_gangsta/internal/config"
	reporthttp "github.com/dhegas/saas_gangsta/internal/domains/report/delivery/http"
	reportrepo "github.com/dhegas/saas_gangsta/internal/domains/report/repository"
	reportusecase "github.com/dhegas/saas_gangsta/internal/domains/report/usecase"
	"github.com/dhegas/saas_gangsta/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterReportRoutes(api *gin.RouterGroup, cfg *config.Config, db *gorm.DB) {
	reportRepo := reportrepo.NewPartnerReportRepository(db)
	reportUC := reportusecase.NewPartnerReportUsecase(reportRepo)
	reportHandler := reporthttp.NewPartnerReportHandler(reportUC)

	reportRoutes := api.Group("/reports")
	reportRoutes.Use(
		middleware.JWTAuth(cfg),
		middleware.RoleGuard("PARTNER"),
		middleware.TenantGuard(),
	)

	reportRoutes.GET("/revenue", reportHandler.GetRevenue)
	reportRoutes.GET("/top-menus", reportHandler.GetTopMenus)
	reportRoutes.GET("/orders-by-table", reportHandler.GetOrdersByTable)
	reportRoutes.GET("/daily-summary", reportHandler.GetDailySummary)
}
