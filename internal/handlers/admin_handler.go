package handlers

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/muhammedshamil8/url-shortener/internal/response"
)

// AdminListUsers lists all registered users in the database.
func (h *Handler) AdminListUsers(c *gin.Context) {
	users, err := h.repo.GetAllUsers()
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve users")
		return
	}
	response.OK(c, gin.H{
		"users": users,
	})
}

// AdminDeleteUser deletes a user from the database.
func (h *Handler) AdminDeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	err = h.repo.DeleteUser(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalServerError(c, "Failed to delete user")
		return
	}

	response.OK(c, gin.H{
		"status": "success",
	})
}
