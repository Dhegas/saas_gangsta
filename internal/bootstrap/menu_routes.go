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
	menuRepo := menurepo.NewMerchantMenuRepository(db)
	menuUC := menuusecase.NewMerchantMenuUsecase(menuRepo)
	menuHandler := menuhttp.NewMerchantMenuHandler(menuUC)

	// Customer / Public Routes (bisa baca tanpa TenantGuard / JWT MITRA, tapi baca tenant_id dari query)
	customerMenuRoutes := api.Group("/menus")
	customerMenuRoutes.Use(
		middleware.JWTAuth(cfg),
		middleware.RoleGuards("MITRA", "BASIC", "ADMIN"),
	)
	customerMenuRoutes.GET("", menuHandler.GetAllMenus)
	customerMenuRoutes.GET("/:id", menuHandler.GetMenuByID)

	// Mitra specific routes (bisa tulis)
	mitraMenuRoutes := api.Group("/menus")
	mitraMenuRoutes.Use(
		middleware.JWTAuth(cfg),
		middleware.RoleGuard("MITRA"),
		middleware.TenantGuard(),
	)
	mitraMenuRoutes.POST("", menuHandler.CreateMenu)
	mitraMenuRoutes.PUT("/:id", menuHandler.UpdateMenu)
	mitraMenuRoutes.DELETE("/:id", menuHandler.SoftDeleteMenu)
	mitraMenuRoutes.PATCH("/:id/toggle-available", menuHandler.ToggleMenuAvailable)
}
