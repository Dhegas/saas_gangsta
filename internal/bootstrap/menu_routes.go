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

	menuRoutes := api.Group("/menus")
	menuRoutes.Use(
		middleware.JWTAuth(cfg),
		middleware.RoleGuard("MITRA"),
		middleware.TenantGuard(),
	)

	menuRoutes.POST("", menuHandler.CreateMenu)
	menuRoutes.GET("", menuHandler.GetAllMenus)
	menuRoutes.GET("/:id", menuHandler.GetMenuByID)
	menuRoutes.PUT("/:id", menuHandler.UpdateMenu)
	menuRoutes.DELETE("/:id", menuHandler.SoftDeleteMenu)
	menuRoutes.PATCH("/:id/toggle-available", menuHandler.ToggleMenuAvailable)
}
