package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/muhammedshamil8/url-shortener/internal/config"
	"github.com/muhammedshamil8/url-shortener/internal/handlers"
	"github.com/muhammedshamil8/url-shortener/internal/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Setup(r *gin.Engine, h *handlers.Handler, cfg *config.Config) {
	r.Use(gin.Recovery())
	r.Use(middleware.CORS(cfg))
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	rateLimiter := middleware.NewRateLimiter()
	r.Use(middleware.RateLimit(rateLimiter))

	jwtMiddleware := middleware.AuthMiddleware(cfg)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/:code", h.RedirectHandler)

	api := r.Group("/api/v1")
	{
		api.GET("/live", h.LiveHandler)
		api.GET("/ready", h.ReadyHandler)

		api.POST("/shorten", h.ShortenHandler)
	}

	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/register", h.RegisterHandler)
		auth.POST("/login", h.LoginHandler)
		auth.POST("/refresh", h.RefreshHandler)
	}

	authRoutes := r.Group("/api/v1", jwtMiddleware)
	{
		authRoutes.GET("/me", h.GetProfileHandler)
		authRoutes.GET("/my/urls", h.ListUserURLs)
		authRoutes.PUT("/my/urls/:id", h.UpdateURL)
		authRoutes.DELETE("/my/urls/:id", h.DeleteURL)
	}

	adminRoutes := r.Group("/api/v1/admin", jwtMiddleware, middleware.AdminOnly())
	{
		adminRoutes.GET("/urls", h.ListAllHandler)
		adminRoutes.DELETE("/urls/:id", h.DeleteHandler)
		adminRoutes.GET("/users", h.AdminListUsers)
		adminRoutes.DELETE("/users/:id", h.AdminDeleteUser)
	}
}
