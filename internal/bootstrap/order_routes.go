package bootstrap

import (
	"github.com/dhegas/saas_gangsta/internal/config"
	orderhttp "github.com/dhegas/saas_gangsta/internal/domains/order/delivery/http"
	orderrepo "github.com/dhegas/saas_gangsta/internal/domains/order/repository"
	orderusecase "github.com/dhegas/saas_gangsta/internal/domains/order/usecase"
	authrepo "github.com/dhegas/saas_gangsta/internal/domains/user/auth/repository"
	"github.com/dhegas/saas_gangsta/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterOrderRoutes(api *gin.RouterGroup, cfg *config.Config, db *gorm.DB) {
	orderRepo := orderrepo.NewPartnerOrderRepository(db)
	authRepo := authrepo.NewAuthRepository(db)
	orderUC := orderusecase.NewPartnerOrderUsecase(orderRepo, authRepo, cfg)
	orderHandler := orderhttp.NewPartnerOrderHandler(orderUC)

	// Customer Order (Self-Order publik dari QR code / slug) menggunakan unified usecase
	custOrderHandler := orderhttp.NewCustomerOrderHandler(orderUC)

	// Register rute publik untuk customer membuat order
	publicTenantOrderRoutes := api.Group("/public/tenant/:tenantSlug", middleware.TenantResolver(db))
	publicTenantOrderRoutes.GET("/orders", custOrderHandler.GetPublicOrders)
	publicTenantOrderRoutes.GET("/orders/:orderId", custOrderHandler.GetOrderStatus)

	// Customer Routes (untuk membuat order dari QR code)
	customerOrderRoutes := api.Group("/orders")
	customerOrderRoutes.Use(
		middleware.JWTAuth(cfg), // Tetap butuh auth (token login customer atau guest token), atur sesuai kebutuhan sistem
		middleware.RoleGuards("CUSTOMER", "PARTNER", "ADMIN"),
	)
	customerOrderRoutes.POST("/tenant/:tenantSlug", middleware.TenantResolver(db), orderHandler.CreateOrder)
	customerOrderRoutes.GET("", middleware.TenantGuard(db), orderHandler.GetAllOrders)
	customerOrderRoutes.GET("/:id", middleware.TenantGuard(db), orderHandler.GetOrderByID) // Customer mungkin butuh melihat struk detailnya

	// Partner specific routes (untuk memanajemen pesanan masuk)
	partnerOrderRoutes := api.Group("/orders")
	partnerOrderRoutes.Use(
		middleware.JWTAuth(cfg),
		middleware.RoleGuard("PARTNER"),
		middleware.TenantGuard(db),
	)
	partnerOrderRoutes.PATCH("/:id/status", orderHandler.UpdateOrderStatus)
	partnerOrderRoutes.DELETE("/:id", orderHandler.SoftDeleteOrder)
}
