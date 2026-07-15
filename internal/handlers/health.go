package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muhammedshamil8/url-shortener/internal/response"
)

// Health godoc
//
//	@Summary		Health Check
//	@Description	Check if the API is running
//	@Tags			Health
//	@Produce		json
//	@Success		200	{object}	models.SuccessResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/api/v1/live [get]
func (h *Handler) LiveHandler(c *gin.Context) {
	response.OK(c, gin.H{
		"status": "alive",
	})
}

// Ready godoc
//
//	@Summary		Ready Check
//	@Description	Check if App is ready to serve traffic
//	@Tags		Health
//	@Produce		json
//	@Success		200	{object}	models.SuccessResponse
//	@Failure		503	{object}	models.ErrorResponse
//	@Router		/api/v1/ready [get]
func (h *Handler) ReadyHandler(c *gin.Context) {
	if err := h.repo.Health(); err != nil {
		response.Error(
			c,
			http.StatusServiceUnavailable,
			"database is unavailable",
		)
		return
	}
	response.OK(c, gin.H{
		"status": "ready",
	})
}
