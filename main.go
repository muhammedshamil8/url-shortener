package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	."github.com/muhammedshamil8/url-shortener/internal/database"
	."github.com/muhammedshamil8/url-shortener/internal/handlers"
	."github.com/muhammedshamil8/url-shortener/internal/repository"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env")
	}
	DB , err := InitDB()
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer DB.Close()
	Init(DB)
	if err := MigrateUrlTable(); err != nil {
		log.Fatal("Error migrating tables:", err)
	}

	r := gin.Default()
	r.GET("/health", HealthCheckHandler)

	r.POST("/shorten", ShortenHandler)
	r.GET("/:code", RedirectHandler)	
	r.GET("/urls", ListAllHandler)
	r.DELETE("/:id", DeleteHandler)

	port := os.Getenv("APP_PORT")
	if port == "" {
		log.Println("APP_PORT is not set")
		log.Println("Automatic setting to 8080")
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server: " + err.Error())
	}

}
