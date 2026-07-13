package handlers

import (
	"database/sql"
	"errors"
	"fmt"

	// "log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/muhammedshamil8/url-shortener/internal/repository"
	"github.com/muhammedshamil8/url-shortener/internal/models"
	"github.com/muhammedshamil8/url-shortener/internal/utils"
)

const (
	MaxRetries      = 5
	UniqueViolation = "23505"
)

func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to URL Shortener Service",
	})
}

func ShortenHandler(c *gin.Context) {
	var req models.ShortenRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	if err := utils.ValidateURL(req.URL); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid URL",
		})
		return
	}

	for i := 0; i < MaxRetries; i++ {
		shortCode, err := utils.GenerateShortCode()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to generate short code",
			})
			return
		}
		id, err := repository.CreateShortURL(shortCode, req.URL)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == UniqueViolation {
				continue
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to create short url",
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"id":           id,
			"original_url": req.URL,
			"short_code":   shortCode,
			"short_url":    os.Getenv("BASE_URL") + shortCode,
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"message": fmt.Sprintf("Failed to create short url after %d attempts", MaxRetries),
	})

}

func RedirectHandler(c *gin.Context) {
	code := c.Param("code")
	url, err := repository.GetURLByCode(code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "URL not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Database error",
		})
		return
	}
	// if err := repository.IncrementClickCount(code); err != nil {
	// 	log.Println("Error incrementing click count:", err)
	// }
	c.Redirect(http.StatusSeeOther, url)
}

func DeleteHandler(c *gin.Context) {
	idstr := c.Param("id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid id",
		})
		return
	}
	err = repository.DeleteURL(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Failed to delete url",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "URL deleted successfully",
	})
}

func ListAllHandler(c *gin.Context) {
	urls, err := repository.GetAllURLs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get all urls",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"urls": urls,
	})
}
