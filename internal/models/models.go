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
}

type SuccessResponse struct {
	Status    string `json:"status" example:"success"`
	Data      any    `json:"data"`
	RequestID string `json:"request_id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type ErrorResponse struct {
	Status    string `json:"status" example:"error"`
	Message   string `json:"message" example:"Invalid URL"`
	RequestID string `json:"request_id" example:"550e8400-e29b-41d4-a716-446655440000"`
}

type ListOptions struct {
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
	Sort   string `form:"sort"`
	Order  string `form:"order"`
	Search string `form:"search"`

	MinClicks int `form:"min_clicks"`
	MaxClicks int `form:"max_clicks"`
}

func (o *ListOptions) Normalize() {
	if o.Page < 1 {
		o.Page = 1
	}
	if o.Limit < 1 || o.Limit > 100 {
		o.Limit = 20
	}
}
