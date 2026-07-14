package handlers

import (
	"database/sql"
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
	r := gin.Default()
	h := New(repo, config.Config{})

	r.GET("/health/api", h.HealthCheckHandler)
	r.POST("/shorten", h.ShortenHandler)
	r.GET("/:code", h.RedirectHandler)
	r.GET("/urls/all", h.ListAllHandler)
	r.DELETE("/:id", h.DeleteHandler)
	return r
}

func TestHealthCheckHandler(t *testing.T) {
	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Health Check",
			path:           "/health/api",
			expectedStatus: http.StatusOK,
			expectedBody:   "Welcome to URL Shortener Service",
		},
		{
			name:           "Not Found",
			path:           "/health/unknown",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404 page not found",
		},
	}

	r := setupTestRouter(&MockRepository{})

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
		name           string
		body           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Empty Body",
			body:           "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid request body",
		},
		{
			name:           "Malformed JSON",
			body:           `{"url":`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid request body",
		},
		{
			name:           "Missing URL Field",
			body:           `{}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid request body",
		},
		{
			name:           "Invalid URL",
			body:           `{"url":"hello"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid URL",
		},
	}

	r := setupTestRouter(&MockRepository{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body := strings.NewReader(tt.body)

			req, err := http.NewRequest(
				http.MethodPost,
				"/shorten",
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
			if !strings.Contains(recorder.Body.String(), tt.expectedBody) {
				t.Errorf("test %q: body %q, want %q", tt.name, recorder.Body.String(), tt.expectedBody)
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
		repo           *MockRepository
	}{
		{
			name:           "Redirect",
			path:           "/abc",
			expectedStatus: http.StatusSeeOther,
			expectedBody:   "https://google.com",
			repo: &MockRepository{
				GetURLByCodeFunc: func(code string) (string, error) {
					return "https://google.com", nil
				},
			},
		},
		{
			name:           "Not Found",
			path:           "/xyz",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "URL not found",
			repo: &MockRepository{
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
		repo           *MockRepository
	}{
		{
			name:           "Delete",
			path:           "/1",
			expectedStatus: http.StatusOK,
			expectedBody:   "URL deleted successfully",
			repo: &MockRepository{
				DeleteURLFunc: func(id int) error {
					return nil
				},
			},
		},
		{
			name:           "Not Found",
			path:           "/xyz",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid id",
			repo: &MockRepository{
				DeleteURLFunc: func(id int) error {
					return sql.ErrNoRows
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
		name           string
		path           string
		expectedStatus int
		expectedBody   any
		repo           *MockRepository
	}{
		{
			name:           "List All",
			path:           "/urls/all",
			expectedStatus: http.StatusOK,
			expectedBody: []models.URL{
				{
					ID:          1,
					OriginalURL: "https://google.com",
					ShortCode:   "abc",
					CreatedAt:   time.Now(),
					ClickCount:  0,
				},
			},
			repo: &MockRepository{
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
			name:           "Not Found",
			path:           "/urls/xyz",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "404 page not found",
			repo: &MockRepository{
				GetAllURLsFunc: func() ([]models.URL, error) {
					return []models.URL{}, nil
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

			switch bodyVal := tt.expectedBody.(type) {
			case []models.URL:
				for _, url := range bodyVal {
					if !strings.Contains(recorder.Body.String(), url.OriginalURL) {
						t.Fatalf("response body = %q, want it to contain %q",
							recorder.Body.String(),
							url.OriginalURL,
						)
					}
				}
			case string:
				if !strings.Contains(recorder.Body.String(), bodyVal) {
					t.Fatalf("response body = %q, want it to contain %q",
						recorder.Body.String(),
						bodyVal,
					)
				}
			}
		})
	}
}
