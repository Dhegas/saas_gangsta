package bootstrap

import (
	"net/http"

	"github.com/dhegas/saas_gangsta/internal/config"
	"github.com/dhegas/saas_gangsta/internal/middleware"
	"github.com/dhegas/saas_gangsta/internal/common/response"
	authhttp "github.com/dhegas/saas_gangsta/internal/domains/user/auth/delivery/http"
	"github.com/gin-gonic/gin"
)

func RegisterPartnerRoutes(api *gin.RouterGroup, cfg *config.Config, authHandler *authhttp.AuthHandler) {
	partnerRoutes := api.Group("/partner")
	partnerRoutes.Use(
		middleware.JWTAuth(cfg),
		middleware.RoleGuard("PARTNER"),
	)

	partnerRoutes.POST("/tenants", authHandler.CreatePartnerTenant)
	partnerRoutes.GET("/tenants", authHandler.ListPartnerTenants)

	partnerTenantScoped := partnerRoutes.Group("")
	partnerTenantScoped.Use(middleware.TenantGuard())
	partnerTenantScoped.GET("/me", func(c *gin.Context) {
		response.Success(c, http.StatusOK, "Partner context valid", gin.H{
			"role":     "PARTNER",
			"tenantId": c.GetString("tenantId"),
		})
	})
}
