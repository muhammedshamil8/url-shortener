package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/muhammedshamil8/url-shortener/internal/database"
	"github.com/muhammedshamil8/url-shortener/internal/handlers"
	"github.com/muhammedshamil8/url-shortener/internal/repository"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env")
	}
	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.New(db)

	h := handlers.New(repo)

	if err := database.MigrateUrlTable(db); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/health/api", h.HealthCheckHandler)

	r.POST("/shorten", h.ShortenHandler)
	r.GET("/urls/all", h.ListAllHandler)
	r.DELETE("/:id", h.DeleteHandler)
	r.GET("/:code", h.RedirectHandler)

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
