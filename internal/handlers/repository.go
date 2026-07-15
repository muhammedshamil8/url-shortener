package handlers

import "github.com/muhammedshamil8/url-shortener/internal/models"

type Repository interface {
	URLRepository
	UserRepository
}

type URLRepository interface {
	CreateShortURL(shortCode, url string, userID *int) (int64, error)
	GetURLByCode(code string) (string, error)
	DeleteURL(id int) error
	GetAllURLs(opts models.ListOptions) ([]models.URL, error)

	Health() error
}

type UserRepository interface {
	CreateUser(username, email, passwordHash string) (int64, error)
	GetUserByEmail(email string) (*models.User, error)
	GetAllURLsByUserEmail(email string) ([]models.URL, error)
	DeleteUserURL(id int) error
	GetAllUsers() ([]models.User, error)
	DeleteUser(id int) error
}
