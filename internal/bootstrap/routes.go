package bootstrap

import (
	"context"
	"net/http"
	"time"

	_ "github.com/dhegas/saas_gangsta/docs"
	"github.com/dhegas/saas_gangsta/internal/config"
	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/middleware"
	"github.com/dhegas/saas_gangsta/internal/common/response"
	"github.com/dhegas/saas_gangsta/internal/infrastructure/database"
	authhttp "github.com/dhegas/saas_gangsta/internal/domains/user/auth/delivery/http"
	authrepo "github.com/dhegas/saas_gangsta/internal/domains/user/auth/repository"
	authusecase "github.com/dhegas/saas_gangsta/internal/domains/user/auth/usecase"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func registerRoutes(router *gin.Engine, cfg *config.Config, db *gorm.DB, redisClient *redis.Client) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/health", func(c *gin.Context) {
		response.Success(c, http.StatusOK, "Service is healthy", gin.H{
			"status":    "ok",
			"service":   cfg.AppName,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	readinessHandler := func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		if err := database.IsReady(ctx, db); err != nil {
			apperrors.Write(c, apperrors.New("INTERNAL_ERROR", "Database is not ready", http.StatusInternalServerError, nil))
			return
		}

		if redisClient == nil {
			apperrors.Write(c, apperrors.New("INTERNAL_ERROR", "Redis is not configured or not reachable", http.StatusInternalServerError, nil))
			return
		}

		if err := database.IsRedisReady(ctx, redisClient); err != nil {
			apperrors.Write(c, apperrors.New("INTERNAL_ERROR", "Redis is not ready", http.StatusInternalServerError, nil))
			return
		}

		response.Success(c, http.StatusOK, "Service is ready", gin.H{
			"status":    "ok",
			"service":   cfg.AppName,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	}

	router.GET("/ready", readinessHandler)

	authRepository := authrepo.NewAuthRepository(db)
	authUC := authusecase.NewAuthUsecase(authRepository, cfg)
	authHandler := authhttp.NewAuthHandler(authUC)

	api := router.Group("/api/v1")
	{
		registerAuthRoutes(api, cfg, authHandler)
		RegisterCustomerRoutes(api, cfg)
		RegisterMerchantRoutes(api, cfg, authHandler)
		RegisterAdminRoutes(api, db)
		RegisterTenantProfileRoutes(api, cfg, db)

		api.GET("/health", func(c *gin.Context) {
			response.Success(c, http.StatusOK, "API is healthy", gin.H{
				"status":    "ok",
				"service":   cfg.AppName,
				"timestamp": time.Now().UTC().Format(time.RFC3339),
			})
		})
		api.GET("/ready", readinessHandler)
	}

	apiNoVersion := router.Group("/api")
	{
		RegisterTenantProfileRoutes(apiNoVersion, cfg, db)
	}
}

func registerAuthRoutes(api *gin.RouterGroup, cfg *config.Config, authHandler *authhttp.AuthHandler) {
	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/refresh", authHandler.Refresh)

		authProtected := authRoutes.Group("")
		authProtected.Use(middleware.JWTAuth(cfg))
		{
			authProtected.POST("/subscribe", middleware.RoleGuard("customer"), authHandler.Subscribe)
			authProtected.POST("/logout", authHandler.Logout)
			authProtected.GET("/me", authHandler.Me)
		}
	}
}
