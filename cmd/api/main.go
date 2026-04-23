package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dhegas/saas_gangsta/internal/bootstrap"
)

// @title saas_gangsta Backend API
// @version 1.0
// @description Backend API Service untuk platform SaaS POS & Self-Order UMKM kuliner.
// @BasePath /api/v1
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	app, err := bootstrap.New()
	if err != nil {
		log.Fatalf("bootstrap app: %v", err)
	}

	server := &http.Server{
		Addr:              ":" + app.Config.AppPort,
		Handler:           app.Router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		app.Logger.Info("server started", "port", app.Config.AppPort, "env", app.Config.AppEnv)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.Logger.Error("server crashed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		app.Logger.Error("server shutdown error", "error", err)
	}

	if app.Redis != nil {
		if err := app.Redis.Close(); err != nil {
			app.Logger.Error("redis shutdown error", "error", err)
		}
	}

	if app.DB != nil {
		sqlDB, err := app.DB.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				app.Logger.Error("database shutdown error", "error", err)
			}
		}
	}

	app.Logger.Info("server stopped")
}
