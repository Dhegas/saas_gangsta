package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/dhegas/saas_gangsta/internal/common/config"
	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/common/middleware"
	"github.com/dhegas/saas_gangsta/internal/common/response"
	"github.com/dhegas/saas_gangsta/internal/database"
	logpkg "github.com/dhegas/saas_gangsta/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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
		return nil, fmt.Errorf("bootstrap db: %w", err)
	}

	redisClient, err := database.ConnectRedis(cfg.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("bootstrap redis: %w", err)
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
	router.GET("/openapi.yaml", func(c *gin.Context) {
		c.File("./docs/openapi.yaml")
	})

	router.GET("/swagger", func(c *gin.Context) {
		c.File("./docs/swagger.html")
	})

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
