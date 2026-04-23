package bootstrap

import (
	"net/http"

	"github.com/dhegas/saas_gangsta/internal/config"
	"github.com/dhegas/saas_gangsta/internal/middleware"
	"github.com/dhegas/saas_gangsta/internal/common/response"
	"github.com/gin-gonic/gin"
)

func RegisterCustomerRoutes(api *gin.RouterGroup, cfg *config.Config) {
	customerRoutes := api.Group("/customer")
	customerRoutes.Use(
		middleware.JWTAuth(cfg),
		middleware.RoleGuard("BASIC"),
		middleware.TenantGuard(),
	)

	customerRoutes.GET("/me", func(c *gin.Context) {
		response.Success(c, http.StatusOK, "Customer context valid", gin.H{
			"role":     "BASIC",
			"tenantId": c.GetString("tenantId"),
		})
	})
}
