package handlers

import (
	"database/sql"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/muhammedshamil8/url-shortener/internal/auth"
	"github.com/muhammedshamil8/url-shortener/internal/models"
	"github.com/muhammedshamil8/url-shortener/internal/response"
)

// Register godoc
//
//	@Summary	Register a new account
//	@Description	Register a new account
//	@Tags	Users
//	@Param	request	body	models.RegisterRequest	true	"Request body"
//	@Produce	json
//	@Success	200	{object}	models.SuccessResponse
//	@Failure	400	{object}	models.ErrorResponse
//	@Failure	500	{object}	models.ErrorResponse
//	@Router	/register [post]
func (h *Handler) RegisterHandler(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.BindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		response.InternalServerError(c, "Failed to hash password")
		return
	}

	id, err := h.repo.CreateUser(req.Username, req.Email, hashedPassword)
	if err != nil {
		response.InternalServerError(c, "Failed to create user")
		return
	}
	response.OK(c, gin.H{
		"id":       id,
		"username": req.Username,
		"email":    req.Email,
	})
}

// Login godoc
//
//	@Summary	Login to account
//	@Description	Login to account
//	@Tags	Users
//	@Param	request	body	models.LoginRequest	true	"Request body"
//	@Produce	json
//	@Success	200	{object}	models.SuccessResponse
//	@Failure	400	{object}	models.ErrorResponse
//	@Failure	401	{object}	models.ErrorResponse
//	@Failure	500	{object}	models.ErrorResponse
//	@Router	/login [post]
func (h *Handler) LoginHandler(c *gin.Context) {
	var req models.LoginRequest
	if err := c.BindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	user, err := h.repo.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.Unauthorized(c, "Invalid credentials")
			return
		}
		response.InternalServerError(c, "Database error")
		return
	}

	if err := auth.CheckPasswordHash(req.Password, user.PasswordHash); err != nil {
		response.Unauthorized(c, "Invalid credentials")
		return
	}

	accessToken, err := auth.GenerateToken(user.ID, user.Email, h.cfg.JWT.AccessTokenSecret, h.cfg.JWT.AccessTokenExpiry)
	if err != nil {
		response.InternalServerError(c, "Failed to generate access token")
		return
	}

	refreshToken, err := auth.GenerateToken(user.ID, user.Email, h.cfg.JWT.RefreshTokenSecret, h.cfg.JWT.RefreshTokenExpiry)
	if err != nil {
		response.InternalServerError(c, "Failed to generate refresh token")
		return
	}

	response.OK(c, models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         *user,
	})

}

// Refresh godoc
//
//	@Summary		Refresh access token
//	@Description	Refresh access token using a refresh token
//	@Tags			Users
//	@Param			request	body	models.RefreshRequest	true	"Request body"
//	@Produce		json
//	@Success		200	{object}	models.RefreshResponse
//	@Failure		400	{object}	models.ErrorResponse
//	@Failure		401	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/refresh [post]
func (h *Handler) RefreshHandler(c *gin.Context) {
	var req models.RefreshRequest
	if err := c.BindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	claims, err := auth.ValidateToken(req.RefreshToken, h.cfg.JWT.RefreshTokenSecret)
	if err != nil {
		response.Unauthorized(c, "Invalid or expired refresh token")
		return
	}

	userID := claims.UserID
	email := claims.Email

	newAccessToken, err := auth.GenerateToken(userID, email, h.cfg.JWT.AccessTokenSecret, h.cfg.JWT.AccessTokenExpiry)
	if err != nil {
		response.InternalServerError(c, "Failed to generate access token")
		return
	}

	newRefreshToken, err := auth.GenerateToken(userID, email, h.cfg.JWT.RefreshTokenSecret, h.cfg.JWT.RefreshTokenExpiry)
	if err != nil {
		response.InternalServerError(c, "Failed to generate refresh token")
		return
	}

	response.OK(c, models.RefreshResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	})
}

