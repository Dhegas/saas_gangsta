package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	_ "github.com/dhegas/saas_gangsta/docs"
	"github.com/dhegas/saas_gangsta/internal/common/config"
	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/common/middleware"
	"github.com/dhegas/saas_gangsta/internal/common/response"
	"github.com/dhegas/saas_gangsta/internal/database"
	authhttp "github.com/dhegas/saas_gangsta/internal/modules/auth/delivery/http"
	authrepo "github.com/dhegas/saas_gangsta/internal/modules/auth/repository"
	authusecase "github.com/dhegas/saas_gangsta/internal/modules/auth/usecase"
	logpkg "github.com/dhegas/saas_gangsta/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

type App struct {
	Config *config.Config
	Logger *slog.Logger
	DB     *gorm.DB
	Redis  *redis.Client
	Router *gin.Engine
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	log := logpkg.New(cfg.AppEnv)

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Warn("bootstrap db: continuing without database", "error", err)
		db = nil
	}

	redisClient, err := database.ConnectRedis(cfg.RedisURL)
	if err != nil {
		log.Warn("bootstrap redis: continuing without redis", "error", err)
		redisClient = nil
	}

	router := gin.New()
	router.Use(
		middleware.CORS(cfg),
		middleware.Logger(log),
		middleware.Recovery(log),
	)

	registerRoutes(router, cfg, db, redisClient)

	return &App{
		Config: cfg,
		Logger: log,
		DB:     db,
		Redis:  redisClient,
		Router: router,
	}, nil
}

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

	api := router.Group("/api/v1")
	{
		authRepository := authrepo.NewAuthRepository(db)
		authUC := authusecase.NewAuthUsecase(authRepository, cfg)
		authHandler := authhttp.NewAuthHandler(authUC)

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

		customerRoutes := api.Group("/customer")
		customerRoutes.Use(
			middleware.JWTAuth(cfg),
			middleware.RoleGuard("customer"),
			middleware.TenantGuard(),
		)
		{
			customerRoutes.GET("/me", func(c *gin.Context) {
				response.Success(c, http.StatusOK, "Customer context valid", gin.H{
					"role":     "customer",
					"tenantId": c.GetString("tenantId"),
				})
			})
		}

		merchantRoutes := api.Group("/merchant")
		merchantRoutes.Use(
			middleware.JWTAuth(cfg),
			middleware.RoleGuard("merchant"),
		)
		{
			merchantRoutes.POST("/tenants", authHandler.CreateMerchantTenant)
			merchantRoutes.GET("/tenants", authHandler.ListMerchantTenants)

			merchantTenantScoped := merchantRoutes.Group("")
			merchantTenantScoped.Use(middleware.TenantGuard())
			merchantTenantScoped.GET("/me", func(c *gin.Context) {
				response.Success(c, http.StatusOK, "Merchant context valid", gin.H{
					"role":     "merchant",
					"tenantId": c.GetString("tenantId"),
				})
			})
		}

		RegisterAdminRoutes(api, db)

		api.GET("/health", func(c *gin.Context) {
			response.Success(c, http.StatusOK, "API is healthy", gin.H{
				"status":    "ok",
				"service":   cfg.AppName,
				"timestamp": time.Now().UTC().Format(time.RFC3339),
			})
		})
		api.GET("/ready", readinessHandler)
	}
}
