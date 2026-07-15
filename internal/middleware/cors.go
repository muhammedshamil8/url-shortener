package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/muhammedshamil8/url-shortener/internal/config"
)

func CORS(cfg *config.Config) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: cfg.Server.AllowedOrigins,
		AllowMethods: []string{
			"GET",
			"POST",
			"DELETE",
		},

		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},

		ExposeHeaders: []string{
			"X-Request-ID",
		},

		AllowCredentials: false,

		MaxAge: 12 * time.Hour,
	})
}
