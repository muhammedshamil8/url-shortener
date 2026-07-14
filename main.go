package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/muhammedshamil8/url-shortener/internal/config"
	"github.com/muhammedshamil8/url-shortener/internal/database"
	"github.com/muhammedshamil8/url-shortener/internal/handlers"
	"github.com/muhammedshamil8/url-shortener/internal/logger"
	"github.com/muhammedshamil8/url-shortener/internal/middleware"
	"github.com/muhammedshamil8/url-shortener/internal/repository"
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

	h := handlers.New(repo)

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
	if err := r.Run(":" + port); err != nil {
		logger.Log.Error("Failed to start server", "error", err)
		os.Exit(1)
	}

}
