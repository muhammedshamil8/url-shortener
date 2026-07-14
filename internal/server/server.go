package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muhammedshamil8/url-shortener/internal/config"
)

func New(cfg *config.Config, r *gin.Engine) *http.Server {
	return &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
}
