package repository

import (
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/muhammedshamil8/url-shortener/internal/config"
	"github.com/muhammedshamil8/url-shortener/internal/database"
	"github.com/muhammedshamil8/url-shortener/internal/models"
)

func setupTestDB(t *testing.T) *Repository {
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
		CreatedAt:   time.Now(),
		ClickCount:  0,
	}
}

func TestCreateShortURL(t *testing.T) {
	repo := setupTestDB(t)
	URL := testURL()
	id, err := repo.CreateShortURL(URL.ShortCode, URL.OriginalURL)
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

	resetDB := func() {
		_, err := repo.db.Exec("TRUNCATE TABLE urls RESTART IDENTITY")
		if err != nil {
			t.Fatal(err)
		}
	}

	populateDB := func() []models.URL {
		resetDB()
		data := []struct {
			code string
			url  string
		}{
			{"abc", "https://google.com"},
			{"xyz", "https://github.com"},
			{"def", "https://golang.org"},
			{"mno", "https://apple.com"},
			{"pqr", "https://microsoft.com"},
		}

		created := make([]models.URL, 0, len(data))
		baseTime := time.Now().Add(-10 * time.Minute)
		for idx, item := range data {
			id, err := repo.CreateShortURL(item.code, item.url)
			if err != nil {
				t.Fatalf("failed to seed url: %v", err)
			}
			customTime := baseTime.Add(time.Duration(idx) * time.Minute)
			_, err = repo.db.Exec("UPDATE urls SET created_at = $1 WHERE id = $2", customTime, id)
			if err != nil {
				t.Fatalf("failed to update created_at: %v", err)
			}
			var url models.URL
			err = repo.db.QueryRow("SELECT id, short_code, original_url, created_at, click_count FROM urls WHERE id = $1", id).
				Scan(&url.ID, &url.ShortCode, &url.OriginalURL, &url.CreatedAt, &url.ClickCount)
			if err != nil {
				t.Fatalf("failed to query seeded url: %v", err)
			}
			created = append(created, url)
		}
		return created
	}

	t.Run("Empty Table", func(t *testing.T) {
		resetDB()
		urls, err := repo.GetAllURLs(models.ListOptions{
			Page:  1,
			Limit: 10,
		})
		if err != nil {
			t.Fatalf("unexpected error on empty table: %v", err)
		}
		if len(urls) != 0 {
			t.Fatalf("expected 0 URLs, got %d", len(urls))
		}
	})

	t.Run("Pagination and Limits", func(t *testing.T) {
		created := populateDB()

		urls, err := repo.GetAllURLs(models.ListOptions{
			Page:  1,
			Limit: 2,
		})
		if err != nil {
			t.Fatalf("failed to get page 1: %v", err)
		}
		if len(urls) != 2 {
			t.Fatalf("expected 2 URLs, got %d", len(urls))
		}
		if urls[0].ID != created[4].ID || urls[1].ID != created[3].ID {
			t.Fatalf("unexpected pagination results on page 1")
		}

		urls, err = repo.GetAllURLs(models.ListOptions{
			Page:  2,
			Limit: 2,
		})
		if err != nil {
			t.Fatalf("failed to get page 2: %v", err)
		}
		if len(urls) != 2 {
			t.Fatalf("expected 2 URLs, got %d", len(urls))
		}
		if urls[0].ID != created[2].ID || urls[1].ID != created[1].ID {
			t.Fatalf("unexpected pagination results on page 2")
		}

		urls, err = repo.GetAllURLs(models.ListOptions{
			Page:  3,
			Limit: 2,
		})
		if err != nil {
			t.Fatalf("failed to get page 3: %v", err)
		}
		if len(urls) != 1 {
			t.Fatalf("expected 1 URL on page 3, got %d", len(urls))
		}
		if urls[0].ID != created[0].ID {
			t.Fatalf("unexpected pagination results on page 3")
		}

		urls, err = repo.GetAllURLs(models.ListOptions{
			Page:  4,
			Limit: 2,
		})
		if err != nil {
			t.Fatalf("failed to get page 4: %v", err)
		}
		if len(urls) != 0 {
			t.Fatalf("expected 0 URLs on page 4, got %d", len(urls))
		}
	})

	t.Run("Limit Validation", func(t *testing.T) {
		populateDB()
		// Limit 0 should default to 20
		urls, err := repo.GetAllURLs(models.ListOptions{
			Page:  1,
			Limit: 0,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(urls) != 5 {
			t.Fatalf("expected 5 URLs, got %d", len(urls))
		}

		// Limit 500 should default to 20
		urls, err = repo.GetAllURLs(models.ListOptions{
			Page:  1,
			Limit: 500,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(urls) != 5 {
			t.Fatalf("expected 5 URLs, got %d", len(urls))
		}
	})

	t.Run("Page is 0", func(t *testing.T) {
		created := populateDB()
		urls, err := repo.GetAllURLs(models.ListOptions{
			Page:  0,
			Limit: 2,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(urls) != 2 {
			t.Fatalf("expected 2 URLs, got %d", len(urls))
		}
		if urls[0].ID != created[4].ID || urls[1].ID != created[3].ID {
			t.Fatalf("expected page 0 to default to page 1")
		}
	})

	t.Run("Sorting", func(t *testing.T) {
		created := populateDB()

		urls, err := repo.GetAllURLs(models.ListOptions{
			Page:  1,
			Limit: 5,
			Sort:  "created_at",
			Order: "ASC",
		})
		if err != nil {
			t.Fatalf("failed sorting created_at ASC: %v", err)
		}
		if urls[0].ID != created[0].ID || urls[4].ID != created[4].ID {
			t.Fatalf("sorting by created_at ASC failed")
		}

		urls, err = repo.GetAllURLs(models.ListOptions{
			Page:  1,
			Limit: 5,
			Sort:  "short_code",
			Order: "ASC",
		})
		if err != nil {
			t.Fatalf("failed sorting short_code ASC: %v", err)
		}
		if urls[0].ShortCode != "abc" || urls[1].ShortCode != "def" || urls[4].ShortCode != "xyz" {
			t.Fatalf("sorting by short_code ASC failed: first=%s, second=%s, last=%s", urls[0].ShortCode, urls[1].ShortCode, urls[4].ShortCode)
		}

		_, err = repo.GetURLByCode("def")
		if err != nil {
			t.Fatalf("failed to increment click count: %v", err)
		}

		urls, err = repo.GetAllURLs(models.ListOptions{
			Page:  1,
			Limit: 5,
			Sort:  "click_count",
			Order: "DESC",
		})
		if err != nil {
			t.Fatalf("failed sorting click_count DESC: %v", err)
		}
		if urls[0].ShortCode != "def" {
			t.Fatalf("expected def to be first in click_count DESC, got %s", urls[0].ShortCode)
		}
	})

	t.Run("Invalid Sort/Order Defaults", func(t *testing.T) {
		populateDB()

		urls, err := repo.GetAllURLs(models.ListOptions{
			Page:  1,
			Limit: 5,
			Sort:  "invalid_column_name",
			Order: "INVALID_DIR",
		})
		if err != nil {
			t.Fatalf("failed invalid sort/order validation: %v", err)
		}
		if len(urls) != 5 {
			t.Fatalf("expected 5 URLs, got %d", len(urls))
		}
	})

	t.Run("SQL Injection Protection", func(t *testing.T) {
		populateDB()

		urls, err := repo.GetAllURLs(models.ListOptions{
			Page:  1,
			Limit: 5,
			Sort:  "id; DROP TABLE urls; --",
			Order: "ASC",
		})
		if err != nil {
			t.Fatalf("unexpected error on SQL injection test for Sort: %v", err)
		}
		if len(urls) != 5 {
			t.Fatalf("expected 5 URLs, got %d", len(urls))
		}

		urls, err = repo.GetAllURLs(models.ListOptions{
			Page:  1,
			Limit: 5,
			Sort:  "created_at",
			Order: "ASC; DROP TABLE urls; --",
		})
		if err != nil {
			t.Fatalf("unexpected error on SQL injection test for Order: %v", err)
		}
		if len(urls) != 5 {
			t.Fatalf("expected 5 URLs, got %d", len(urls))
		}

		var count int
		err = repo.db.QueryRow("SELECT COUNT(*) FROM urls").Scan(&count)
		if err != nil {
			t.Fatalf("failed to query urls count (table might have been dropped): %v", err)
		}
		if count != 5 {
			t.Fatalf("expected 5 URLs to still exist, got %d", count)
		}
	})

	t.Run("Search: com, limit: 2, page: 1", func(t *testing.T) {
		_ = populateDB()
		urls, err := repo.GetAllURLs(models.ListOptions{
			Page:   1,
			Limit:  2,
			Search: "com",
		})
		if err != nil {
			t.Fatalf("failed search: %v", err)
		}
		if len(urls) != 2 {
			t.Fatalf("expected 2 URLs, got %d", len(urls))
		}
	})

	tests := []struct {
		name     string
		search   string
		expected int
	}{
		{"Search URL", "google", 1},
		{"Search URL (case insensitive)", "GOOGLE", 1},
		{"Search Code", "abc", 1},
		{"Search Missing", "nothing", 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_ = populateDB()
			urls, err := repo.GetAllURLs(models.ListOptions{
				Page:   1,
				Limit:  10,
				Search: test.search,
			})
			if err != nil {
				t.Fatalf("failed search: %v", err)
			}
			if len(urls) != test.expected {
				t.Fatalf("expected %d URLs, got %d", test.expected, len(urls))
			}
		})
	}

	t.Run("Click Count Filtering", func(t *testing.T) {
		_ = populateDB() // Seeding 5 URLs

		// Set clicks:
		// "abc" -> 5 clicks
		for i := 0; i < 5; i++ {
			if _, err := repo.GetURLByCode("abc"); err != nil {
				t.Fatalf("failed to increment click: %v", err)
			}
		}
		// "xyz" -> 3 clicks
		for i := 0; i < 3; i++ {
			if _, err := repo.GetURLByCode("xyz"); err != nil {
				t.Fatalf("failed to increment click: %v", err)
			}
		}
		// "def" -> 1 click
		if _, err := repo.GetURLByCode("def"); err != nil {
			t.Fatalf("failed to increment click: %v", err)
		}

		// MinClicks: 2 -> should return abc (5) and xyz (3)
		urls, err := repo.GetAllURLs(models.ListOptions{
			Page:      1,
			Limit:     10,
			MinClicks: 2,
		})
		if err != nil {
			t.Fatalf("failed to get min_clicks: %v", err)
		}
		if len(urls) != 2 {
			t.Fatalf("expected 2 URLs, got %d", len(urls))
		}

		// MaxClicks: 2 -> should return def (1), mno (0), pqr (0)
		urls, err = repo.GetAllURLs(models.ListOptions{
			Page:      1,
			Limit:     10,
			MaxClicks: 2,
		})
		if err != nil {
			t.Fatalf("failed to get max_clicks: %v", err)
		}
		if len(urls) != 3 {
			t.Fatalf("expected 3 URLs, got %d", len(urls))
		}

		// MinClicks: 1, MaxClicks: 4 -> should return xyz (3) and def (1)
		urls, err = repo.GetAllURLs(models.ListOptions{
			Page:      1,
			Limit:     10,
			MinClicks: 1,
			MaxClicks: 4,
		})
		if err != nil {
			t.Fatalf("failed to get range: %v", err)
		}
		if len(urls) != 2 {
			t.Fatalf("expected 2 URLs, got %d", len(urls))
		}
	})

	t.Run("Date Range Filtering", func(t *testing.T) {
		created := populateDB() // 5 URLs starting from baseTime (-10m) offset by 1m each
		baseTime := created[0].CreatedAt

		// MinDate: baseTime + 90 seconds -> should return created[2] (baseTime + 2m), created[3] (baseTime + 3m), created[4] (baseTime + 4m)
		urls, err := repo.GetAllURLs(models.ListOptions{
			Page:    1,
			Limit:   10,
			MinDate: baseTime.Add(90 * time.Second),
		})
		if err != nil {
			t.Fatalf("failed to get min_date: %v", err)
		}
		if len(urls) != 3 {
			t.Fatalf("expected 3 URLs (MinDate), got %d", len(urls))
		}

		// MaxDate: baseTime + 150 seconds -> should return created[0], created[1], created[2] (since created[2] is +2m = 120s <= 150s)
		urls, err = repo.GetAllURLs(models.ListOptions{
			Page:    1,
			Limit:   10,
			MaxDate: baseTime.Add(150 * time.Second),
		})
		if err != nil {
			t.Fatalf("failed to get max_date: %v", err)
		}
		if len(urls) != 3 {
			t.Fatalf("expected 3 URLs (MaxDate), got %d", len(urls))
		}

		// Range: baseTime + 90s to baseTime + 210s -> should return created[2] (baseTime + 2m = 120s) and created[3] (baseTime + 3m = 180s)
		urls, err = repo.GetAllURLs(models.ListOptions{
			Page:    1,
			Limit:   10,
			MinDate: baseTime.Add(90 * time.Second),
			MaxDate: baseTime.Add(210 * time.Second),
		})
		if err != nil {
			t.Fatalf("failed to get date range: %v", err)
		}
		if len(urls) != 2 {
			t.Fatalf("expected 2 URLs (Range), got %d", len(urls))
		}
	})
}
