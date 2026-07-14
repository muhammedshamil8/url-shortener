package handlers

import "github.com/muhammedshamil8/url-shortener/internal/models"

type URLRepository interface {
	CreateShortURL(shortCode, url string) (int64, error)
	GetURLByCode(code string) (string, error)
	DeleteURL(id int) error
	GetAllURLs(opts models.ListOptions) ([]models.URL, error)

	Health() error
}
