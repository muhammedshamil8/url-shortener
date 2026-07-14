package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/muhammedshamil8/url-shortener/internal/config"
	"github.com/muhammedshamil8/url-shortener/internal/models"
)

func setupTestRouter(repo URLRepository) *gin.Engine {
	r := gin.New()
	h := New(repo, config.Config{})

	r.GET("/api/v1/live", h.LiveHandler)
	r.GET("/api/v1/ready", h.ReadyHandler)
	r.POST("/api/v1/shorten", h.ShortenHandler)
	r.GET("/:code", h.RedirectHandler)
	r.GET("/api/v1/urls", h.ListAllHandler)
	r.DELETE("/api/v1/:id", h.DeleteHandler)
	return r
}

func TestLiveHandler(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Health Check",
			path:           "/api/v1/live",
			expectedStatus: http.StatusOK,
			expectedBody:   "alive",
		},
	}

	repo := &FakeRepository{}
	r := setupTestRouter(repo)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, tt.path, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			r.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf("got %d, want %d", recorder.Code, tt.expectedStatus)
			}
			if !strings.Contains(recorder.Body.String(), tt.expectedBody) {
				t.Fatalf("response body = %q, want it to contain %q",
					recorder.Body.String(),
					tt.expectedBody,
				)
			}
		})
	}
}

func TestShortenHandler(t *testing.T) {
	tests := []struct {
		name             string
		body             string
		expectedStatus   int
		expectedContains string
	}{
		{
			name:             "Empty Body",
			body:             "",
			expectedStatus:   http.StatusBadRequest,
			expectedContains: "Invalid request body",
		},
		{
			name:             "Malformed JSON",
			body:             `{"url":`,
			expectedStatus:   http.StatusBadRequest,
			expectedContains: "Invalid request body",
		},
		{
			name:             "Missing URL Field",
			body:             `{}`,
			expectedStatus:   http.StatusBadRequest,
			expectedContains: "Invalid request body",
		},
		{
			name:             "Invalid URL",
			body:             `{"url":"hello"}`,
			expectedStatus:   http.StatusBadRequest,
			expectedContains: "Invalid URL",
		},
	}
	repo := &FakeRepository{}
	r := setupTestRouter(repo)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body := strings.NewReader(tt.body)

			req, err := http.NewRequest(
				http.MethodPost,
				"/api/v1/shorten",
				body,
			)
			if err != nil {
				t.Fatalf("test %q: failed to create request: %v", tt.name, err)
			}
			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Errorf("test %q: status code %d, want %d", tt.name, recorder.Code, tt.expectedStatus)
			}
			if !strings.Contains(recorder.Body.String(), tt.expectedContains) {
				t.Errorf("test %q: body %q, want %q", tt.name, recorder.Body.String(), tt.expectedContains)
			}
		})
	}
}

func TestRedirectHandler(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
		repo           *FakeRepository
	}{
		{
			name:           "Redirect",
			path:           "/abc",
			expectedStatus: http.StatusSeeOther,
			expectedBody:   "https://google.com",
			repo: &FakeRepository{
				GetURLByCodeFunc: func(code string) (string, error) {
					return "https://google.com", nil
				},
			},
		},
		{
			name:           "db down",
			path:           "/xyz",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Database error",
			repo: &FakeRepository{
				GetURLByCodeFunc: func(code string) (string, error) {
					return "", errors.New("database error")
				},
			},
		},
		{
			name:           "Not Found",
			path:           "/abcd",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "URL not found",
			repo: &FakeRepository{
				GetURLByCodeFunc: func(code string) (string, error) {
					return "", sql.ErrNoRows
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupTestRouter(tt.repo)

			recorder := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, tt.path, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			r.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf("got %d, want %d", recorder.Code, tt.expectedStatus)
			}
			if !strings.Contains(recorder.Body.String(), tt.expectedBody) {
				t.Fatalf("response body = %q, want it to contain %q",
					recorder.Body.String(),
					tt.expectedBody,
				)
			}
		})
	}
}

func TestDeleteHandler(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
		repo           *FakeRepository
	}{
		{
			name:           "Delete",
			path:           "/api/v1/1",
			expectedStatus: http.StatusOK,
			expectedBody:   "URL deleted successfully",
			repo: &FakeRepository{
				DeleteURLFunc: func(id int) error {
					return nil
				},
			},
		},
		{
			name:           "Invalid ID",
			path:           "/api/v1/xyz",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "",
			repo:           &FakeRepository{},
		},
		{
			name:           "URL not found",
			path:           "/api/v1/2",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "URL not found",
			repo: &FakeRepository{
				DeleteURLFunc: func(id int) error {
					return sql.ErrNoRows
				},
			},
		},
		{
			name:           "db down",
			path:           "/api/v1/2",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Failed to delete url",
			repo: &FakeRepository{
				DeleteURLFunc: func(id int) error {
					return errors.New("database error")
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			r := setupTestRouter(tt.repo)

			req, err := http.NewRequest(http.MethodDelete, tt.path, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			r.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf("got %d, want %d", recorder.Code, tt.expectedStatus)
			}
			if !strings.Contains(recorder.Body.String(), tt.expectedBody) {
				t.Fatalf("response body = %q, want it to contain %q",
					recorder.Body.String(),
					tt.expectedBody,
				)
			}
		})
	}
}

func TestListAllHandler(t *testing.T) {
	tests := []struct {
		name             string
		path             string
		expectedStatus   int
		expectedContains []string
		repo             *FakeRepository
	}{
		{
			name:           "List All",
			path:           "/api/v1/urls",
			expectedStatus: http.StatusOK,
			expectedContains: []string{
				"https://google.com",
				"abc",
			},
			repo: &FakeRepository{
				GetAllURLsFunc: func() ([]models.URL, error) {
					return []models.URL{
						{
							ID:          1,
							OriginalURL: "https://google.com",
							ShortCode:   "abc",
							CreatedAt:   time.Now(),
							ClickCount:  0,
						},
					}, nil
				},
			},
		},
		{
			name:             "Repository Error",
			path:             "/api/v1/urls",
			expectedStatus:   http.StatusInternalServerError,
			expectedContains: []string{"Failed to get all urls"},
			repo: &FakeRepository{
				GetAllURLsFunc: func() ([]models.URL, error) {
					return nil, errors.New("database error")
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			r := setupTestRouter(tt.repo)

			req, err := http.NewRequest(http.MethodGet, tt.path, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			r.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf("got %d, want %d", recorder.Code, tt.expectedStatus)
			}
			for _, s := range tt.expectedContains {
				if !strings.Contains(recorder.Body.String(), s) {
					t.Fatalf(
						"response body = %q, want it to contain %q",
						recorder.Body.String(),
						s,
					)
				}
			}
		})
	}
}

func TestReadyHandler(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
		repo           *FakeRepository
	}{
		{
			name:           "Ready",
			path:           "/api/v1/ready",
			expectedStatus: http.StatusOK,
			expectedBody:   "ready",
			repo: &FakeRepository{
				HealthFunc: func() error {
					return nil
				},
			},
		},
		{
			name:           "Not Ready",
			path:           "/api/v1/ready",
			expectedStatus: http.StatusServiceUnavailable,
			expectedBody:   "database is unavailable",
			repo: &FakeRepository{
				HealthFunc: func() error {
					return errors.New("database unavailable")
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			r := setupTestRouter(tt.repo)

			req, err := http.NewRequest(http.MethodGet, tt.path, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}
			r.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf("got %d, want %d", recorder.Code, tt.expectedStatus)
			}
			if !strings.Contains(recorder.Body.String(), tt.expectedBody) {
				t.Fatalf("response body = %q, want it to contain %q",
					recorder.Body.String(),
					tt.expectedBody,
				)
			}
		})
	}
}
