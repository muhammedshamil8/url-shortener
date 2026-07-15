package repository

import (
	"database/sql"

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

func (r *Repository) DeleteUserURL(id int) error {
	result, err := r.db.Exec(
		`DELETE FROM urls WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *Repository) UpdateUserURL(id int, email string, newURL string) error {
	result, err := r.db.Exec(
		`UPDATE urls u
		 SET original_url = $1
		 FROM users usr
		 WHERE u.user_id = usr.id AND u.id = $2 AND usr.email = $3`,
		newURL, id, email,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *Repository) GetAllURLsByUserEmail(email string) ([]models.URL, error) {
	rows, err := r.db.Query(
		`SELECT u.id, u.short_code, u.original_url, u.click_count, u.created_at, u.user_id 
		FROM urls u 
		JOIN users usr ON u.user_id = usr.id 
		WHERE usr.email = $1`,
		email,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []models.URL
	for rows.Next() {
		var url models.URL
		err := rows.Scan(&url.ID, &url.ShortCode, &url.OriginalURL, &url.ClickCount, &url.CreatedAt, &url.UserID)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

func (r *Repository) GetAllUsers() ([]models.User, error) {
	rows, err := r.db.Query(
		`SELECT id, username, email, created_at FROM users ORDER BY id ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Repository) DeleteUser(id int) error {
	result, err := r.db.Exec(
		`DELETE FROM users WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
