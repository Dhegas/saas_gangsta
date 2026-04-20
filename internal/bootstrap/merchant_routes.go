package bootstrap

import (
	"net/http"

	"github.com/dhegas/saas_gangsta/internal/config"
	"github.com/dhegas/saas_gangsta/internal/middleware"
	"github.com/dhegas/saas_gangsta/internal/common/response"
	authhttp "github.com/dhegas/saas_gangsta/internal/domains/user/auth/delivery/http"
	"github.com/gin-gonic/gin"
)

func RegisterMerchantRoutes(api *gin.RouterGroup, cfg *config.Config, authHandler *authhttp.AuthHandler) {
	merchantRoutes := api.Group("/merchant")
	merchantRoutes.Use(
		middleware.JWTAuth(cfg),
		middleware.RoleGuard("merchant"),
	)

	merchantRoutes.POST("/tenants", authHandler.CreateMerchantTenant)
	merchantRoutes.GET("/tenants", authHandler.ListMerchantTenants)

	merchantTenantScoped := merchantRoutes.Group("")
	merchantTenantScoped.Use(middleware.TenantGuard())
	merchantTenantScoped.GET("/me", func(c *gin.Context) {
		response.Success(c, http.StatusOK, "Merchant context valid", gin.H{
			"role":     "merchant",
			"tenantId": c.GetString("tenantId"),
		})
	})
}
