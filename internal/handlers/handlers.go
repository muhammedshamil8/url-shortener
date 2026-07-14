package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/muhammedshamil8/url-shortener/internal/models"
	"github.com/muhammedshamil8/url-shortener/internal/response"
	"github.com/muhammedshamil8/url-shortener/internal/utils"
)

const (
	MaxRetries      = 5
	UniqueViolation = "23505"
)

type Handler struct {
    repo URLRepository
}

func New(repo URLRepository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler)HealthCheckHandler(c *gin.Context) {
	response.OK(c, gin.H{
		"message": "Welcome to URL Shortener Service",
	})
}

func (h *Handler)ShortenHandler(c *gin.Context) {
	var req models.ShortenRequest
	if err := c.BindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := utils.ValidateURL(req.URL); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid URL")
		return
	}

	for i := 0; i < MaxRetries; i++ {
		shortCode, err := utils.GenerateShortCode()
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "Failed to generate short code")
			return
		}
		id, err := h.repo.CreateShortURL(shortCode, req.URL)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == UniqueViolation {
				continue
			}
			response.Error(c, http.StatusInternalServerError, "Failed to create short url")
			return
		}
		response.Created(c, gin.H{
			"id":           id,
			"original_url": req.URL,
			"short_code":   shortCode,
			"short_url":    os.Getenv("BASE_URL") + shortCode,
		})
		return
	}

	response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to create short url after %d attempts", MaxRetries))

}

func (h *Handler)RedirectHandler(c *gin.Context) {
	code := c.Param("code")
	url, err := h.repo.GetURLByCode(code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Error(c, http.StatusNotFound, "URL not found")
			return
		}

		response.Error(c, http.StatusInternalServerError, "Database error")
		return
	}
	// if err := repository.IncrementClickCount(code); err != nil {
	// 	log.Println("Error incrementing click count:", err)
	// }
	c.Redirect(http.StatusSeeOther, url)
}

func (h *Handler)DeleteHandler(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid id")
		return
	}
	err = h.repo.DeleteURL(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Failed to delete url")
		return
	}
	response.OK(c, gin.H{
		"message": "URL deleted successfully",
	})
}

func (h *Handler)ListAllHandler(c *gin.Context) {
	urls, err := h.repo.GetAllURLs()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get all urls")
		return
	}
	response.OK(c, gin.H{
		"urls": urls,
	})
}
