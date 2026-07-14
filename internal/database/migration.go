package database

import "database/sql"

func MigrateUserTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) NOT NULL UNIQUE,
		password TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func MigrateUrlTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS urls (
		id SERIAL PRIMARY KEY,
		short_code VARCHAR(20) NOT NULL UNIQUE,
		original_url TEXT NOT NULL,
		click_count INT DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE

	);
	`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
