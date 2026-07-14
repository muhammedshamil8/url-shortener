package repository

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/muhammedshamil8/url-shortener/internal/config"
	"github.com/muhammedshamil8/url-shortener/internal/database"
	"github.com/muhammedshamil8/url-shortener/internal/models"
)

func setupTestUserDB(t *testing.T) *Repository {
	t.Helper()

	err := godotenv.Load("../../.env")
	if err != nil {
		t.Log("Failed to load .env")
	}
	cfg := config.Load()
	db, err := database.InitDB(cfg.DB)
	if err != nil {
		t.Fatal(err)
	}
	_, _ = db.Exec("DROP TABLE IF EXISTS urls CASCADE")
	_, _ = db.Exec("DROP TABLE IF EXISTS users CASCADE")
	if err := database.MigrateUserTable(db); err != nil {
		t.Fatal(err)
	}
	if err := database.MigrateUrlTable(db); err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("TRUNCATE TABLE urls, users RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		db.Close()
	})
	return New(db)
}

func TestCreateUser(t *testing.T) {
	repo := setupTestUserDB(t)

	username := "testuser"
	email := "test@example.com"
	passwordHash := "hashed_password"

	id, err := repo.CreateUser(username, email, passwordHash)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	if id == 0 {
		t.Fatalf("expected id > 0, got %d", id)
	}

	var user models.User
	err = repo.db.QueryRow(
		"SELECT username, email, password_hash FROM users WHERE id = $1",
		id,
	).Scan(&user.Username, &user.Email, &user.PasswordHash)
	if err != nil {
		t.Fatalf("failed to fetch user from DB: %v", err)
	}

	if user.Username != username || user.Email != email || user.PasswordHash != passwordHash {
		t.Fatalf("saved user doesn't match")
	}
}

func TestGetUserByEmail(t *testing.T) {
	repo := setupTestUserDB(t)

	username := "testuser"
	email := "test@example.com"
	passwordHash := "hashed_password"

	_, err := repo.CreateUser(username, email, passwordHash)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	user, err := repo.GetUserByEmail(email)
	if err != nil {
		t.Fatalf("failed to get user by email: %v", err)
	}

	if user.Username != username || user.Email != email || user.PasswordHash != passwordHash {
		t.Fatalf("fetched user doesn't match")
	}
}

func TestGetAllURLsByUserEmail(t *testing.T) {
	repo := setupTestUserDB(t)

	username := "testuser"
	email := "test@example.com"
	passwordHash := "hashed_password"

	userID, err := repo.CreateUser(username, email, passwordHash)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Insert urls with user_id
	_, err = repo.db.Exec(
		`INSERT INTO urls (short_code, original_url, user_id) 
		VALUES ($1, $2, $3), ($4, $5, $6)`,
		"abc", "https://google.com", userID,
		"def", "https://yahoo.com", userID,
	)
	if err != nil {
		t.Fatalf("failed to insert test urls: %v", err)
	}

	urls, err := repo.GetAllURLsByUserEmail(email)
	if err != nil {
		t.Fatalf("failed to get all URLs by user email: %v", err)
	}

	if len(urls) != 2 {
		t.Fatalf("expected 2 urls, got %d", len(urls))
	}
	if urls[0].ShortCode != "abc" || urls[1].ShortCode != "def" {
		t.Fatalf("retrieved urls don't match inserted urls")
	}
}

func TestDeleteUserURL(t *testing.T) {
	repo := setupTestUserDB(t)

	username := "testuser"
	email := "test@example.com"
	passwordHash := "hashed_password"

	userID, err := repo.CreateUser(username, email, passwordHash)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	var urlID int
	err = repo.db.QueryRow(
		`INSERT INTO urls (short_code, original_url, user_id) 
		VALUES ($1, $2, $3) RETURNING id`,
		"abc", "https://google.com", userID,
	).Scan(&urlID)
	if err != nil {
		t.Fatalf("failed to insert test url: %v", err)
	}

	err = repo.DeleteUserURL(urlID)
	if err != nil {
		t.Fatalf("failed to delete user url: %v", err)
	}

	var count int
	err = repo.db.QueryRow("SELECT COUNT(*) FROM urls WHERE id = $1", urlID).Scan(&count)
	if err != nil {
		t.Fatalf("failed to query urls count: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected url to be deleted, but it exists")
	}
}

func TestGetAllUsers(t *testing.T) {
	repo := setupTestUserDB(t)

	_, err := repo.CreateUser("user1", "user1@example.com", "hash1")
	if err != nil {
		t.Fatalf("failed to create user 1: %v", err)
	}

	_, err = repo.CreateUser("user2", "user2@example.com", "hash2")
	if err != nil {
		t.Fatalf("failed to create user 2: %v", err)
	}

	users, err := repo.GetAllUsers()
	if err != nil {
		t.Fatalf("failed to get all users: %v", err)
	}

	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
	if users[0].Username != "user1" || users[1].Username != "user2" {
		t.Fatalf("unexpected user details retrieved")
	}
}

func TestDeleteUser(t *testing.T) {
	repo := setupTestUserDB(t)

	id, err := repo.CreateUser("user1", "user1@example.com", "hash1")
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	err = repo.DeleteUser(int(id))
	if err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}

	var count int
	err = repo.db.QueryRow("SELECT COUNT(*) FROM users WHERE id = $1", id).Scan(&count)
	if err != nil {
		t.Fatalf("failed to query users count: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected user to be deleted, but it exists")
	}
}
