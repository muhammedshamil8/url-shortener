package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/muhammedshamil8/url-shortener/internal/auth"
	"github.com/muhammedshamil8/url-shortener/internal/config"
	"github.com/muhammedshamil8/url-shortener/internal/models"
	"github.com/muhammedshamil8/url-shortener/internal/response"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			response.Unauthorized(c, "Authorization token not provided")
			c.Abort()
			return
		}

		const bearerPrefix = "Bearer "
		if len(token) < len(bearerPrefix) || token[:len(bearerPrefix)] != bearerPrefix {
			response.Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}
		token = token[len(bearerPrefix):]
		claims, err := auth.ValidateToken(token, cfg.JWT.AccessTokenSecret)
		if err != nil {
			response.Unauthorized(c, "Invalid token")
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}

func GetClaims(c *gin.Context) *models.Claims {
	return c.MustGet("claims").(*models.Claims)
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := GetClaims(c)

		if claims.Role != "admin" {
			response.Forbidden(c, "Forbidden")
			c.Abort()
			return
		}

		c.Next()
	}
}
