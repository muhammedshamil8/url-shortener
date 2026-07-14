package repository

import (
	"github.com/muhammedshamil8/url-shortener/internal/logger"
	"github.com/muhammedshamil8/url-shortener/internal/models"
)

func (r *Repository) CreateUser(username, email, passwordHash string) (int64, error) {
	var id int64
	err := r.db.QueryRow(
		`INSERT INTO users (username, email, password_hash) 
		VALUES ($1, $2, $3) RETURNING id`,
		username, email, passwordHash,
	).Scan(&id)

	if err != nil {
		logger.Log.Error("Error inserting into database", "error", err)
		return 0, err
	}
	return id, nil
}

func (r *Repository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(
		`SELECT id, username, email, password_hash, created_at 
		FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
