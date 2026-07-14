package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/muhammedshamil8/url-shortener/internal/handlers"
	"github.com/muhammedshamil8/url-shortener/internal/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(r *gin.Engine, h *handlers.Handler) {
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/:code", h.RedirectHandler)

	api := r.Group("/api/v1")
	{
		api.GET("/health", h.HealthCheckHandler)
		api.POST("/shorten", h.ShortenHandler)
		api.GET("/urls", h.ListAllHandler)
		api.DELETE("/urls/:id", h.DeleteHandler)
	}
}
