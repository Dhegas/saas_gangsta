package bootstrap

import (
	"github.com/dhegas/saas_gangsta/internal/config"
	categoryhttp "github.com/dhegas/saas_gangsta/internal/domains/category/delivery/http"
	categoryrepo "github.com/dhegas/saas_gangsta/internal/domains/category/repository"
	categoryusecase "github.com/dhegas/saas_gangsta/internal/domains/category/usecase"
	"github.com/dhegas/saas_gangsta/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterCategoryRoutes(api *gin.RouterGroup, cfg *config.Config, db *gorm.DB) {
	categoryRepo := categoryrepo.NewPartnerCategoryRepository(db)
	categoryUC := categoryusecase.NewPartnerCategoryUsecase(categoryRepo)
	categoryHandler := categoryhttp.NewPartnerCategoryHandler(categoryUC)

	categoryRoutes := api.Group("/categories")
	categoryRoutes.Use(
		middleware.JWTAuth(cfg),
		middleware.RoleGuard("PARTNER"),
		middleware.TenantGuard(),
	)

	categoryRoutes.POST("", categoryHandler.CreateCategory)
	categoryRoutes.GET("", categoryHandler.GetAllCategories)
	categoryRoutes.GET("/:id", categoryHandler.GetCategoryByID)
	categoryRoutes.PUT("/:id", categoryHandler.UpdateCategory)
	categoryRoutes.DELETE("/:id", categoryHandler.SoftDeleteCategory)
	categoryRoutes.PATCH("/:id/toggle-active", categoryHandler.ToggleCategoryActive)
	categoryRoutes.PATCH("/reorder", categoryHandler.ReorderCategories)
}
