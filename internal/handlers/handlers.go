package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/muhammedshamil8/url-shortener/internal/config"
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
	cfg  config.Config
}

func New(repo URLRepository, cfg config.Config) *Handler {
	return &Handler{repo: repo, cfg: cfg}
}

// Shorten godoc
//
//		@Summary	Shorten a URL
//		@Description	Create a short URL from a long URL
//		@Tags	URLs
//		@Accept	json
//		@Produce	json
//	 @Param	request	body	models.ShortenRequest	true	"Request body"
//		@Success	200	{object}	models.SuccessResponse
//		@Failure	400	{object}	models.ErrorResponse
//		@Failure	500	{object}	models.ErrorResponse
//		@Router	/shorten [post]
func (h *Handler) ShortenHandler(c *gin.Context) {
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
			"short_url":    h.cfg.Server.BaseURL + "/" + shortCode,
		})
		return
	}

	response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to create short url after %d attempts", MaxRetries))

}

// Redirect godoc
//
//		@Summary	Redirect to original URL
//		@Description	Redirect to original URL based on short code
//		@Tags	URLs
//		@Produce	json
//		@Param	code	path	string	true	"Short code"
//	 @Success 303 {string} string "Redirect"
//		@Failure	404	{object}	models.ErrorResponse
//		@Failure	500	{object}	models.ErrorResponse
//		@Router	/{code} [get]
func (h *Handler) RedirectHandler(c *gin.Context) {
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

// Delete godoc
//
//	@Summary	Delete a URL
//	@Description	Delete a URL based on id
//	@Tags	URLs
//	@Produce	json
//	@Param	id	path	int	true	"URL ID"
//	@Success	200	{object}	models.SuccessResponse
//	@Failure	400	{object}	models.ErrorResponse
//	@Failure	404	{object}	models.ErrorResponse
//	@Failure	500	{object}	models.ErrorResponse
//	@Router	/{id} [delete]
func (h *Handler) DeleteHandler(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid id")
		return
	}
	err = h.repo.DeleteURL(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Error(c, http.StatusNotFound, "URL not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to delete url")
		return
	}
	response.OK(c, gin.H{
		"message": "URL deleted successfully",
	})
}

// ListAll godoc
//
//	@Summary	List all URLs
//	@Description	List all URLs
//	@Tags	URLs
//	@Produce	json
//	@Success	200	{object}	models.SuccessResponse
//	@Failure	404	{object}	models.ErrorResponse
//	@Failure	500	{object}	models.ErrorResponse
//	@Router	/urls/all [get]
func (h *Handler) ListAllHandler(c *gin.Context) {
	var opts models.ListOptions
	if err := c.ShouldBindQuery(&opts); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}
	urls, err := h.repo.GetAllURLs(opts)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get all urls")
		return
	}
	response.OK(c, gin.H{
		"urls": urls,
	})
}
