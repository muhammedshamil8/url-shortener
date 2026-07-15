package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/muhammedshamil8/url-shortener/internal/response"
)

const (
	RequestIDHeader = "X-Request-ID"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader(RequestIDHeader)

		if rid == "" {
			rid = uuid.NewString()
		}

		c.Set(response.RequestIDKey, rid)
		c.Writer.Header().Set(RequestIDHeader, rid)

		c.Next()
	}
}
