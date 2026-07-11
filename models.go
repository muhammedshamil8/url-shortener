package main

type ShortenRequest struct {
	URL string `json:"url" binding:"required"`
}