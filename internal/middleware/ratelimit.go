package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muhammedshamil8/url-shortener/internal/logger"
	"github.com/muhammedshamil8/url-shortener/internal/response"
)

func RateLimit(rl *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		visitor := rl.getVisitor(ip)

		if !visitor.limiter.Allow() {
			response.Error(
				c,
				http.StatusTooManyRequests,
				"rate limit exceeded",
			)

			logger.Log.Warn(
				"rate limit exceeded",
				"ip", ip,
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
				"user_agent", c.Request.UserAgent(),
				"request_id", c.Value("request_id"),
			)

			c.Abort()
			return
		}

		c.Next()
	}
}
