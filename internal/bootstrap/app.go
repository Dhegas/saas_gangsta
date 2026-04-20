package bootstrap

import (
	"fmt"
	"log/slog"

	"github.com/dhegas/saas_gangsta/internal/config"
	"github.com/dhegas/saas_gangsta/internal/middleware"
	"github.com/dhegas/saas_gangsta/internal/infrastructure/database"
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
