package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/muhammedshamil8/url-shortener/internal/auth"
	"github.com/muhammedshamil8/url-shortener/internal/cache"
	"github.com/muhammedshamil8/url-shortener/internal/config"
	"github.com/muhammedshamil8/url-shortener/internal/logger"
	"github.com/muhammedshamil8/url-shortener/internal/models"
	"github.com/muhammedshamil8/url-shortener/internal/response"
	"github.com/muhammedshamil8/url-shortener/internal/utils"
)

const (
	MaxRetries      = 5
	UniqueViolation = "23505"
)

type Handler struct {
	repo Repository
	cache cache.Cache
	cfg  config.Config
}

func New(repo Repository, cache cache.Cache, cfg config.Config) *Handler {
	return &Handler{repo: repo, cache: cache, cfg: cfg}
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
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := utils.ValidateURL(req.URL); err != nil {
		response.BadRequest(c, "Invalid URL")
		return
	}

	var userIDPtr *int
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		const bearerPrefix = "Bearer "
		if len(authHeader) > len(bearerPrefix) && authHeader[:len(bearerPrefix)] == bearerPrefix {
			tokenString := authHeader[len(bearerPrefix):]
			claims, err := auth.ValidateToken(tokenString, h.cfg.JWT.AccessTokenSecret)
			if err == nil {
				uID := claims.UserID
				userIDPtr = &uID
			}
		}
	}

	for i := 0; i < MaxRetries; i++ {
		shortCode, err := utils.GenerateShortCode()
		if err != nil {
			response.InternalServerError(c, "Failed to generate short code")
			return
		}
		id, err := h.repo.CreateShortURL(shortCode, req.URL, userIDPtr)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == UniqueViolation {
				continue
			}
			response.InternalServerError(c, "Failed to create short url")
			return
		}
		if h.cache != nil {
			h.cache.Set(c, cache.URLCacheKey(shortCode), req.URL, time.Hour)
		}
		response.Created(c, gin.H{
			"id":           id,
			"original_url": req.URL,
			"short_code":   shortCode,
			"short_url":    strings.TrimSuffix(h.cfg.Server.BaseURL, "/") + "/" + shortCode,
		})
		return
	}

	response.InternalServerError(c, fmt.Sprintf("Failed to create short url after %d attempts", MaxRetries))

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
	// use cache to get url
	var url string
	var err error
	if h.cache != nil {
		url, err = h.cache.Get(c, cache.URLCacheKey(code))
	} else {
		err = errors.New("cache not configured")
	}

	// if cache miss get from database and set in cache
	if err != nil {
		logger.Log.Info("database lookup")
		url, err = h.repo.GetURLByCode(code)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				response.NotFound(c, "URL not found")
				return
			}
			response.InternalServerError(c, "Database error")
			return
		}
		if h.cache != nil {
			h.cache.Set(c, cache.URLCacheKey(code), url, time.Hour)
		}
	} else {
		logger.Log.Info("cache hit ")
	}
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
		response.BadRequest(c, "Invalid id")
		return
	}
	// Get short code before deleting from database
	code, _ := h.repo.GetCodeByID(id)

	err = h.repo.DeleteURL(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(c, "URL not found")
			return
		}
		response.InternalServerError(c, "Failed to delete url")
		return
	}

	// Invalidate cache
	if code != "" && h.cache != nil {
		_ = h.cache.Delete(c, cache.URLCacheKey(code))
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
//	@Param	start_date	query	string	false	"Start date"
//	@Param	end_date	query	string	false	"End date"
//	@Param	limit	query	int	false	"Limit"
//	@Param	offset	query	int	false	"Offset"
//	@Param	sort_by	query	string	false	"Sort by"
//	@Param	sort_order	query	string	false	"Sort order"
//	@Param	min_clicks	query	int	false	"Minimum clicks"
//	@Param	max_clicks	query	int	false	"Maximum clicks"
//	@Param	search	query	string	false	"Search"
//	@Produce	json
//	@Success	200	{object}	models.SuccessResponse
//	@Failure	404	{object}	models.ErrorResponse
//	@Failure	500	{object}	models.ErrorResponse
//	@Router	/urls/all [get]
func (h *Handler) ListAllHandler(c *gin.Context) {
	var opts models.ListOptions
	if err := c.ShouldBindQuery(&opts); err != nil {
		response.BadRequest(c, "Invalid query parameters")
		return
	}
	urls, err := h.repo.GetAllURLs(opts)
	if err != nil {
		response.InternalServerError(c, "Failed to get all urls")
		return
	}
	response.OK(c, gin.H{
		"urls": urls,
	})
}
