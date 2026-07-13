package models

type ShortenRequest struct {
	URL string `json:"url" binding:"required"`
}