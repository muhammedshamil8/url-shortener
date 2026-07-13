package repository

import (
	"database/sql"
	"log"
	"time"
)

type URL struct {
	ID          int64     `json:"id"`
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
	click_count int       `json:"click_count"`
}

var db *sql.DB

func Init(database *sql.DB) {
	db = database
}

func CreateShortURL(shortCode, url string) (int64, error) {
	var id int64
	err := db.QueryRow("INSERT INTO urls (short_code, original_url) VALUES ($1, $2) RETURNING id", shortCode, url).Scan(&id)
	if err != nil {
		log.Println("Error inserting into database:", err)
		return 0, err
	}
	return id, nil
}

func GetURLByCode(code string) (string, error) {
	var url string
	err := db.QueryRow("SELECT original_url FROM urls WHERE short_code = $1", code).Scan(&url)
	if err != nil {
		log.Println("Error getting url from database:", err)
		return "", err
	}
	return url, nil
}

func DeleteURL(id int) error {
	_, err := db.Exec("DELETE FROM urls WHERE id = $1", id)
	if err != nil {
		log.Println("Error deleting from database:", err)
		return err
	}
	return nil
}

func GetAllURLs() ([]URL, error) {
	var urls []URL
	rows, err := db.Query("SELECT id, short_code, original_url, created_at, click_count FROM urls")
	if err != nil {
		log.Println("Error getting urls from database:", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var url URL
		if err := rows.Scan(&url.ID, &url.ShortCode, &url.OriginalURL, &url.CreatedAt, &url.click_count); err != nil {
			log.Println("Error getting urls from database:", err)
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}
