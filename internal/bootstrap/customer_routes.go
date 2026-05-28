package bootstrap

import (
	"net/http"

	"github.com/dhegas/saas_gangsta/internal/common/response"
	"github.com/dhegas/saas_gangsta/internal/config"
	orderhttp "github.com/dhegas/saas_gangsta/internal/domains/order/delivery/http"
	orderrepo "github.com/dhegas/saas_gangsta/internal/domains/order/repository"
	orderusecase "github.com/dhegas/saas_gangsta/internal/domains/order/usecase"
	authrepo "github.com/dhegas/saas_gangsta/internal/domains/user/auth/repository"
	"github.com/dhegas/saas_gangsta/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterCustomerRoutes(api *gin.RouterGroup, cfg *config.Config, db *gorm.DB) {
	// Order
	orderRepo := orderrepo.NewPartnerOrderRepository(db)
	authRepo := authrepo.NewAuthRepository(db)
	orderUC := orderusecase.NewPartnerOrderUsecase(orderRepo, authRepo, cfg)
	orderHandler := orderhttp.NewPartnerOrderHandler(orderUC)

	// Base group dengan autentikasi CUSTOMER
	customerRoutes := api.Group("/customer")
	customerRoutes.Use(
		middleware.JWTAuth(cfg),
		middleware.RoleGuard("CUSTOMER"),
	)

	// Ping / context validation
	customerRoutes.GET("/me", func(c *gin.Context) {
		response.Success(c, http.StatusOK, "Customer context valid", gin.H{
			"role": "CUSTOMER",
		})
	})

	// Order — dibuat oleh customer yang sudah login
	customerOrderRoutes := customerRoutes.Group("/orders")
	{
		// Membuat order baru (TenantResolver untuk resolve tenantId dari slug)
		customerOrderRoutes.POST("/tenant/:tenantSlug", middleware.TenantResolver(db), orderHandler.CreateOrder)

		// Melihat daftar order
		customerOrderRoutes.GET("", middleware.TenantGuard(db), orderHandler.GetAllOrders)

		// Melihat detail order (struk)
		customerOrderRoutes.GET("/:id", middleware.TenantGuard(db), orderHandler.GetOrderByID)
	}
}
