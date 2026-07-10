package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func main() {
	InitDB()
	
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the URL Shortener API"})
	})
	
	// r.POST("/shorten", shortenHandler)
	// r.GET("/:code", redirectHandler)
	// r.GET("/urls", listAllHandler)
	// r.DELETE("/:code", deleteHandler)
	
	r.Run(":8080")
}