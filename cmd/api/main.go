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
	app.Logger.Info("server stopped")
}
