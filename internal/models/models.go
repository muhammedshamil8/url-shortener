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

type User struct {
	ID        int       `json:"id" example:"1"`
	Username  string    `json:"username" example:"shamil"`
	Password  string    `json:"password" example:"password"`
	Email     string    `json:"email" example:"shamil@[EMAIL_ADDRESS]"`
	CreatedAt time.Time `json:"created_at"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User        User   `json:"user"`
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

	MinClicks int       `form:"min_clicks"`
	MaxClicks int       `form:"max_clicks"`
	MinDate   time.Time `form:"min_date"`
	MaxDate   time.Time `form:"max_date"`
}

func (o *ListOptions) Normalize() {
	if o.Page < 1 {
		o.Page = 1
	}
	if o.Limit < 1 || o.Limit > 100 {
		o.Limit = 20
	}
}
