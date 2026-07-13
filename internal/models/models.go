package models

import "time"

type ShortenRequest struct {
	URL string `json:"url" binding:"required"`
}

type URL struct {
	ID            int       `json:"id" db:"id"`
	OriginalURL   string    `json:"original_url" db:"original_url"`
	ShortCode     string    `json:"short_code" db:"short_code"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	ClickCount    int       `json:"click_count" db:"click_count"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}


type HealthResponse struct {
	Status string `json:"status"`
}
