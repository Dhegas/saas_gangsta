package bootstrap

import (
	"github.com/dhegas/saas_gangsta/internal/config"
	tenanthttp "github.com/dhegas/saas_gangsta/internal/domains/tenant/delivery/http"
	tenantrepo "github.com/dhegas/saas_gangsta/internal/domains/tenant/repository"
	tenantuc "github.com/dhegas/saas_gangsta/internal/domains/tenant/usecase"
	"github.com/dhegas/saas_gangsta/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterTenantProfileRoutes(api *gin.RouterGroup, cfg *config.Config, db *gorm.DB) {
	repo := tenantrepo.NewTenantProfileRepository(db)
	uc := tenantuc.NewTenantProfileUsecase(repo)
	handler := tenanthttp.NewTenantProfileHandler(uc)

	routes := api.Group("/tenant-profiles")
	routes.Use(
		middleware.JWTAuth(cfg),
		middleware.RoleGuard("MITRA"),
		middleware.TenantGuard(),
	)

	// Collection routes
	routes.POST("", handler.Create)
	routes.GET("", handler.List)

	// Single-resource routes
	routes.GET("/:id", handler.GetByID)
	routes.PUT("/:id", handler.Update)
	routes.DELETE("/:id", handler.Delete)
	routes.PATCH("/:id/toggle-active", handler.ToggleActive)
}
