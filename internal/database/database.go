package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func InitDB() (*sql.DB,error) {
	var err error
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"))
	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		return nil,err
	}

	// Verify the connection is alive
	if err := DB.Ping(); err != nil {
		return nil,err
	}

	log.Println("Successfully connected to PostgreSQL")

	return DB, nil
}




