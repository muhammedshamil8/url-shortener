package handlers

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/muhammedshamil8/url-shortener/internal/response"
	"github.com/muhammedshamil8/url-shortener/internal/utils"
)

// GetProfile godoc
//
//	@Summary	Get profile
//	@Description	Get profile
//	@Tags	Users
//	@Produce	json
//	@Success	200	{object}	models.SuccessResponse
//	@Failure	404	{object}	models.ErrorResponse
//	@Failure	500	{object}	models.ErrorResponse
//	@Router	/me [get]
func (h *Handler) GetProfileHandler(c *gin.Context) {
	email := c.GetString("email")
	user, err := h.repo.GetUserByEmail(email)
	if err != nil {
		response.InternalServerError(c, "Failed to get user")
		return
	}
	response.OK(c, gin.H{
		"user": user,
	})
}

// ListUserURLs godoc
//
//	@Summary	List user URLs
//	@Description	List user URLs
//	@Tags	Users
//	@Produce	json
//	@Success	200	{object}	models.SuccessResponse
//	@Failure	404	{object}	models.ErrorResponse
//	@Failure	500	{object}	models.ErrorResponse
//	@Router	/me/urls [get]
func (h *Handler) ListUserURLs(c *gin.Context) {
	email := c.GetString("email")
	urls, err := h.repo.GetAllURLsByUserEmail(email)
	if err != nil {
		response.InternalServerError(c, "Failed to get user urls")
		return
	}
	response.OK(c, gin.H{
		"urls": urls,
	})
}

// DeleteURL godoc
//
//	@Summary	Delete URL
//	@Description	Delete URL
//	@Tags	Users
//	@Produce	json
//	@Param	id	path	int	true	"URL ID"
//	@Success	200	{object}	models.SuccessResponse
//	@Failure	400	{object}	models.ErrorResponse
//	@Failure	404	{object}	models.ErrorResponse
//	@Failure	500	{object}	models.ErrorResponse
//	@Router	/me/urls/{id} [delete]
func (h *Handler) DeleteURL(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		response.BadRequest(c, "Invalid id")
		return
	}
	// Get short code before deleting from database
	code, _ := h.repo.GetCodeByID(id)

	err = h.repo.DeleteUserURL(id)
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
		_ = h.cache.Delete(c, code)
	}

	response.OK(c, gin.H{
		"message": "URL deleted successfully",
	})
}

// UpdateURL godoc
//
//	@Summary	Update URL
//	@Description	Update URL target original URL
//	@Tags	Users
//	@Accept	json
//	@Produce	json
//	@Param	id	path	int	true	"URL ID"
//	@Param	request	body	models.ShortenRequest	true	"Request body"
//	@Success	200	{object}	models.SuccessResponse
//	@Failure	400	{object}	models.ErrorResponse
//	@Failure	404	{object}	models.ErrorResponse
//	@Failure	500	{object}	models.ErrorResponse
//	@Router	/me/urls/{id} [put]
func (h *Handler) UpdateURL(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		response.BadRequest(c, "Invalid id")
		return
	}

	var req struct {
		URL string `json:"url" binding:"required"`
	}
	if err := c.BindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := utils.ValidateURL(req.URL); err != nil {
		response.BadRequest(c, "Invalid URL")
		return
	}

	// Get short code before updating to invalidate cache
	code, _ := h.repo.GetCodeByID(id)

	err = h.repo.UpdateUserURL(id, c.GetString("email"), req.URL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(c, "URL not found or unauthorized")
			return
		}
		response.InternalServerError(c, "Failed to update url")
		return
	}

	// Invalidate cache
	if code != "" && h.cache != nil {
		_ = h.cache.Delete(c, code)
	}

	response.OK(c, gin.H{
		"message": "URL updated successfully",
	})
}
