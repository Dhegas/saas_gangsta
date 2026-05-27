package bootstrap

import (
	"github.com/dhegas/saas_gangsta/internal/config"
	orderhttp "github.com/dhegas/saas_gangsta/internal/domains/order/delivery/http"
	orderrepo "github.com/dhegas/saas_gangsta/internal/domains/order/repository"
	orderusecase "github.com/dhegas/saas_gangsta/internal/domains/order/usecase"
	"github.com/dhegas/saas_gangsta/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterOrderRoutes(api *gin.RouterGroup, cfg *config.Config, db *gorm.DB) {
	orderRepo := orderrepo.NewPartnerOrderRepository(db)
	orderUC := orderusecase.NewPartnerOrderUsecase(orderRepo)
	orderHandler := orderhttp.NewPartnerOrderHandler(orderUC)

	// Customer Management routes (per order)
	customerRepo := orderrepo.NewCustomerRepository(db)
	customerUC := orderusecase.NewCustomerUsecase(customerRepo)
	customerHandler := orderhttp.NewCustomerHandler(customerUC)

	// Customer Order (Self-Order publik dari QR code / slug) menggunakan unified usecase
	custOrderHandler := orderhttp.NewCustomerOrderHandler(orderUC)

	// Register rute publik untuk customer membuat order
	publicTenantOrderRoutes := api.Group("/public/tenant/:tenantSlug", middleware.TenantResolver(db))
	publicTenantOrderRoutes.POST("/orders", custOrderHandler.CreateOrder)
	publicTenantOrderRoutes.GET("/orders", custOrderHandler.GetPublicOrders)
	publicTenantOrderRoutes.GET("/orders/:orderId", custOrderHandler.GetOrderStatus)

	// Customer Routes (untuk membuat order dari QR code)
	customerOrderRoutes := api.Group("/orders")
	customerOrderRoutes.Use(
		middleware.JWTAuth(cfg), // Tetap butuh auth (token login customer atau guest token), atur sesuai kebutuhan sistem
		middleware.RoleGuards("CUSTOMER", "PARTNER", "ADMIN"),
	)
	customerOrderRoutes.POST("/tenant/:tenantSlug", middleware.TenantResolver(db), orderHandler.CreateOrder)
	customerOrderRoutes.GET("/:id", orderHandler.GetOrderByID) // Customer mungkin butuh melihat struk detailnya

	// Customer sub-resource routes: POST, GET, PUT /api/orders/:id/customer
	customerOrderRoutes.POST("/:id/customer", customerHandler.CreateCustomer)
	customerOrderRoutes.GET("/:id/customer", customerHandler.GetCustomer)
	customerOrderRoutes.PUT("/:id/customer", customerHandler.UpdateCustomer)

	// Partner specific routes (untuk memanajemen pesanan masuk)
	partnerOrderRoutes := api.Group("/orders")
	partnerOrderRoutes.Use(
		middleware.JWTAuth(cfg),
		middleware.RoleGuard("PARTNER"),
		middleware.TenantGuard(db),
	)
	partnerOrderRoutes.GET("", orderHandler.GetAllOrders)
	partnerOrderRoutes.PATCH("/:id/status", orderHandler.UpdateOrderStatus)
	partnerOrderRoutes.DELETE("/:id", orderHandler.SoftDeleteOrder)
}
