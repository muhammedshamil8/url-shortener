package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const RequestIDKey = "request_id"

func Success(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{
		"status":     "success",
		"data":       data,
		"request_id": c.GetString(RequestIDKey),
	})
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"status":     "error",
		"message":    message,
		"request_id": c.GetString(RequestIDKey),
	})
}

func OK(c *gin.Context, data any) {
	Success(c, http.StatusOK, data)
}

func Created(c *gin.Context, data any) {
	Success(c, http.StatusCreated, data)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message)
}

func InternalServerError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message)
}

func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

func UnprocessableEntity(c *gin.Context, message string) {
	Error(c, http.StatusUnprocessableEntity, message)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message)
}
