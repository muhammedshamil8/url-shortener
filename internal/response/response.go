package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muhammedshamil8/url-shortener/internal/middleware"
)

func Success(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{
		"status":     "success",
		"data":       data,
		"request_id": c.GetString(middleware.RequestIDKey),
	})
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"status":     "error",
		"message":    message,
		"request_id": c.GetString(middleware.RequestIDKey),
	})
}

func OK(c *gin.Context, data any) {
	Success(c, http.StatusOK, data)
}

func Created(c *gin.Context, data any) {
	Success(c, http.StatusCreated, data)
}
