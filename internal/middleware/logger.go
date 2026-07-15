package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muhammedshamil8/url-shortener/internal/logger"
	"github.com/muhammedshamil8/url-shortener/internal/response"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()
		requestID := c.GetString(response.RequestIDKey)
		logger.Log.Info(
			"request",
			"request_id", requestID,
			"path", c.Request.URL.Path,
			"method", c.Request.Method,
			"status", c.Writer.Status(),
			"ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
			"duration", time.Since(start),
		)
	}
}
