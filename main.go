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
	"github.com/muhammedshamil8/url-shortener/internal/config"
	"github.com/muhammedshamil8/url-shortener/internal/database"
	"github.com/muhammedshamil8/url-shortener/internal/handlers"
	"github.com/muhammedshamil8/url-shortener/internal/logger"
	"github.com/muhammedshamil8/url-shortener/internal/middleware"
	"github.com/muhammedshamil8/url-shortener/internal/repository"

	_ "github.com/muhammedshamil8/url-shortener/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	logger.Init()

	logger.Log.Info("Server starting")
	err := godotenv.Load()
	if err != nil {
		logger.Log.Error("Failed to load .env", "error", err)
		os.Exit(1)
	}
	cfg := config.Load()
	db, err := database.InitDB(cfg.DB)
	if err != nil {
		logger.Log.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	repo := repository.New(db)

	h := handlers.New(repo, *cfg)

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
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/health/api", h.HealthCheckHandler)
	r.POST("/shorten", h.ShortenHandler)
	r.GET("/urls/all", h.ListAllHandler)
	r.DELETE("/:id", h.DeleteHandler)
	r.GET("/:code", h.RedirectHandler)

	port := cfg.Server.Port
	if port == "" {
		logger.Log.Warn("APP_PORT is not set")
		logger.Log.Info("Automatic setting to 8080")
		port = "8080"
	}

	logger.Log.Info("Server running on port " + port)
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			logger.Log.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// -------------------------------------------------------
	// 2️⃣ Graceful shutdown (listen for Ctrl+C)
	// -------------------------------------------------------
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
