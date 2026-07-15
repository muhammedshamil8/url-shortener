package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/muhammedshamil8/url-shortener/internal/logger"
	"github.com/muhammedshamil8/url-shortener/internal/models"
)

type Repository struct {
	db *sql.DB
}

const healthCheckTimeout = time.Second

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateShortURL(shortCode, url string, userID *int) (int64, error) {
	var id int64
	err := r.db.QueryRow(
		`INSERT INTO urls (short_code, original_url, user_id) 
		VALUES ($1, $2, $3) RETURNING id`,
		shortCode, url, userID,
	).Scan(&id)

	if err != nil {
		logger.Log.Error("Error inserting into database", "error", err)
		return 0, err
	}
	return id, nil
}

func (r *Repository) GetURLByCode(code string) (string, error) {
	var url string

	err := r.db.QueryRow(
		`UPDATE urls
		 SET click_count = click_count + 1
		 WHERE short_code = $1
		 RETURNING original_url`,
		code,
	).Scan(&url)

	if err != nil {
		return "", err
	}

	return url, nil
}

func (r *Repository) GetCodeByID(id int) (string, error) {
	var code string
	err := r.db.QueryRow("SELECT short_code FROM urls WHERE id = $1", id).Scan(&code)
	if err != nil {
		return "", err
	}
	return code, nil
}

func (r *Repository) DeleteURL(id int) error {
	_, err := r.db.Exec("DELETE FROM urls WHERE id = $1", id)
	if err != nil {
		logger.Log.Error("Error deleting from database", "error", err)
		return err
	}
	return nil
}

func (r *Repository) GetAllURLs(opts models.ListOptions) ([]models.URL, error) {
	var urls []models.URL

	opts.Normalize()
	query := `
		SELECT id, short_code, original_url, created_at, click_count
		FROM urls
	`

	args := []any{}
	conditions := []string{}

	sortColumn := "created_at"
	switch opts.Sort {
	case "created_at":
		sortColumn = "created_at"
	case "click_count":
		sortColumn = "click_count"
	case "short_code":
		sortColumn = "short_code"
	default:
		sortColumn = "created_at"
	}

	orderDirection := "DESC"
	switch opts.Order {
	case "ASC":
		orderDirection = "ASC"
	case "DESC":
		orderDirection = "DESC"
	default:
		orderDirection = "DESC"
	}

	if opts.Search != "" {
		param := len(args) + 1
		conditions = append(
			conditions,
			fmt.Sprintf(
				"(original_url ILIKE $%d OR short_code ILIKE $%d)",
				param,
				param,
			),
		)

		args = append(args, "%"+opts.Search+"%")
	}

	if opts.MinClicks > 0 {
		param := len(args) + 1
		conditions = append(conditions, fmt.Sprintf("click_count >= $%d", param))
		args = append(args, opts.MinClicks)
	}

	if opts.MaxClicks > 0 {
		param := len(args) + 1
		conditions = append(conditions, fmt.Sprintf("click_count <= $%d", param))
		args = append(args, opts.MaxClicks)
	}

	if !opts.MinDate.IsZero() {
		param := len(args) + 1
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", param))
		args = append(args, opts.MinDate)
	}

	if !opts.MaxDate.IsZero() {
		param := len(args) + 1
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", param))
		args = append(args, opts.MaxDate)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	limitPos := len(args) + 1
	offsetPos := len(args) + 2
	offset := (opts.Page - 1) * opts.Limit

	query += fmt.Sprintf(
		" ORDER BY %s %s LIMIT $%d OFFSET $%d",
		sortColumn,
		orderDirection,
		limitPos,
		offsetPos,
	)

	args = append(args, opts.Limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		logger.Log.Error("Error getting urls from database", "error", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var url models.URL
		if err := rows.Scan(&url.ID, &url.ShortCode, &url.OriginalURL, &url.CreatedAt, &url.ClickCount); err != nil {
			logger.Log.Error("Error getting urls from database", "error", err)
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func (r *Repository) Health() error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		healthCheckTimeout,
	)
	defer cancel()
	err := r.db.PingContext(ctx)
	if err != nil {
		logger.Log.Error("Error pinging database", "error", err)
		return err
	}
	return nil
}

// func IncrementClickCount(code string) error {
// 	_, err := db.Exec("UPDATE urls SET click_count = click_count + 1 WHERE short_code = $1", code)
// 	if err != nil {
// 		log.Println("Error incrementing click count:", err)
// 		return err
// 	}
// 	return nil
// }
