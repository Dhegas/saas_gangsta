package bootstrap

import (
	"net/http"

	"github.com/dhegas/saas_gangsta/internal/common/response"
	"github.com/dhegas/saas_gangsta/internal/config"
	tenanthttp "github.com/dhegas/saas_gangsta/internal/domains/tenant/delivery/http"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant/repository"
	"github.com/dhegas/saas_gangsta/internal/domains/tenant/usecase"
	"github.com/dhegas/saas_gangsta/internal/middleware"
	"github.com/dhegas/saas_gangsta/internal/infrastructure/storage"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterPartnerRoutes(api *gin.RouterGroup, cfg *config.Config, db *gorm.DB, imageService storage.ImageService) {
	tenantRepo := repository.NewPartnerTenantRepository(db)
	tenantUsecase := usecase.NewPartnerTenantUsecase(tenantRepo, imageService)
	tenantHandler := tenanthttp.NewPartnerTenantHandler(tenantUsecase)

	partnerRoutes := api.Group("/partner")
	partnerRoutes.Use(
		middleware.JWTAuth(cfg),
		middleware.RoleGuard("PARTNER"),
	)

	partnerRoutes.POST("/tenants", tenantHandler.CreatePartnerTenant)
	partnerRoutes.GET("/tenants", tenantHandler.ListPartnerTenants)
	partnerRoutes.DELETE("/tenants/:id", tenantHandler.SoftDeletePartnerTenant)

	partnerTenantScoped := partnerRoutes.Group("")
	partnerTenantScoped.Use(middleware.TenantGuard())
	partnerTenantScoped.GET("/me", func(c *gin.Context) {
		response.Success(c, http.StatusOK, "Partner context valid", gin.H{
			"role":     "PARTNER",
			"tenantId": c.GetString("tenantId"),
		})
	})
}
