package database

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/muhammedshamil8/url-shortener/internal/config"
	"github.com/muhammedshamil8/url-shortener/internal/logger"
)

var DB *sql.DB

func InitDB(cfg config.DBConfig) (*sql.DB, error) {
	var err error
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.SSLMode)
	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Verify the connection is alive
	if err := DB.Ping(); err != nil {
		return nil, err
	}

	logger.Log.Info("Successfully connected to PostgreSQL")

	return DB, nil
}
