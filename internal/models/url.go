package models

import "time"

type ShortenRequest struct {
	URL string `json:"url" binding:"required" example:"https://google.com"`
}

type URL struct {
	ID          int       `json:"id" example:"1"`
	OriginalURL string    `json:"original_url" example:"https://google.com"`
	ShortCode   string    `json:"short_code" example:"abc123"`
	CreatedAt   time.Time `json:"created_at"`
	ClickCount  int       `json:"click_count" example:"5"`
	UserID      int       `json:"user_id" example:"1"`
}
