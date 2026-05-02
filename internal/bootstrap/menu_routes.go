package bootstrap

import (
	"github.com/dhegas/saas_gangsta/internal/config"
	menuhttp "github.com/dhegas/saas_gangsta/internal/domains/menu/delivery/http"
	menurepo "github.com/dhegas/saas_gangsta/internal/domains/menu/repository"
	menuusecase "github.com/dhegas/saas_gangsta/internal/domains/menu/usecase"
	"github.com/dhegas/saas_gangsta/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterMenuRoutes(api *gin.RouterGroup, cfg *config.Config, db *gorm.DB) {
	menuRepo := menurepo.NewPartnerMenuRepository(db)
	menuUC := menuusecase.NewPartnerMenuUsecase(menuRepo)
	menuHandler := menuhttp.NewPartnerMenuHandler(menuUC)

	// Customer / Public Routes (bisa baca tanpa TenantGuard / JWT Partner, tapi baca tenant_id dari query)
	customerMenuRoutes := api.Group("/menus")
	customerMenuRoutes.Use(
		middleware.JWTAuth(cfg),
		middleware.RoleGuards("PARTNER", "CUSTOMER", "ADMIN"),
	)
	customerMenuRoutes.GET("", menuHandler.GetAllMenus)
	customerMenuRoutes.GET("/:id", menuHandler.GetMenuByID)

	// Partner specific routes (bisa tulis)
	partnerMenuRoutes := api.Group("/menus")
	partnerMenuRoutes.Use(
		middleware.JWTAuth(cfg),
		middleware.RoleGuard("PARTNER"),
		middleware.TenantGuard(),
	)
	partnerMenuRoutes.POST("", menuHandler.CreateMenu)
	partnerMenuRoutes.PUT("/:id", menuHandler.UpdateMenu)
	partnerMenuRoutes.DELETE("/:id", menuHandler.SoftDeleteMenu)
	partnerMenuRoutes.PATCH("/:id/toggle-available", menuHandler.ToggleMenuAvailable)
}
