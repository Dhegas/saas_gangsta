package bootstrap

import (
	"context"
	"net/http"
	"time"

	docs "github.com/dhegas/saas_gangsta/docs"
	"github.com/dhegas/saas_gangsta/internal/common/cache"
	apperrors "github.com/dhegas/saas_gangsta/internal/common/errors"
	"github.com/dhegas/saas_gangsta/internal/common/response"
	"github.com/dhegas/saas_gangsta/internal/config"
	tenantrepo "github.com/dhegas/saas_gangsta/internal/domains/tenant/repository"
	"github.com/dhegas/saas_gangsta/internal/domains/user/auth"
	authhttp "github.com/dhegas/saas_gangsta/internal/domains/user/auth/delivery/http"
	authrepo "github.com/dhegas/saas_gangsta/internal/domains/user/auth/repository"
	authusecase "github.com/dhegas/saas_gangsta/internal/domains/user/auth/usecase"
	userhttp "github.com/dhegas/saas_gangsta/internal/domains/user/management/delivery/http"
	userrepo "github.com/dhegas/saas_gangsta/internal/domains/user/management/repository"
	userusecase "github.com/dhegas/saas_gangsta/internal/domains/user/management/usecase"
	"github.com/dhegas/saas_gangsta/internal/infrastructure/database"
	"github.com/dhegas/saas_gangsta/internal/infrastructure/websocket"
	"github.com/dhegas/saas_gangsta/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func registerRoutes(router *gin.Engine, cfg *config.Config, db *gorm.DB, redisClient *redis.Client, wsHub *websocket.Hub) {
	docs.SwaggerInfo.BasePath = "/api/v1"

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
			apperrors.Write(c, apperrors.New("INTERNAL_ERROR", "Database is not ready", http.StatusInternalServerError))
			return
		}

		// Redis is commented out, so we skip Redis readiness checks
		// if redisClient == nil {
		// 	apperrors.Write(c, apperrors.New("INTERNAL_ERROR", "Redis is not configured or not reachable", http.StatusInternalServerError))
		// 	return
		// }

		// if err := database.IsRedisReady(ctx, redisClient); err != nil {
		// 	apperrors.Write(c, apperrors.New("INTERNAL_ERROR", "Redis is not ready", http.StatusInternalServerError))
		// 	return
		// }

		response.Success(c, http.StatusOK, "Service is ready", gin.H{
			"status":    "ok",
			"service":   cfg.AppName,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	}

	router.GET("/ready", readinessHandler)

	// Register WebSocket route
	router.GET("/ws", func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
			return
		}

		claims, err := auth.ParseAccessToken(token, cfg.JWTSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		conn, err := websocket.Upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		var clientID string
		if claims.Role == "CUSTOMER" {
			clientID = claims.Subject
		} else if claims.Role == "PARTNER" {
			// Check if tenant_id query parameter is provided to support branch-specific WebSocket registration
			reqTenantID := c.Query("tenant_id")
			if reqTenantID == "" {
				reqTenantID = c.Query("tenantId")
			}

			if reqTenantID != "" {
				// Verify that this partner user owns/has access to this tenant_id
				var count int64
				err := db.Table("tenants").
					Where("id = NULLIF(?, '')::uuid AND user_id = NULLIF(?, '')::uuid AND deleted_at IS NULL", reqTenantID, claims.Subject).
					Count(&count).Error
				if err == nil && count > 0 {
					clientID = reqTenantID
				} else {
					// Unauthorized access to this tenant ID or invalid tenant ID
					_ = conn.WriteJSON(gin.H{"error": "Unauthorized or invalid tenant access"})
					conn.Close()
					return
				}
			} else {
				clientID = claims.TenantID
			}
		} else {
			clientID = claims.Subject
		}

		if clientID == "" {
			conn.Close()
			return
		}

		wsHub.Register(clientID, conn)

		go func() {
			defer func() {
				wsHub.Unregister(clientID, conn)
				conn.Close()
			}()
			for {
				_, _, err := conn.ReadMessage()
				if err != nil {
					break
				}
			}
		}()
	})

	authRepository := authrepo.NewAuthRepository(db)
	authUC := authusecase.NewAuthUsecase(authRepository, cfg)
	authHandler := authhttp.NewAuthHandler(authUC)
	tenantRepo := tenantrepo.NewAdminTenantRepository(db)
	userRepository := userrepo.NewUserRepository(db)
	userUC := userusecase.NewUserUsecase(userRepository, tenantRepo)
	userHandler := userhttp.NewUserHandler(userUC)

	localCache := cache.NewLocalCache()

	api := router.Group("/api/v1")
	{
		registerAuthRoutes(api, cfg, authHandler)
		registerUserRoutes(api, cfg, db, userHandler)
		RegisterPublicRoutes(api, cfg, db, localCache)
		RegisterCustomerRoutes(api, cfg, db, localCache, wsHub)
		RegisterPartnerRoutes(api, cfg, db, localCache, wsHub)
		RegisterAdminRoutes(api, cfg, db, localCache)

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

func registerAuthRoutes(api *gin.RouterGroup, cfg *config.Config, authHandler *authhttp.AuthHandler) {
	authRoutes := api.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/refresh", authHandler.Refresh)

		authProtected := authRoutes.Group("")
		authProtected.Use(middleware.JWTAuth(cfg))
		{
			authProtected.POST("/logout", authHandler.Logout)
			authProtected.GET("/me", authHandler.Me)
		}
	}
}

func registerUserRoutes(api *gin.RouterGroup, cfg *config.Config, db *gorm.DB, userHandler *userhttp.UserHandler) {
	userRoutes := api.Group("/users")
	userRoutes.Use(middleware.JWTAuth(cfg), middleware.TenantGuard(db), middleware.RoleGuards("PARTNER", "ADMIN"))
	{
		userRoutes.GET("", userHandler.ListUsers)
		userRoutes.GET("/:id", userHandler.GetUserDetail)
		userRoutes.PUT("/:id", userHandler.UpdateUser)
		userRoutes.DELETE("/:id", userHandler.SoftDeleteUser)
		userRoutes.PATCH("/:id/toggle-active", userHandler.ToggleActiveUser)
	}
}
