package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"github.com/muhammedshamil8/url-shortener/repository"
	"strconv"
)

func healthCheckHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome to URL Shortener Service",
	})
}


func shortenHandler(c *gin.Context) {
	// recive url from request body
	var req ShortenRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	url := req.URL

	// generate short code
	shortCode, err := generateShortCode()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to generate short code",
		})
		return
	}

	err = repository.CreateShortURL(shortCode, url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create short url",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"short_url": os.Getenv("BASE_URL") + shortCode,
	})

}

func redirectHandler(c *gin.Context) {
	code := c.Param("code")
	url, err := repository.GetURLByCode(code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "URL not found",
		})
		return
	}
	c.Redirect(http.StatusSeeOther, url)
}

func deleteHandler(c *gin.Context) {
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

func listAllHandler(c *gin.Context) {
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