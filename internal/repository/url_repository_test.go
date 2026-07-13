package repository

import (
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/muhammedshamil8/url-shortener/internal/database"
	"github.com/muhammedshamil8/url-shortener/internal/models"
)

func setupTestDB(t *testing.T) *Repository {
	t.Helper()

	// Load .env if present, but ignore error if missing (e.g. in CI)
	_ = godotenv.Load("../../.env")
	db, err := database.InitDB()
	if err != nil {
		t.Fatal(err)
	}
	if err := database.MigrateUrlTable(db); err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("TRUNCATE TABLE urls RESTART IDENTITY")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		db.Close()
	})
	return New(db)
}

func testURL() models.URL {
    return models.URL{
        OriginalURL: "https://google.com",
        ShortCode:   "abc",
		CreatedAt:     time.Now(),
		ClickCount:    0,
    }
}

func TestCreateShortURL(t *testing.T) {
	repo := setupTestDB(t)
	URL := testURL()
	id, err := repo.CreateShortURL(URL.ShortCode, URL.OriginalURL); 
	if err != nil {
		t.Fatalf("failed to create short URL: %v", err)
	}
	if id == 0 {
		t.Fatalf("expected id > 0, got %d", id)
	}
	var url models.URL
	err = repo.db.QueryRow(
		"SELECT short_code, original_url, click_count FROM urls WHERE id = $1",
		id).Scan(&url.ShortCode, &url.OriginalURL, &url.ClickCount)
	if err != nil {
		t.Fatalf("failed to get URL by code: %v", err)
	}
	if url.ShortCode != "abc" {
		t.Fatalf("expected short_code %q, got %q", "abc", url.ShortCode)
	}
	if url.OriginalURL != "https://google.com" {
		t.Fatalf("expected original_url %q, got %q", "https://google.com", url.OriginalURL)
	}
	if url.ClickCount != 0 {
		t.Fatalf("expected click_count %d, got %d", 0, url.ClickCount)
	}


}

func TestGetURLByCode(t *testing.T) {
	repo := setupTestDB(t)
	url := testURL()
	if id, err := repo.CreateShortURL(url.ShortCode, url.OriginalURL); err != nil {
		t.Fatalf("failed to create short URL: %v", err)
	} else {
		t.Log("Created short URL with ID:", id)
	}
	url2, err := repo.GetURLByCode(url.ShortCode)
	if err != nil {
		t.Fatalf("failed to get URL by code: %v", err)
	}
	if url2 != url.OriginalURL {
		t.Fatalf("got %q, want %q", url2, url.OriginalURL)
	}
}

func TestDeleteURL(t *testing.T) {
	repo := setupTestDB(t)
	url := testURL()
	if id, err := repo.CreateShortURL(url.ShortCode, url.OriginalURL); err != nil {
		t.Fatalf("failed to create short URL: %v", err)
	} else {
		t.Log("Created short URL with ID:", id)
		if err := repo.DeleteURL(int(id)); err != nil {
			t.Fatalf("failed to delete URL: %v", err)
		} else {
			var count int

			err = repo.db.QueryRow(
			"SELECT COUNT(*) FROM urls WHERE id = $1",
			id,
		).Scan(&count)

			if err != nil {
				t.Fatal(err)
			}

			if count != 0 {
				t.Fatalf("expected 0 rows, got %d", count)
			}
			t.Log("Deleted short URL with ID:", id)
		}
	}
}

func TestGetAllURLs(t *testing.T) {
	repo := setupTestDB(t)
	repo.CreateShortURL("abc","https://google.com")
	repo.CreateShortURL("xyz","https://github.com")
	repo.CreateShortURL("def","https://golang.org")

	urls, err := repo.GetAllURLs()
	if err != nil {
		t.Fatalf("failed to get all URLs: %v", err)
	}
	if len(urls) != 3 {
    t.Fatalf("expected 3 URLs, got %d", len(urls))
}
}
