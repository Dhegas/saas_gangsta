package bootstrap

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/dhegas/saas_gangsta/internal/common/response"
	"github.com/dhegas/saas_gangsta/internal/config"
	"github.com/dhegas/saas_gangsta/internal/infrastructure/database"
	"github.com/dhegas/saas_gangsta/internal/middleware"
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

	ginMode := strings.ToLower(strings.Trim(strings.TrimSpace(os.Getenv("GIN_MODE")), "\"'"))
	if ginMode == gin.DebugMode || ginMode == gin.ReleaseMode || ginMode == gin.TestMode {
		gin.SetMode(ginMode)
	} else if cfg.AppEnv == "production" {
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
	router.HandleMethodNotAllowed = true
	router.Use(
		middleware.CORS(cfg),
		middleware.Logger(log),
		middleware.Recovery(log),
	)
	router.NoRoute(func(c *gin.Context) {
		response.Error(c, http.StatusNotFound, "Route not found", gin.H{
			"code": "NOT_FOUND",
			"path": c.Request.URL.Path,
		})
	})
	router.NoMethod(func(c *gin.Context) {
		response.Error(c, http.StatusMethodNotAllowed, "Method not allowed", gin.H{
			"code":   "METHOD_NOT_ALLOWED",
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
		})
	})

	registerRoutes(router, cfg, db, redisClient)

	return &App{
		Config: cfg,
		Logger: log,
		DB:     db,
		Redis:  redisClient,
		Router: router,
	}, nil
}
