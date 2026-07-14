// @title URL Shortener API
// @version 1.0
// @description A production-ready URL shortener built with Go and Gin.
// @BasePath /
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/muhammedshamil8/url-shortener/docs"
	"github.com/muhammedshamil8/url-shortener/internal/config"
	"github.com/muhammedshamil8/url-shortener/internal/database"
	"github.com/muhammedshamil8/url-shortener/internal/handlers"
	"github.com/muhammedshamil8/url-shortener/internal/logger"
	"github.com/muhammedshamil8/url-shortener/internal/repository"
	"github.com/muhammedshamil8/url-shortener/internal/routes"
	"github.com/muhammedshamil8/url-shortener/internal/server"
)

func main() {
	logger.Init()

	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			logger.Log.Warn(".env file not found")
		}
	}

	cfg := config.Load()

	if err := cfg.Validate(); err != nil {
		logger.Log.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	logger.Log.Info("Server starting", "environment", cfg.Env)

	db, err := database.InitDB(cfg.DB)
	if err != nil {
		logger.Log.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	repo := repository.New(db)

	h := handlers.New(repo, *cfg)

	if err := database.MigrateUserTable(db); err != nil {
		logger.Log.Error("Failed to migrate database", "error", err)
		os.Exit(1)
	}

	if err := database.MigrateUrlTable(db); err != nil {
		logger.Log.Error("Failed to migrate database", "error", err)
		os.Exit(1)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	if err := r.SetTrustedProxies(nil); err != nil {
		logger.Log.Error("Failed to set trusted proxies", "error", err)
		os.Exit(1)
	}

	routes.Setup(r, h, cfg)

	logger.Log.Info("Server running on port " + cfg.Server.Port)
	srv := server.New(cfg, r)
	go func() {
		if err := srv.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			logger.Log.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait for signal
	<-quit
	logger.Log.Info("Shutting down server...")

	// Create a timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Log.Info("Server exited properly")

}
