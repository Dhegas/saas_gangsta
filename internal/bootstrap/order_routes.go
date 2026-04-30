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
	orderRepo := orderrepo.NewMerchantOrderRepository(db)
	orderUC := orderusecase.NewMerchantOrderUsecase(orderRepo)
	orderHandler := orderhttp.NewMerchantOrderHandler(orderUC)

	// Customer Management routes (per order)
	customerRepo := orderrepo.NewCustomerRepository(db)
	customerUC := orderusecase.NewCustomerUsecase(customerRepo)
	customerHandler := orderhttp.NewCustomerHandler(customerUC)

	// Customer Routes (untuk membuat order dari QR code)
	customerOrderRoutes := api.Group("/orders")
	customerOrderRoutes.Use(
		middleware.JWTAuth(cfg), // Tetap butuh auth (token login customer atau guest token), atur sesuai kebutuhan sistem
		middleware.RoleGuards("BASIC", "MITRA", "ADMIN"),
	)
	customerOrderRoutes.POST("", orderHandler.CreateOrder)
	customerOrderRoutes.GET("/:id", orderHandler.GetOrderByID) // Customer mungkin butuh melihat struk detailnya

	// Customer sub-resource routes: POST, GET, PUT /api/orders/:id/customer
	customerOrderRoutes.POST("/:id/customer", customerHandler.CreateCustomer)
	customerOrderRoutes.GET("/:id/customer", customerHandler.GetCustomer)
	customerOrderRoutes.PUT("/:id/customer", customerHandler.UpdateCustomer)

	// Mitra specific routes (untuk memanajemen pesanan masuk)
	mitraOrderRoutes := api.Group("/orders")
	mitraOrderRoutes.Use(
		middleware.JWTAuth(cfg),
		middleware.RoleGuard("MITRA"),
		middleware.TenantGuard(),
	)
	mitraOrderRoutes.GET("", orderHandler.GetAllOrders)
	mitraOrderRoutes.PATCH("/:id/status", orderHandler.UpdateOrderStatus)
	mitraOrderRoutes.DELETE("/:id", orderHandler.SoftDeleteOrder)
}
