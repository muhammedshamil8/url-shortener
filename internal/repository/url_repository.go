package repository

import (
	"database/sql"

	"github.com/muhammedshamil8/url-shortener/internal/logger"
	"github.com/muhammedshamil8/url-shortener/internal/models"
)

type Repository struct {
    db *sql.DB
}

func New(db *sql.DB) *Repository {
    return &Repository{db: db}
}


func (r *Repository)CreateShortURL(shortCode, url string) (int64, error) {
	var id int64
	err := r.db.QueryRow(
		`INSERT INTO urls (short_code, original_url) 
		VALUES ($1, $2) RETURNING id`,
		shortCode, url,
	).Scan(&id)

	if err != nil {
		logger.Log.Error("Error inserting into database", "error", err)
		return 0, err
	}
	return id, nil
}

func (r *Repository)GetURLByCode(code string) (string, error) {
	var url string

	err := r.db.QueryRow(
		`UPDATE urls
		 SET click_count = click_count + 1
		 WHERE short_code = $1
		 RETURNING original_url`,
		code,
	).Scan(&url)

	if err != nil {
		return "", err
	}

	return url, nil
}

func (r *Repository)DeleteURL(id int) error {
	_, err := r.db.Exec("DELETE FROM urls WHERE id = $1", id)
	if err != nil {
		logger.Log.Error("Error deleting from database", "error", err)
		return err
	}
	return nil
}

func (r *Repository)GetAllURLs() ([]models.URL, error) {
	var urls []models.URL
	rows, err := r.db.Query("SELECT id, short_code, original_url, created_at, click_count FROM urls")
	if err != nil {
		logger.Log.Error("Error getting urls from database", "error", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var url models.URL
		if err := rows.Scan(&url.ID, &url.ShortCode, &url.OriginalURL, &url.CreatedAt, &url.ClickCount); err != nil {
			logger.Log.Error("Error getting urls from database", "error", err)
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}

// func IncrementClickCount(code string) error {
// 	_, err := db.Exec("UPDATE urls SET click_count = click_count + 1 WHERE short_code = $1", code)
// 	if err != nil {
// 		log.Println("Error incrementing click count:", err)
// 		return err
// 	}
// 	return nil
// }
